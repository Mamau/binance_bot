package logging

import (
	"io"

	"binance_bot/pkg/log"
)

type Option func(*options)

type options struct {
	level       string
	env         string
	noTimestamp bool
	writer      io.Writer
	prettify    bool
	fallback    log.Logger
}

func Level(lvl string) Option {
	return func(o *options) {
		o.level = lvl
	}
}

func Env(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func NoTimestamp(flag bool) Option {
	return func(o *options) {
		o.noTimestamp = flag
	}
}

func Writer(w io.Writer) Option {
	return func(o *options) {
		o.writer = w
	}
}

func Prettify(flag bool) Option {
	return func(o *options) {
		o.prettify = flag
	}
}

func Fallback(fallback log.Logger) Option {
	return func(o *options) {
		o.fallback = fallback
	}
}
