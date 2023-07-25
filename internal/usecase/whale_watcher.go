package usecase

import (
	telegramClient "binance_bot/internal/clients/telegram"
	"binance_bot/internal/config"
	"binance_bot/internal/entity"
	"binance_bot/pkg/log"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var commonCache = cache.New(cache.NoExpiration, cache.NoExpiration)

type WhaleWatcher struct {
	service *WhaleHome
	logger  log.Logger
	cache   *cache.Cache
	tg      *telegramClient.Client
	conf    *config.Config
	done    chan struct{}
}

func NewWhaleWatcher(s *WhaleHome, l log.Logger, tg *telegramClient.Client, conf *config.Config) *WhaleWatcher {
	return &WhaleWatcher{
		service: s,
		logger:  l,
		tg:      tg,
		conf:    conf,
		cache:   commonCache,
	}
}

func (ww *WhaleWatcher) Run() {
	defer close(ww.done)
	ticker := time.NewTicker(time.Second * 11)

	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				data, err := ww.service.GetWhaleAction(ctx)
				if err != nil {
					ww.logger.Err(err).Msg("error while get whale action by daemon")
					return
				}
				for i := range data {
					ww.saveOrUpdateCache(data[i])
				}
			}()
		}
	}
}
func (ww *WhaleWatcher) saveOrUpdateCache(data entity.WhaleAction) {
	if cachedWhale, found := ww.cache.Get(data.Hash); found {
		wa, ok := cachedWhale.(entity.WhaleAction)
		if !ok {
			ww.logger.Err(errors.New("error convert whaleAction")).Msg("error while convert whaleAction")
			return
		}

		if data.Date.After(wa.Date) {
			ww.cache.Set(data.Hash, data, cache.NoExpiration)
			ww.logger.Info().
				Str("name", data.WhaleName).
				Str("hash", data.Hash).
				Msg("update cache")

			var buffer bytes.Buffer
			buffer.WriteString("<b>Обновление</b>\n")
			buffer.WriteString(fmt.Sprintf("<b>Имя</b>: %s\n", data.WhaleName))
			buffer.WriteString(fmt.Sprintf("<b>Позиция</b>: %s\n", data.WhalePosition))
			buffer.WriteString(fmt.Sprintf("<b>Тип</b>: %s\n", data.Type))
			buffer.WriteString(fmt.Sprintf("<b>ETH</b>: %f\n", data.ValueETH))
			buffer.WriteString(fmt.Sprintf("<b>Дата</b>: %s\n", data.Date.Format("02.01.2006 15:04:05")))
			buffer.WriteString("\n")

			if err := ww.tg.SendMessage(ww.conf.App.TelegramUserID, buffer.String()); err != nil {
				ww.logger.Err(err).Msg("error while send message for update")
				return
			}
		}
		return
	}

	ww.cache.Set(data.Hash, data, cache.NoExpiration)
}
func (ww *WhaleWatcher) Terminate(ctx context.Context) error {
	ww.logger.Info().Msg("terminating observe whale watcher")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ww.done:
		return nil
	}
}
