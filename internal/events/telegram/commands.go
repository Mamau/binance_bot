package telegram

import (
	"binance_bot/internal/entity"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	HelpCmd    = "/help"
	GetInfoCmd = "/get_info"
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
	default:
		return p.tg.SendMessage(userID, msgUnknownCommand)
	}
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
