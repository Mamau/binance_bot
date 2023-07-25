package main

import (
	"binance_bot/internal/config"
	"binance_bot/internal/consumer/eventconsumer"
	"binance_bot/internal/usecase"
	"binance_bot/pkg/application"
	"binance_bot/pkg/log"
	"binance_bot/pkg/transport/http"
	"os"
)

var id, _ = os.Hostname()

func createApp(
	hs *http.Server,
	c *config.Config,
	logger log.Logger,
	watcher *usecase.WhaleWatcher,
	consumer *eventconsumer.Consumer,
) *application.App {
	return application.New(
		application.ID(id),
		application.Name(c.App.Name),
		application.Location(c.App.TZ),
		application.Logger(logger),
		application.Servers(
			hs,
		),
		application.Daemons(
			watcher,
			consumer,
		),
	)
}

func main() {
	a, cleanup, err := newApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = a.Run(); err != nil {
		panic(err)
	}
}
