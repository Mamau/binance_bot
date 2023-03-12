package entity

import (
	"bytes"
	"fmt"
	"strings"
)

type BetType string
type BetName string

const (
	BTC   BetName = "BTCUSDT"
	SOL   BetName = "SOLUSDT"
	ETH   BetName = "ETHUSDT"
	Empty BetName = "undefined"
)
const (
	Long      BetType = "long"
	Short     BetType = "short"
	Undefined BetType = "undefined"
)

type Position struct {
	Symbol     BetName `json:"symbol"`
	Amount     float64 `json:"amount"`
	BetType    BetType `json:"bet_type"`
	EntryPrice float64 `json:"entry_price"`
	Pnl        float64 `json:"pnl"`
	Leverage   int     `json:"leverage"`
}

func NewPosition(
	symbol string,
	amount float64,
	entryPrice float64,
	pnl float64,
	leverage int,
) Position {
	return Position{
		Symbol:     determineBetName(symbol),
		BetType:    determineBetType(amount),
		Amount:     amount,
		EntryPrice: entryPrice,
		Pnl:        pnl,
		Leverage:   leverage,
	}
}
func (p Position) ToTelegramMessage() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<b>Валюта</b>: %55s\n", string(p.Symbol)))
	buffer.WriteString(fmt.Sprintf("<b>Тип ставки</b>: <b>%51s</b>\n", strings.ToUpper(string(p.BetType))))
	buffer.WriteString(fmt.Sprintf("<b>Плечи</b>: %65d\n", p.Leverage))
	buffer.WriteString(fmt.Sprintf("<b>Сколько потратил</b>: %34f\n", (p.EntryPrice*p.Amount)/float64(p.Leverage)))
	buffer.WriteString(fmt.Sprintf("<b>Цена входа</b>: %50f\n", p.EntryPrice))
	buffer.WriteString(fmt.Sprintf("<b>Сколько купил крипты</b>: %33f\n", p.Amount))
	buffer.WriteString(fmt.Sprintf("<b>Pnl</b>: %65f\n", p.Pnl))

	return buffer.String()
}
func determineBetName(name string) BetName {
	switch name {
	case string(BTC):
		return BTC
	case string(ETH):
		return ETH
	case string(SOL):
		return SOL
	default:
		return Empty
	}
}
func determineBetType(amount float64) BetType {
	if amount == 0 {
		return Undefined
	}

	if amount < 0 {
		return Short
	}
	return Long
}
