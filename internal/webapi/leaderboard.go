package webapi

import (
	"binance_bot/internal/config"
	"binance_bot/internal/entity"
	"binance_bot/pkg/log"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type LeaderBoard struct {
	log  log.Logger
	conf *config.Config
}

func NewWebData(log log.Logger, conf *config.Config) *LeaderBoard {
	return &LeaderBoard{
		log:  log,
		conf: conf,
	}
}

func (l *LeaderBoard) GetPositions(ctx context.Context, encryptedUid []string) []*entity.LeaderBoard {
	var ops []entity.OtherPosition
	for _, v := range encryptedUid {
		ops = append(ops, entity.OtherPosition{
			EncryptedUid: v,
			TradeType:    l.conf.App.TradeType,
		})
	}

	var result []*entity.LeaderBoard
	var wg sync.WaitGroup
	for _, v := range ops {
		wg.Add(1)
		go func(op entity.OtherPosition) {
			defer wg.Done()
			response, err := l.doRequest(ctx, op)
			if err != nil {
				l.log.Err(err).Msg("error while do request in cycle")
				return
			}

			defer func() { _ = response.Body.Close() }()

			var lb entity.LeaderBoard
			if err := json.NewDecoder(response.Body).Decode(&lb); err != nil {
				l.log.Err(err).Msg("error while decode in cycle")
				return
			}
			result = append(result, &lb)
		}(v)
	}
	wg.Wait()

	return result
}

func (l *LeaderBoard) GetPosition(ctx context.Context, encryptedUid, tradeType string) (*entity.LeaderBoard, error) {
	op := entity.OtherPosition{
		EncryptedUid: encryptedUid,
		TradeType:    tradeType,
	}

	response, err := l.doRequest(ctx, op)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	var lb entity.LeaderBoard
	if err := json.NewDecoder(response.Body).Decode(&lb); err != nil {
		return nil, err
	}

	return &lb, nil
}

func (l *LeaderBoard) doRequest(ctx context.Context, op entity.OtherPosition) (*http.Response, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "www.binance.com",
		Path:   "bapi/futures/v1/public/future/leaderboard/getOtherPosition",
	}

	result, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(result))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}

	client := http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= http.StatusNetworkAuthenticationRequired {
		var lberror entity.LeaderBoardError
		if err := json.NewDecoder(resp.Body).Decode(&lberror); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("probably remote server is down, status code: %d, message: %s", resp.StatusCode, lberror.Message)
	}

	return resp, nil
}
