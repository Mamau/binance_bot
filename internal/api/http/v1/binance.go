package v1

import (
	"binance_bot/internal/api/http/v1/request"
	"binance_bot/internal/usecase"
	"binance_bot/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type binanceRoutes struct {
	logger log.Logger
	board  *usecase.LeaderBoard
}

func newBinanceRoutes(group *gin.RouterGroup, logger log.Logger, board *usecase.LeaderBoard) {
	b := &binanceRoutes{
		logger: logger,
		board:  board,
	}
	group.GET("/leader-position", b.leaderPosition)
	group.GET("/statistic", b.statistic)
}

func (b *binanceRoutes) statistic(c *gin.Context) {
	data := b.board.GetStatistic(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"data": data.Data})
}

func (b *binanceRoutes) leaderPosition(c *gin.Context) {
	var r request.LeaderPosition
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := b.board.GetLeaderByRequest(c.Request.Context(), r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
