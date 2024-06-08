package servers

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Swagger struct {
	httpServer *http.Server
	port       string
}

func NewSwagger(port string) *Swagger {
	return &Swagger{
		port: port,
	}
}

func (s *Swagger) Start() error {
	router := mux.NewRouter().StrictSlash(true)
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	s.httpServer = &http.Server{
		Addr:    s.port,
		Handler: router,
	}

	var err error
	log.Println("start swagger UI", s.port)

	go func() {
		err = s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: serve swagger UI: %v", err)
		}
	}()

	return nil
}

func (s *Swagger) Stop(ctx context.Context) error {
	defer log.Println("swagger server is stopped")
	return s.httpServer.Shutdown(ctx)
}
