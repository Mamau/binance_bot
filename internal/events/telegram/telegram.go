package telegram

import (
	telegramClient "binance_bot/internal/clients/telegram"
	"binance_bot/internal/events"
	"binance_bot/internal/usecase"
	"errors"
	"fmt"
)

type Processor struct {
	tg                 *telegramClient.Client
	leaderBoardService *usecase.LeaderBoard
	offset             int
}
type Meta struct {
	UserId    int
	Username  string
	FirstName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(client *telegramClient.Client, lbs *usecase.LeaderBoard) *Processor {
	return &Processor{
		tg:                 client,
		offset:             0,
		leaderBoardService: lbs,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("cant get events %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, v := range updates {
		res = append(res, event(v))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("cant process message %w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("cant process message %w", err)
	}
	if err := p.doCmd(event.Text, meta.UserId, meta.Username, meta.FirstName); err != nil {
		if err := p.tg.SendMessage(meta.UserId, err.Error()); err != nil {
			return err
		}
		return fmt.Errorf("cant process message %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("error while type meta assertion %w", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegramClient.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			UserId:    upd.Message.From.ID,
			Username:  upd.Message.From.Username,
			FirstName: upd.Message.From.FirstName,
		}
	}
	return res
}

func fetchText(upd telegramClient.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegramClient.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
