package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"binance_bot/pkg/log"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		ctx:               context.Background(),
		sigs:              []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT},
		stopTimeout:       10 * time.Second,
		daemonStopTimeout: 10 * time.Second,
		logger:            log.NewDiscardLogger(),
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

func (a *App) Run() error {
	if err := a.setLocation(); err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(a.ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			sctx, cancel := context.WithTimeout(a.opts.ctx, a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start()
		})
	}
	wg.Wait()

	a.startDaemons()

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := a.Stop()
				if err != nil {
					return fmt.Errorf("failed to stop app: %s", err)
				}
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *App) startDaemons() {
	for _, dmn := range a.opts.daemons {
		go func(d Daemon) {
			d.Run()
		}(dmn)
	}
}

func (a *App) Stop() error {
	wg := sync.WaitGroup{}
	for _, dmn := range a.opts.daemons {
		wg.Add(1)
		go func(d Daemon) {
			ctx, cancel := context.WithTimeout(a.ctx, a.opts.daemonStopTimeout)
			defer cancel()
			if err := d.Terminate(ctx); err != nil {
				a.opts.logger.Err(err).Msg("error while terminate daemon")
			}
			wg.Done()
		}(dmn)
	}
	wg.Wait()

	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) ID() string { return a.opts.id }

func (a *App) Name() string { return a.opts.name }

func (a *App) Version() string { return a.opts.version }

func (a *App) setLocation() error {
	loc, err := time.LoadLocation(a.opts.location)
	if err != nil {
		return fmt.Errorf("error while set timezone: %s", err.Error())
	}
	time.Local = loc

	return nil
}
