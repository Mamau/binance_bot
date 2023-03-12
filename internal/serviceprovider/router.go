package serviceprovider

import (
	"binance_bot/internal/config"
	app "binance_bot/pkg/application"
	"binance_bot/pkg/log"
	"binance_bot/pkg/router"
	"binance_bot/pkg/router/middleware/logging"
	"binance_bot/pkg/router/middleware/recoverer"
	"binance_bot/pkg/router/middleware/servertiming"
	"github.com/gin-gonic/gin"
)

func NewBaseRouter(conf *config.Config, logger log.Logger, version app.BuildVersion) *gin.Engine {
	return router.New(
		router.Logger(logger),
		router.DocPath(conf.App.SwaggerFolder),
		router.BuildCommit(version.Commit),
		router.BuildTime(version.Time),
		router.Middlewares(
			recoverer.New(
				recoverer.Logger(logger),
			),
			servertiming.New(),
			//timeout.New(timeout.Timeout(30*time.Second)),
			logging.New(
				logging.Level(conf.App.LogLevel),
				logging.Env(conf.App.Environment),
				logging.Fallback(logger),
				logging.Prettify(conf.App.PrettyLogs),
			),
		))
}
