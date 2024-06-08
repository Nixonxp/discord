package servers

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"math/rand/v2"
	"net/http"
)

type Metrics struct {
	e    *echo.Echo
	port string
}

func NewMetrics(port string) *Metrics {
	return &Metrics{
		port: port,
	}
}

func (s *Metrics) Start() error {
	s.e = echo.New()
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	s.e.GET("/health", func(c echo.Context) error {
		status := http.StatusOK
		statusMessage := "OK"

		if !isServiceOk(10) {
			status = http.StatusInternalServerError
			statusMessage = "Error"
		}

		return c.JSON(status, struct{ Status string }{Status: statusMessage})
	})

	s.e.GET("/ready", func(c echo.Context) error {
		status := http.StatusOK
		statusMessage := "OK"

		if !isServiceOk(5) {
			status = http.StatusInternalServerError
			statusMessage = "Error"
		}

		return c.JSON(status, struct{ Status string }{Status: statusMessage})
	})

	log.Println("start metrics", s.port)
	go func() {
		err := s.e.Start(s.port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: serve http: %v", err)
		}
	}()

	return nil
}

func (s *Metrics) Stop(ctx context.Context) error {
	defer log.Println("metrics server is stopped")
	return s.e.Shutdown(ctx)
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
