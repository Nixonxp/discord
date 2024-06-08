package server

import (
	"context"
	config "github.com/Nixonxp/discord/channel/configs"
	"github.com/Nixonxp/discord/channel/internal/app/errors"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	MainServer            *MainServer
	Resources             *Resources
	TerminationTimeout    time.Duration
	InitializationTimeout time.Duration
	cfg                   *config.Config

	appState  int32
	mux       sync.Mutex
	err       error
	halt      chan struct{}
	appCtx    context.Context
	appCancel context.CancelFunc
}

const (
	defaultTerminationTimeout    = time.Second
	defaultInitializationTimeout = time.Second * 15
)

func NewApplication(cfg *config.Config) *Application {
	var srv MainServer
	var svc = Resources{
		Services: []Service{
			&srv.tracer,
			&srv.logger,
			&srv.mongo,
		},
		ShutdownTimeout: terminationTimeout,
		Cfg:             cfg,
	}

	return &Application{
		MainServer:         &srv,
		Resources:          &svc,
		TerminationTimeout: terminationTimeout,
		cfg:                cfg,
	}
}

func (a *Application) Run() error {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if err := a.init(ctx); err != nil {
		a.err = err
		return err
	}

	var servicesRunning = make(chan struct{})

	log.Print("application starting...")
	a.setError(a.run(ctx))

	// shutdown
	close(servicesRunning)
	if a.Resources != nil {
		select {
		case <-servicesRunning:
		case <-time.After(a.TerminationTimeout):
		}
		log.Print("resources release...")
		a.setError(a.Resources.Release())
	}
	log.Print("application stopped")
	return a.getError()
}

func (a *Application) init(ctx context.Context) error {
	if a.TerminationTimeout == 0 {
		a.TerminationTimeout = defaultTerminationTimeout
	}
	if a.InitializationTimeout == 0 {
		a.InitializationTimeout = defaultInitializationTimeout
	}
	a.halt = make(chan struct{})
	a.appCtx, a.appCancel = context.WithCancel(context.Background())
	if a.Resources != nil {
		ctx, cancel := context.WithTimeout(ctx, a.InitializationTimeout)
		defer cancel()
		return a.Resources.Init(ctx)
	}
	return nil
}

func (a *Application) run(ctx context.Context) error {
	defer a.Shutdown()

	var errRun = make(chan error, 1)
	go func() {
		defer close(errRun)
		if err := a.MainServer.AppStart(ctx, a.cfg); err != nil {
			errRun <- err
		}
	}()
	var errHld = make(chan error, 1)
	go func() {
		defer close(errHld)
		select {
		// wait for os signal
		case <-ctx.Done():
			select {
			case <-time.After(a.TerminationTimeout):
				errHld <- errors.ErrTermTimeout
			case <-a.appCtx.Done():
				// ok
			}
		}
	}()
	select {
	case err, ok := <-errRun:
		if ok && err != nil {
			return err
		}
	case err, ok := <-errHld:
		if ok && err != nil {
			return err
		}
	case <-a.appCtx.Done():
		// shutdown
	}
	return nil
}

// Shutdown stops the application immediately. At this point, all calculations should be completed.
func (a *Application) Shutdown() {
	a.appCancel()
}

func (a *Application) getError() error {
	var err error
	a.mux.Lock()
	err = a.err
	a.mux.Unlock()
	return err
}

func (a *Application) setError(err error) {
	if err == nil {
		return
	}
	a.mux.Lock()
	if a.err == nil {
		a.err = err
	}
	a.mux.Unlock()
	a.Shutdown()
}
