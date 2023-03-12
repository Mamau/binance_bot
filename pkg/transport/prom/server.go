package prom

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"binance_bot/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*http.Server
	opts options
}

func New(opts ...Option) *Server {
	o := options{
		port:    "9100",
		handle:  "/metrics",
		guiPort: "9090",
		logger:  log.NewDiscardLogger(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	mux := http.NewServeMux()
	mux.Handle(o.handle, promhttp.Handler())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", o.port),
		Handler: mux,
	}

	return &Server{
		opts:   o,
		Server: srv,
	}
}
func (s *Server) Start() error {
	s.opts.logger.Info().Msg("prometheus server started")

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.opts.logger.Info().Msg("server prometheus stopping")
	return s.Server.Shutdown(ctx)
}
