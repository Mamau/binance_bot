package response

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, status int, v any) {
	c.JSON(status, map[string]any{"data": v})
}

func RenderError(c *gin.Context, err error) {
	var netErr net.Error

	switch {
	case errors.As(err, &netErr):
		c.JSON(http.StatusGatewayTimeout, map[string]any{"error": http.StatusText(http.StatusGatewayTimeout)})
	default:
		c.JSON(http.StatusInternalServerError, map[string]any{"error": http.StatusText(http.StatusInternalServerError)})
	}
}

func RenderErrorf(c *gin.Context, status int, format string, a ...any) {
	c.JSON(status, map[string]any{"error": fmt.Sprintf(format, a...)})
}
