package application

import (
	"context"
	"os"
	"time"

	"binance_bot/pkg/log"
)

type Option func(o *options)

type options struct {
	id       string
	name     string
	version  string
	location string

	ctx  context.Context
	sigs []os.Signal

	logger            log.Logger
	servers           []Server
	daemons           []Daemon
	stopTimeout       time.Duration
	daemonStopTimeout time.Duration
}

func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

func ID(id string) Option {
	return func(o *options) { o.id = id }
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Location(loc string) Option {
	return func(o *options) {
		o.location = loc
	}
}

func Servers(srv ...Server) Option {
	return func(o *options) { o.servers = srv }
}

func Daemons(dmn ...Daemon) Option {
	return func(o *options) { o.daemons = dmn }
}

func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

func DaemonStopTimeout(t time.Duration) Option {
	return func(o *options) { o.daemonStopTimeout = t }
}
