package serviceprovider

import (
	v1 "binance_bot/internal/api/http/v1"
	telegram2 "binance_bot/internal/clients/telegram"
	"binance_bot/internal/config"
	"binance_bot/internal/consumer/eventconsumer"
	"binance_bot/internal/events/telegram"
	"binance_bot/internal/usecase"
	"binance_bot/internal/webapi"

	"binance_bot/pkg/application"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewBaseRouter,
	usecase.NewWhaleHome,
	webapi.NewBlockchain,
	usecase.NewWhaleWatcher,
	application.GetBuildVersion,
	config.GetConfig,
	config.GetAppConfig,
	NewHttp,
	v1.NewRouter,
	NewLogger,
	webapi.NewWebData,
	usecase.NewLeaderBoard,
	usecase.NewLeaderBoardWatcher,
	telegram.NewProcessor,
	telegram2.NewClient,
	eventconsumer.NewConsumer,
)
