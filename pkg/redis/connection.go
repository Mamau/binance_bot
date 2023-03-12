package redis

import (
	"fmt"

	"binance_bot/pkg/log"

	r "github.com/go-redis/redis/v8"
)

func New(opts ...Option) *r.Client {
	o := options{
		db:     0,
		logger: NewRedisLogger(log.NewDiscardLogger()),
	}

	for _, opt := range opts {
		opt(&o)
	}

	r.SetLogger(o.logger)

	rc := r.NewClient(&r.Options{
		Addr:     fmt.Sprintf("%s:%s", o.host, o.port),
		DB:       o.db,
		Username: o.username,
		Password: o.password,
	})

	return rc
}
