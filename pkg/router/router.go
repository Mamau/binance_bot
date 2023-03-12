package router

import (
	"binance_bot/pkg/log"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func New(opts ...Option) *gin.Engine {
	o := options{
		logger:                 log.NewDiscardLogger(),
		docPath:                "undefined",
		middlewares:            []gin.HandlerFunc{},
		handleMethodNotAllowed: true,
		enableContextFallback:  true,
		pprof:                  false,
		pprofPrefix:            "debug/pprof",
	}
	for _, opt := range opts {
		opt(&o)
	}

	engine := gin.New()
	engine.HandleMethodNotAllowed = o.handleMethodNotAllowed
	engine.ContextWithFallback = o.enableContextFallback
	engine.Use(o.middlewares...)

	h := builtinHandlers{
		logger:      o.logger,
		docPath:     o.docPath,
		buildCommit: o.buildCommit,
		buildTime:   o.buildTime,
	}

	engine.GET("/", h.root)
	engine.GET("/live", h.livenessProbe)
	engine.GET("/doc", h.renderDoc)

	if o.pprof {
		pprof.Register(engine, o.pprofPrefix)
	}

	return engine
}
