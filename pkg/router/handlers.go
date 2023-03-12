package router

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"binance_bot/pkg/router/response"

	"binance_bot/pkg/log"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

type builtinHandlers struct {
	logger      log.Logger
	docPath     string
	buildCommit string
	buildTime   time.Time
}

func (h *builtinHandlers) livenessProbe(c *gin.Context) {
	response.Render(c, http.StatusOK, "Service Alive")
}

func (h *builtinHandlers) printVersion(c *gin.Context) {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	response.Render(c, http.StatusOK, struct {
		Hostname string `json:"hostname"`
		Commit   string `json:"commit"`
		Time     string `json:"time"`
	}{
		Hostname: hostname,
		Commit:   h.buildCommit,
		Time:     h.buildTime.In(time.UTC).Format(time.RFC3339),
	})
}

func (h *builtinHandlers) renderDoc(c *gin.Context) {
	filename := filepath.Join(h.docPath, "api.swagger.json")

	f, err := os.Open(filename)
	if err != nil {
		h.logger.Err(err).Msgf("could not open file: %v", filename)
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		h.logger.Err(err).Msgf("could not read file: %v", filename)
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	mime := mimetype.Detect(b)
	c.Writer.Header().Add("Content-Type", mime.String())

	if _, err := c.Writer.Write(b); err != nil {
		h.logger.Err(err).Msgf("could not write data")
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func (h *builtinHandlers) root(c *gin.Context) {
	response.Render(c, http.StatusOK, http.StatusText(http.StatusOK))
}

func (h *builtinHandlers) notFound(c *gin.Context) {
	response.Render(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
