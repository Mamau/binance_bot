package http

import (
	"fmt"
	"net/http"
	"time"

	"binance_bot/pkg/log"
)

type Option func(o *options)

type options struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	handler      http.Handler
	logger       log.Logger
	addr         string
}

func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

func Handler(h http.Handler) Option {
	return func(o *options) { o.handler = h }
}

func ReadTimeout(t time.Duration) Option {
	return func(o *options) { o.readTimeout = t }
}

func WriteTimeout(t time.Duration) Option {
	return func(o *options) { o.writeTimeout = t }
}

func Addr(a string) Option {
	return func(o *options) { o.addr = fmt.Sprintf(":%s", a) }
}
