package recoverer

import (
	"binance_bot/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func New(opts ...Option) gin.HandlerFunc {
	o := options{
		logger: log.NewDiscardLogger(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					o.logger.Error().Msgf("panic recovered: %v", x)
				case error:
					o.logger.Err(x).Msgf("panic recovered")
				default:
					o.logger.Error().Msgf("panic recovered (%T): %v", x, x)
				}

				c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				c.Abort()
			}
		}()

		c.Next()
	}
}
