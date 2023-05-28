package entity

import (
	"bytes"
	"fmt"
	"strings"
)

type DetailStatistic struct {
	Name          BetType `json:"name"`
	Count         int     `json:"count"`
	MinEntryPrice float64 `json:"min_entry_price"`
	MaxEntryPrice float64 `json:"max_entry_price"`
}

func NewDetailStatistic(name BetType, count int, minEntryPrice, maxEntryPrice float64) *DetailStatistic {
	return &DetailStatistic{
		Name:          name,
		Count:         count,
		MinEntryPrice: minEntryPrice,
		MaxEntryPrice: maxEntryPrice,
	}
}
func (s *DetailStatistic) ToTelegramMessage() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<b>Название</b>: %s\n", strings.ToUpper(string(s.Name))))
	buffer.WriteString(fmt.Sprintf("<b>Всего ставок</b>: <b>%d</b>\n", s.Count))
	buffer.WriteString(fmt.Sprintf("<b>Мин цена входа</b>: <b>%f</b>\n", s.MinEntryPrice))
	buffer.WriteString(fmt.Sprintf("<b>Макс цена входа</b>: <b>%f</b>\n", s.MaxEntryPrice))

	return buffer.String()
}

type Statistic struct {
	CurrentPrice float64 `json:"current_price"`
	Data         []*DetailStatistic
}

func NewStatistic(currentPrice float64, data ...*DetailStatistic) Statistic {
	return Statistic{
		CurrentPrice: currentPrice,
		Data:         data,
	}
}
func (s Statistic) ToTelegramMessage() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("--------------------\n")
	buffer.WriteString("<b>Биткоин</b>\n")
	buffer.WriteString(fmt.Sprintf("<b>Текущая цена</b>: %f\n", s.CurrentPrice))
	buffer.WriteString("--------------------\n")
	for i := range s.Data {
		buffer.WriteString(s.Data[i].ToTelegramMessage())
		buffer.WriteString("\n")
	}

	return buffer.String()
}
