package redis

import (
	"context"

	"binance_bot/pkg/log"
)

type Logging interface {
	Printf(ctx context.Context, format string, v ...interface{})
}

type Option func(o *options)

type options struct {
	host     string
	port     string
	username string
	password string
	db       int
	logger   Logging
}

func DB(db int) Option {
	return func(o *options) { o.db = db }
}
func Host(host string) Option {
	return func(o *options) { o.host = host }
}
func Port(port string) Option {
	return func(o *options) { o.port = port }
}
func Username(username string) Option {
	return func(o *options) { o.username = username }
}
func Password(password string) Option {
	return func(o *options) { o.password = password }
}
func Logger(l log.Logger) Option {
	return func(o *options) { o.logger = NewRedisLogger(l) }
}
