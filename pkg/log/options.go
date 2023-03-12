package log

import (
	"io"
	"time"
)

type Option func(*options)

type options struct {
	level       string
	env         string
	buildCommit string
	buildTime   time.Time
	noTimestamp bool
	writer      io.Writer
	prettify    bool
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

func BuildCommit(commit string) Option {
	return func(o *options) {
		o.buildCommit = commit
	}
}

func BuildTime(t time.Time) Option {
	return func(o *options) {
		o.buildTime = t
	}
}

func NoTimestamp(no bool) Option {
	return func(o *options) {
		o.noTimestamp = no
	}
}

func Writer(w io.Writer) Option {
	return func(o *options) {
		o.writer = w
	}
}

func Prettify(prettify bool) Option {
	return func(o *options) {
		o.prettify = prettify
	}
}
