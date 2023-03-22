package entity

import (
	"bytes"
	"fmt"
)

type Statistic struct {
	CountLong     int     `json:"count_long"`
	CountShort    int     `json:"count_short"`
	CountBet      int     `json:"count_bet"`
	MinEntryPrice float64 `json:"min_entry_price"`
	MaxEntryPrice float64 `json:"max_entry_price"`
}

func NewStatistic(
	CountBet int,
	CountShort int,
	CountLong int,
	MinEntryPrice float64,
	MaxEntryPrice float64,
) Statistic {
	return Statistic{
		CountBet:      CountBet,
		CountShort:    CountShort,
		CountLong:     CountLong,
		MinEntryPrice: MinEntryPrice,
		MaxEntryPrice: MaxEntryPrice,
	}
}
func (s Statistic) ToTelegramMessage() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<b>Лонгов</b>: %d\n", s.CountLong))
	buffer.WriteString(fmt.Sprintf("<b>Шортов</b>: %d\n", s.CountShort))
	buffer.WriteString(fmt.Sprintf("<b>Всего ставок</b>: <b>%d</b>\n", s.CountBet))
	buffer.WriteString(fmt.Sprintf("<b>Мин цена входа</b>: <b>%f</b>\n", s.MinEntryPrice))
	buffer.WriteString(fmt.Sprintf("<b>Макс цена входа</b>: <b>%f</b>\n", s.MaxEntryPrice))

	return buffer.String()
}
