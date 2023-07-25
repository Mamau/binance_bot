package webapi

import (
	"binance_bot/internal/config"
	"binance_bot/internal/entity"
	"binance_bot/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Blockchain - данныее с этого ресурса https://www.blockchain.com
type Blockchain struct {
	log  log.Logger
	conf *config.Config
}

func NewBlockchain(log log.Logger, conf *config.Config) *Blockchain {
	return &Blockchain{
		log:  log,
		conf: conf,
	}
}

func (b *Blockchain) GetWhale(ctx context.Context, walletAddress string) (entity.AccountTransaction, error) {
	var at entity.AccountTransaction
	path := fmt.Sprintf("v2/eth/data/account/%s/wallet?page=0&size=10", walletAddress)

	response, err := b.doRequest(ctx, path)
	if err != nil {
		return at, err
	}

	defer func() { _ = response.Body.Close() }()

	if err := json.NewDecoder(response.Body).Decode(&at); err != nil {
		return at, err
	}
	return at, nil
}

func (b *Blockchain) doRequest(ctx context.Context, path string) (*http.Response, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.blockchain.info",
		Path:   path,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
