package servertiming

import (
	"binance_bot/pkg/router/httpsnoop"
	"github.com/gin-gonic/gin"
	orig "github.com/mitchellh/go-server-timing"
)

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &Header{&orig.Header{}}
		var headerWritten bool

		headers := c.Writer.Header()

		c.Request = c.Request.WithContext(
			NewContext(c.Request.Context(), h),
		)

		total := h.NewMetric("total").Start()

		c.Writer = httpsnoop.Wrap(c.Writer, httpsnoop.Hooks{
			WriteHeader: func(original httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					total.Stop()
					writeHeader(headers, h)
					headerWritten = true
					original(code)
				}
			},
			Write: func(original httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(b []byte) (int, error) {
					if !headerWritten {
						total.Stop()
						writeHeader(headers, h)
						headerWritten = true
					}
					return original(b)
				}
			},
		})

		c.Next()
	}
}
