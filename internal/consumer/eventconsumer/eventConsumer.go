package eventconsumer

import (
	"binance_bot/internal/config"
	"binance_bot/internal/events"
	"binance_bot/internal/events/telegram"
	"binance_bot/pkg/log"
	"context"
	"time"
)

type Consumer struct {
	fetcher   *telegram.Processor
	batchSize int
	logger    log.Logger
	done      chan struct{}
}

func NewConsumer(f *telegram.Processor, conf *config.Config, l log.Logger) *Consumer {
	return &Consumer{
		fetcher:   f,
		batchSize: conf.App.BatchSize,
		logger:    l,
	}
}

func (c *Consumer) Terminate(ctx context.Context) error {
	c.logger.Info().Msg("terminating observe post leader bot watcher")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return nil
	}
}

func (c *Consumer) Run() {
	c.logger.Info().Msg("event consumer started")
	for {
		select {
		case <-c.done:
			return
		default:
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				c.logger.Err(err).Msg("error while fetch events")
				continue
			}
			if len(gotEvents) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.handleEvents(gotEvents); err != nil {
				c.logger.Err(err).Msg("error while handle events")
				continue
			}
		}
	}

}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		if err := c.fetcher.Process(event); err != nil {
			c.logger.Err(err).Msg("error while process events")
			continue
		}
	}
	return nil
}
