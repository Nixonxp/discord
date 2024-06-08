package server

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/server/configs"
	"github.com/Nixonxp/discord/server/internal/app/errors"
	"log"
	"sync"
	"time"
)

type Service interface {
	Init(ctx context.Context, cfg *config.Config) error
	Close(ctx context.Context) error
	Ident() string
}

type Resources struct {
	Services        []Service
	ShutdownTimeout time.Duration
	Cfg             *config.Config
}

const (
	defaultShutdownTimeout = time.Millisecond * 15000
)

func (s *Resources) Init(ctx context.Context) error {
	if err := s.initAllServices(ctx); err != nil {
		return err
	}

	s.defaultConfigs()
	return nil
}

func (s *Resources) defaultConfigs() {
	if s.ShutdownTimeout == 0 {
		s.ShutdownTimeout = defaultShutdownTimeout
	}
}

func (s *Resources) initAllServices(ctx context.Context) (initError error) {
	initCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range s.Services {
		log.Printf("run service - %s", s.Services[i].Ident())
		if err := s.Services[i].Init(initCtx, s.Cfg); err != nil {
			return fmt.Errorf("error has occurred in the '%s' service: %w", s.Services[i].Ident(), err)
		}
	}

	return nil
}

func (s *Resources) Release() error {
	return s.release()
}

func (s *Resources) release() error {
	shCtx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	var p parallelRun
	for i := range s.Services {
		var service = s.Services[i]
		log.Printf("stop service - %s", s.Services[i].Ident())
		p.do(shCtx, service.Ident(), func(_ context.Context) error {
			return service.Close(shCtx)
		})
	}

	var errCh = make(chan error)
	go func() {
		defer close(errCh)
		if err := p.wait(); err != nil {
			errCh <- err
		}
	}()

	for {
		select {
		case err, ok := <-errCh:
			if ok {
				return err
			}
			return nil
		case <-shCtx.Done():
			return shCtx.Err()
		}
	}
}

type parallelRun struct {
	mux sync.Mutex
	wg  sync.WaitGroup
	err errors.ArrError
}

func (p *parallelRun) do(ctx context.Context, ident string, f func(context.Context) error) {
	p.wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				p.mux.Lock()
				p.err = append(p.err, fmt.Errorf("unhandled error has occurred in the '%s' service: %v", ident, r))
				p.mux.Unlock()
			}
			p.wg.Done()
		}()
		if err := f(ctx); err != nil {
			p.mux.Lock()
			p.err = append(p.err, fmt.Errorf("error has occurred in the '%s' service: %w", ident, err))
			p.mux.Unlock()
		}
	}()
}

func (p *parallelRun) wait() error {
	p.wg.Wait()
	if len(p.err) > 0 {
		return p.err
	}
	return nil
}
