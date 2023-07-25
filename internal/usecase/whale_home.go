package usecase

import (
	"binance_bot/internal/entity"
	"binance_bot/internal/webapi"
	"binance_bot/pkg/log"
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"
)

type WhaleHome struct {
	logger        log.Logger
	cache         *cache.Cache
	blockchainAPI *webapi.Blockchain
}

func NewWhaleHome(logger log.Logger, blockchainAPI *webapi.Blockchain) *WhaleHome {
	return &WhaleHome{
		logger:        logger,
		blockchainAPI: blockchainAPI,
		cache:         commonCache,
	}
}

func (wh *WhaleHome) GetWhaleActionFromCache() ([]entity.WhaleAction, error) {
	var wa = make([]entity.WhaleAction, 0, len(entity.WhaleList))
	for _, v := range entity.WhaleList {
		whale, ok := wh.cache.Get(v.Wallet)
		if !ok {
			continue
		}
		wAction, ok := whale.(entity.WhaleAction)
		if !ok {
			wh.logger.Err(errors.New("error convert whaleAction")).
				Msg("error while convert whaleAction")
			continue
		}
		wa = append(wa, wAction)
	}
	return wa, nil
}
func (wh *WhaleHome) GetWhaleAction(ctx context.Context) ([]entity.WhaleAction, error) {
	var wa = make([]entity.WhaleAction, 0, len(entity.WhaleList))
	data, err := wh.GetWhaleTransactions(ctx)
	if err != nil {
		return nil, err
	}
	for i := range data {
		whale := entity.FindWhale(data[i].Hash)
		for j := range data[i].AccountTransactions {
			transactionType := entity.TransactionType(data[i].AccountTransactions[j].Type)
			if data[i].AccountTransactions[j].Value == "0" {
				break
			}
			num := new(big.Int)
			num, ok := num.SetString(data[i].AccountTransactions[j].Value, 10)
			if !ok {
				return nil, fmt.Errorf("error while convert string to big.Int")
			}

			timeStamp, err := strconv.ParseInt(data[i].AccountTransactions[j].Timestamp, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error while convert time to int64")
			}

			tm := parseTimeStamp(timeStamp)
			if !isDateWithinAWeek(tm) {
				break
			}

			item := entity.NewWhaleAction(whale.Name, whale.Position, transactionType, num, tm, data[i].Hash)
			wa = append(wa, item)
			break
		}
	}
	sort.Sort(entity.ByDate(wa))

	return wa, nil
}

func (wh *WhaleHome) GetWhaleTransactions(ctx context.Context) ([]entity.AccountTransaction, error) {
	var at = make([]entity.AccountTransaction, 0, len(entity.WhaleList))

	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, v := range entity.WhaleList {
		wg.Add(1)
		go func(item entity.Whale) {
			defer wg.Done()
			whale, err := wh.blockchainAPI.GetWhale(ctx, item.Wallet)
			if err != nil {
				wh.logger.Err(err).
					Str("name", item.Name).
					Str("wallet", item.Wallet).Msg("error while get whale")
				return
			}
			mutex.Lock()
			at = append(at, whale)
			mutex.Unlock()
		}(v)
	}
	wg.Wait()

	return at, nil
}

func parseTimeStamp(timeStamp int64) time.Time {
	numStr := strconv.FormatInt(timeStamp, 10)
	numDigits := len(numStr)
	if numDigits > 10 {
		return time.UnixMilli(timeStamp)
	}
	return time.Unix(timeStamp, 0)
}

func isDateWithinAWeek(date time.Time) bool {
	currentTime := time.Now()
	weekAgo := currentTime.AddDate(0, 0, -7)
	return date.After(weekAgo)
}
