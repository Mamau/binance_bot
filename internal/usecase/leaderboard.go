package usecase

import (
	"binance_bot/internal/api/http/v1/request"
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
	if len(lb.Data.OtherPositionRetList) == 0 {
		if _, ok := l.store.Get("empty_key"); ok {
			return
		}

		if err := l.tg.SendMessage(l.conf.App.TelegramUserID, "чувак закрыл сделку или их нет"); err != nil {
			l.logger.Err(err).Msg("error while send message")
		}
		l.store.Set("empty_key", true, cache.NoExpiration)
		return
	}

	for _, v := range lb.Data.OtherPositionRetList {
		pos := entity.NewPosition(v.Symbol, v.Amount, v.EntryPrice, v.Pnl, v.Leverage)

		if _, ok := l.store.Get(getKey(pos)); ok {
			continue
		}
		// чистим все хранилище
		l.store.Flush()
		l.store.Set(getKey(pos), true, cache.NoExpiration)

		if err := l.tg.SendMessage(l.conf.App.TelegramUserID, pos.ToTelegramMessage()); err != nil {
			l.logger.Err(err).Msg("error while send message")
			continue
		}
	}
}
func (l *LeaderBoard) GetStatistic(ctx context.Context) entity.Statistic {
	var uids []string
	uids, err := l.GetAllLeadersUids(ctx)
	if err != nil {
		uids = l.conf.App.EncryptedUids
	}
	data := l.binanceApi.GetPositions(ctx, uids)

	shortStat := entity.NewDetailStatistic(entity.Short, 0, 0, 0)
	longStat := entity.NewDetailStatistic(entity.Long, 0, 0, 0)
	var markPrice float64

	for _, v := range data {
		if len(v.Data.OtherPositionRetList) == 0 {
			continue
		}

		for _, bet := range v.Data.OtherPositionRetList {
			if betType := entity.DetermineBetName(bet.Symbol); betType != entity.BTC {
				continue
			}
			markPrice = bet.MarkPrice
			if betType := entity.DetermineBetType(bet.Amount); betType == entity.Short {
				shortStat.Count++
				if bet.EntryPrice < shortStat.MinEntryPrice || shortStat.MinEntryPrice == 0 {
					shortStat.MinEntryPrice = bet.EntryPrice
				}
				if bet.EntryPrice > shortStat.MaxEntryPrice || shortStat.MaxEntryPrice == 0 {
					shortStat.MaxEntryPrice = bet.EntryPrice
				}
				continue
			}

			longStat.Count++
			if bet.EntryPrice < longStat.MinEntryPrice || longStat.MinEntryPrice == 0 {
				longStat.MinEntryPrice = bet.EntryPrice
			}
			if bet.EntryPrice > longStat.MaxEntryPrice || longStat.MaxEntryPrice == 0 {
				longStat.MaxEntryPrice = bet.EntryPrice
			}
		}
	}

	return entity.NewStatistic(markPrice, longStat, shortStat)
}
func (l *LeaderBoard) GetLeaderByRequest(ctx context.Context, r request.LeaderPosition) (*entity.LeaderBoard, error) {
	encrUid := l.conf.App.EncryptedUid
	trType := l.conf.App.TradeType
	if r.EncryptedUid != "" {
		encrUid = r.EncryptedUid
	}
	if r.TradeType != "" {
		trType = r.TradeType
	}

	lb, err := l.binanceApi.GetPosition(ctx, encrUid, trType)
	if err != nil {
		l.logger.Err(err).Msg("error while get leader by request from binance")
		return nil, err
	}

	return lb, nil
}
func (l *LeaderBoard) GetAllLeadersUids(ctx context.Context) ([]string, error) {
	data, err := l.binanceApi.GetAllLeaders(ctx)
	if err != nil {
		l.logger.Err(err).Msg("error while get position from binance")
		return nil, err
	}
	var uids []string
	for i := range data.Data {
		uids = append(uids, data.Data[i].EncryptedUid)
	}

	return uids, nil
}
func (l *LeaderBoard) GetLeader(ctx context.Context) (*entity.LeaderBoard, error) {
	lb, err := l.binanceApi.GetPosition(ctx, l.conf.App.EncryptedUid, l.conf.App.TradeType)
	if err != nil {
		l.logger.Err(err).Msg("error while get position from binance")
		return nil, err
	}

	return lb, nil
}

func getKey(p entity.Position) string {
	return fmt.Sprintf("%s_%s_%f", p.Symbol, p.BetType, p.EntryPrice)
}
