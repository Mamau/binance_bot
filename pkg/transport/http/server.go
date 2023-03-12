package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"binance_bot/pkg/log"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
	_defaultAddr         = ":8080"
)

type Server struct {
	*http.Server
	opts options
}

func New(opts ...Option) *Server {
	o := options{
		readTimeout:  _defaultReadTimeout,
		writeTimeout: _defaultWriteTimeout,
		addr:         _defaultAddr,
		handler:      http.HandlerFunc(defaultHandler),
		logger:       log.NewDiscardLogger(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	srv := &http.Server{
		Addr:         o.addr,
		Handler:      o.handler,
		ReadTimeout:  o.readTimeout,
		WriteTimeout: o.writeTimeout,
	}

	return &Server{
		opts:   o,
		Server: srv,
	}
}

func (s *Server) Start() error {
	s.opts.logger.Info().Msgf("http server started at port %s", s.opts.addr)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.opts.logger.Info().Msg("server http stopping")
	return s.Server.Shutdown(ctx)
}

func defaultHandler(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("hello, you need a routing!\n")); err != nil {
		fmt.Println(err)
	}
}
