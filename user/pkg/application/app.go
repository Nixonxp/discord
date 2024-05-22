package application

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
)

// Config - server config
type Config struct {
	GRPCPort          string
	GRPCGatewayPort   string
	HTTPPort          string
	SwaggerUIHTTPPort string
}

type App struct {
	cfg *Config

	httpListener      *echo.Echo
	swaggerUIRouter   *mux.Router
	gatewayHttpServer *http.Server
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

	var router *mux.Router
	if cfg.SwaggerUIHTTPPort != "" {
		router = mux.NewRouter().StrictSlash(true)
		sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
		router.PathPrefix("/swaggerui/").Handler(sh)
	}

	return App{
		cfg:             cfg,
		httpListener:    e,
		swaggerUIRouter: router,
	}, nil
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

			log.Println("start serve", lis.Addr())
			if err := srv.Serve(lis); err != nil {
				return fmt.Errorf("server: serve grpc: %v", err)
			}
			return nil
		})
	}

	if a.cfg.HTTPPort != "" {
		group.Go(func() error {
			log.Println("start http server", a.cfg.HTTPPort)
			err := a.httpListener.Start(a.cfg.HTTPPort)
			if err != nil {
				return fmt.Errorf("server: serve http: %v", err)
			}

			return nil
		})
	}

	if a.cfg.SwaggerUIHTTPPort != "" {
		group.Go(func() error {
			log.Println("start swagger UI", a.cfg.SwaggerUIHTTPPort)
			err := http.ListenAndServe(a.cfg.SwaggerUIHTTPPort, a.swaggerUIRouter)
			if err != nil {
				return fmt.Errorf("server: serve swagger UI: %v", err)
			}
			return nil
		})
	}

	if a.cfg.GRPCGatewayPort != "" && a.gatewayHttpServer != nil {
		group.Go(func() error {
			log.Println("start server", a.cfg.GRPCGatewayPort)
			lis, err := net.Listen("tcp", a.cfg.GRPCGatewayPort)
			if err != nil {
				return fmt.Errorf("server: failed to listen: %v", err)
			}

			if err := a.gatewayHttpServer.Serve(lis); err != nil {
				return fmt.Errorf("server: serve grpc gateway: %v", err)
			}
			return nil
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
