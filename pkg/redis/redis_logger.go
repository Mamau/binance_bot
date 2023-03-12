package redis

import (
	"context"
	"fmt"

	"binance_bot/pkg/log"
)

// RedisLogger - переопределение логгера редиса, чтобы было всё в одном формате
type RedisLogger struct {
	logger log.Logger
}

func NewRedisLogger(logger log.Logger) Logging {
	return &RedisLogger{logger}
}

// Printf - метод печати ошибки редиса в stdout
func (r *RedisLogger) Printf(_ context.Context, format string, v ...interface{}) {
	r.logger.Warn().Msg(fmt.Sprintf(format, v...))
}
