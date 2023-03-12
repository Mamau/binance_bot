package router

import (
	"time"

	"binance_bot/pkg/log"
	"github.com/gin-gonic/gin"
)

type Option func(o *options)

type options struct {
	env                    string
	logger                 log.Logger
	docPath                string
	buildCommit            string
	buildTime              time.Time
	middlewares            []gin.HandlerFunc
	handleMethodNotAllowed bool
	enableContextFallback  bool
	pprof                  bool
	pprofPrefix            string
}

func PprofPrefix(prefix string) Option {
	return func(o *options) {
		o.pprofPrefix = prefix
	}
}

func Pprof(enable bool) Option {
	return func(o *options) {
		o.pprof = enable
	}
}

func Env(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func DocPath(path string) Option {
	return func(o *options) {
		o.docPath = path
	}
}

func BuildCommit(buildCommit string) Option {
	return func(o *options) {
		o.buildCommit = buildCommit
	}
}

func BuildTime(t time.Time) Option {
	return func(o *options) {
		o.buildTime = t
	}
}

func Middlewares(middlewares ...gin.HandlerFunc) Option {
	return func(o *options) {
		o.middlewares = middlewares
	}
}

// HandleMethodNotAllowed refers to https://github.com/gin-gonic/gin/blob/v1.8.1/gin.go#L107
func HandleMethodNotAllowed(handle bool) Option {
	return func(o *options) {
		o.handleMethodNotAllowed = handle
	}
}

// EnableContextFallback refers to https://github.com/gin-gonic/gin/blob/v1.8.1/gin.go#L151
func EnableContextFallback(flag bool) Option {
	return func(o *options) {
		o.enableContextFallback = flag
	}
}
