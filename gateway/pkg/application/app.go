package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

// Config - server config
type Config struct {
	GRPCPort          string
	GRPCGatewayPort   string
	HTTPPort          string
	SwaggerUIHTTPPort string
	DebugPort         string
}

type App struct {
	cfg *Config

	httpListener      *echo.Echo
	swaggerUISrv      *http.Server
	gatewayHttpServer *http.Server
	debugSrv          *http.Server
}

func NewApp(cfg *Config) (App, error) {
	var e *echo.Echo
	if cfg.HTTPPort != "" {
		e = echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.GET("/health", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(10) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})

		e.GET("/ready", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(5) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})
	}

	var uiSrv *http.Server
	if cfg.SwaggerUIHTTPPort != "" {
		router := mux.NewRouter().StrictSlash(true)
		sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
		router.PathPrefix("/swaggerui/").Handler(sh)

		uiSrv = &http.Server{
			Addr:    cfg.SwaggerUIHTTPPort,
			Handler: router,
		}
	}

	app := App{
		cfg:          cfg,
		httpListener: e,
		swaggerUISrv: uiSrv,
	}

	if cfg.DebugPort != "" {
		{
			mux := http.NewServeMux()
			// prometheus metrics
			mux.Handle("/metrics", promhttp.Handler())
			// pprof
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

			httpServer := &http.Server{Handler: mux}
			app.debugSrv = httpServer
		}
	}

	return app, nil
}

func (a *App) AddGatewayServer(srv *http.Server) *App {
	a.gatewayHttpServer = srv
	return a
}

func (a *App) Run(ctx context.Context, srv *grpc.Server) error {
	group := errgroup.Group{}

	if srv != nil {
		group.Go(func() error {
			reflection.Register(srv)

			lis, err := net.Listen("tcp", a.cfg.GRPCPort)
			if err != nil {
				return fmt.Errorf("server: failed to listen: %v", err)
			}

			log.Println("start serve grpc server", lis.Addr())

			go func() {
				if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server: serve grpc: %v", err)
				}
			}()

			<-ctx.Done()
			srv.GracefulStop()
			log.Println("grpc server is stopped")

			return nil
		})
	}

	if a.cfg.HTTPPort != "" {
		group.Go(func() error {
			log.Println("start http server", a.cfg.HTTPPort)

			go func() {
				err := a.httpListener.Start(a.cfg.HTTPPort)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server: serve http: %v", err)
				}
			}()

			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			return a.httpListener.Shutdown(shutdownCtx)
		})
	}

	if a.cfg.DebugPort != "" {
		group.Go(func() error {
			log.Println(ctx, "start serve debug", a.cfg.DebugPort)

			lis, err := net.Listen("tcp", a.cfg.DebugPort)
			if err != nil {
				return fmt.Errorf("server: failed to listen debug: %v", err)
			}

			go func() {
				if err := a.debugSrv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server: debug serve internal: %v", err)
				}
			}()

			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			return a.debugSrv.Shutdown(shutdownCtx)
		})
	}

	if a.cfg.SwaggerUIHTTPPort != "" {
		group.Go(func() error {
			log.Println("start swagger UI", a.cfg.SwaggerUIHTTPPort)

			go func() {
				err := a.swaggerUISrv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server: serve swagger UI: %v", err)
				}
			}()

			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			return a.swaggerUISrv.Shutdown(shutdownCtx)
		})
	}

	if a.cfg.GRPCGatewayPort != "" && a.gatewayHttpServer != nil {
		group.Go(func() error {
			log.Println("start server", a.cfg.GRPCGatewayPort)
			go func() {
				if err := a.gatewayHttpServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server: serve grpc gateway: %v", err)
				}
			}()

			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			return a.gatewayHttpServer.Shutdown(shutdownCtx)
		})
	}

	return group.Wait()
}

// isServiceOk в зависимости от входящего значения вернет false, например
// передано 5, тогда (100 / 5 = 20) 20% вероятностью вернется false, для теста сервиса
func isServiceOk(probability int) bool {
	randNumber := rand.IntN(probability-1) + 1

	if randNumber == 1 {
		return false
	}

	return true
}
