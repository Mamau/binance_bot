package logging

import (
	"os"
	"time"

	"binance_bot/pkg/log"
	"github.com/gin-gonic/gin"
)

func New(opts ...Option) gin.HandlerFunc {
	o := options{
		level:    "info",
		writer:   os.Stdout,
		fallback: log.NewDiscardLogger(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	logger := log.NewLogger(
		log.Level(o.level),
		log.NoTimestamp(o.noTimestamp),
		log.Env(o.env),
		log.Writer(log.NewNonBlockingWriter(o.writer, 1000, 10*time.Millisecond, o.fallback)),
		log.Prettify(o.prettify),
	)

	return func(c *gin.Context) {
		c.Next()

		level := "info"

		log.AddContext(c, logger.WithLevel(log.ParseLevel(level))).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("incoming request")
	}
}
