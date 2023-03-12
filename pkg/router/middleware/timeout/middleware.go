package timeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func New(opts ...Option) gin.HandlerFunc {
	o := options{
		timeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), o.timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
