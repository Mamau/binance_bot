package serviceprovider

import (
	"binance_bot/internal/config"
	nh "net/http"

	"binance_bot/pkg/log"
	"binance_bot/pkg/transport/http"
)

func NewHttp(handler nh.Handler, conf *config.Config, logger log.Logger) *http.Server {
	return http.New(
		http.Logger(logger),
		http.Handler(handler),
		http.Addr(conf.App.Port),
	)
}
