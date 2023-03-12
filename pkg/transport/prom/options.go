package prom

import "binance_bot/pkg/log"

type Option func(o *options)

type options struct {
	port    string
	handle  string
	guiPort string
	logger  log.Logger
}

func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}
func Port(port string) Option {
	return func(o *options) { o.port = port }
}

func Handle(h string) Option {
	return func(o *options) { o.handle = h }
}

func GuiPort(gp string) Option {
	return func(o *options) { o.guiPort = gp }
}
