package telegram

import (
	"binance_bot/internal/entity"
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	HelpCmd     = "/help"
	GetInfoCmd  = "/get_info"
	GetWhaleCmd = "/get_whale"
	GetStatCmd  = "/get_statistic"
)

func (p *Processor) doCmd(text string, userID int, username string, firstName string) error {
	command, err := retrieveCommandName(text)
	if err != nil {
		return err
	}
	log.Info().Msgf("got new command '%s' from '%s'", command, username)

	switch command {
	case HelpCmd:
		return p.sendHelp(userID)
	case GetInfoCmd:
		return p.getInfo(userID)
	case GetStatCmd:
		return p.getStatistic(userID)
	case GetWhaleCmd:
		return p.getWhalesTransactions(userID)
	default:
		return p.tg.SendMessage(userID, msgUnknownCommand)
	}
}
func (p *Processor) getWhalesTransactions(userID int) error {
	data, err := p.wh.GetWhaleActionFromCache()
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for i := range data {
		buffer.WriteString(fmt.Sprintf("<b>Имя</b>: %s\n", data[i].WhaleName))
		buffer.WriteString(fmt.Sprintf("<b>Позиция</b>: %s\n", data[i].WhalePosition))
		buffer.WriteString(fmt.Sprintf("<b>Тип</b>: %s\n", data[i].Type))
		buffer.WriteString(fmt.Sprintf("<b>ETH</b>: %f\n", data[i].ValueETH))
		buffer.WriteString(fmt.Sprintf("<b>Дата</b>: %s\n", data[i].Date.Format("02.01.2006 15:04:05")))
		buffer.WriteString(fmt.Sprintf("<b>Общий баланс</b>: %f ETH\n", data[i].Balance))
		buffer.WriteString("\n")
	}

	return p.tg.SendMessage(userID, buffer.String())
}

func (p *Processor) getStatistic(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	data := p.leaderBoardService.GetStatistic(ctx)

	return p.tg.SendMessage(userID, data.ToTelegramMessage())
}
func (p *Processor) getInfo(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	data, err := p.leaderBoardService.GetLeader(ctx)
	if err != nil {
		return err
	}
	if len(data.Data.OtherPositionRetList) == 0 {
		return fmt.Errorf("user dont have bets")
	}
	rate := data.Data.OtherPositionRetList[0]
	pos := entity.NewPosition(rate.Symbol, rate.Amount, rate.EntryPrice, rate.Pnl, rate.Leverage)

	return p.tg.SendMessage(userID, pos.ToTelegramMessage())
}
func (p *Processor) sendHelp(userID int) error {
	return p.tg.SendMessage(userID, msgHelp)
}

func retrieveCommandName(text string) (string, error) {
	errMsg := fmt.Errorf("команда не найдена в сообщении")
	result := strings.Split(strings.TrimSpace(strings.TrimSuffix(text, "\n")), " ")
	if len(result) == 0 {
		return "", errMsg
	}

	if strings.HasPrefix(result[0], "/") == false {
		return "", errMsg
	}

	return result[0], nil
}
