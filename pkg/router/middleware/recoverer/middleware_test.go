package recoverer_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"binance_bot/pkg/log"
	"binance_bot/pkg/router/middleware/recoverer"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecovererMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		Name         string
		Cause        any
		ExpectedLogs string
	}{
		{
			Name:         "panic using string",
			Cause:        "foo",
			ExpectedLogs: `{"level":"error","message":"panic recovered: foo"}`,
		},
		{
			Name:         "panic using error",
			Cause:        fmt.Errorf("bar"),
			ExpectedLogs: `{"level":"error","error":"bar","message":"panic recovered"}`,
		},
		{
			Name:         "panic using unexpected value",
			Cause:        820,
			ExpectedLogs: `{"level":"error","message":"panic recovered (int): 820"}`,
		},
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
			assert.NoError(t, err)

			var buf bytes.Buffer
			logger := log.NewLogger(
				log.Writer(&buf),
				log.NoTimestamp(true),
			)

			engine := gin.New()
			engine.Use(recoverer.New(recoverer.Logger(logger)))
			engine.GET("/", func(c *gin.Context) {
				panic(v.Cause)
			})

			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.Equal(t, http.StatusText(http.StatusInternalServerError), recorder.Body.String())
			assert.Contains(t, buf.String(), v.ExpectedLogs)
			assert.Equal(t, "text/plain; charset=utf-8", recorder.Header().Get("Content-Type"))
		})
	}
}
