package v1

import (
	"binance_bot/internal/usecase"
	"binance_bot/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(
	engine *gin.Engine,
	logger log.Logger,
	board *usecase.LeaderBoard,
) http.Handler {
	commonGroup := engine.Group("/api/v1")
	commonGroup.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	binanceGroup := commonGroup.Group("/binance")
	{
		newBinanceRoutes(binanceGroup, logger, board)
	}

	return engine
}
