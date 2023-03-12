package usecase

import (
	telegramClient "binance_bot/internal/clients/telegram"
	"binance_bot/internal/config"
	"binance_bot/internal/entity"
	"binance_bot/internal/webapi"
	"binance_bot/pkg/log"
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
)

type LeaderBoard struct {
	binanceApi *webapi.LeaderBoard
	logger     log.Logger
	tg         *telegramClient.Client
	store      *cache.Cache
	conf       *config.Config
}

func NewLeaderBoard(ba *webapi.LeaderBoard, logger log.Logger, tg *telegramClient.Client, conf *config.Config) *LeaderBoard {
	return &LeaderBoard{
		binanceApi: ba,
		logger:     logger,
		tg:         tg,
		store:      cache.New(cache.NoExpiration, cache.NoExpiration),
		conf:       conf,
	}
}
func (l *LeaderBoard) NotifyAboutBet(lb *entity.LeaderBoard) {
	for _, v := range lb.Data.OtherPositionRetList {
		pos := entity.NewPosition(v.Symbol, v.Amount, v.EntryPrice, v.Pnl, v.Leverage)

		if _, ok := l.store.Get(getKey(pos)); ok {
			continue
		}
		// чистим все хранилище
		l.store.Flush()
		l.store.Set(getKey(pos), true, cache.NoExpiration)

		if err := l.tg.SendMessage(l.conf.App.TelegramUserID, pos.ToTelegramMessage()); err != nil {
			l.logger.Err(err).Msg("error while send message1")
			continue
		}
	}
}
func (l *LeaderBoard) GetLeader(ctx context.Context) (*entity.LeaderBoard, error) {
	lb, err := l.binanceApi.GetPosition(ctx)
	if err != nil {
		l.logger.Err(err).Msg("error while get position from binance")
		return nil, err
	}

	return lb, nil
}

func getKey(p entity.Position) string {
	return fmt.Sprintf("%s_%s_%f", p.Symbol, p.BetType, p.EntryPrice)
}
