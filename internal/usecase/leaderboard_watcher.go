package usecase

import (
	"binance_bot/pkg/log"
	"context"
	"time"
)

type LeaderBoardWatcher struct {
	service *LeaderBoard
	logger  log.Logger
	done    chan struct{}
}

func NewLeaderBoardWatcher(s *LeaderBoard, l log.Logger) *LeaderBoardWatcher {
	return &LeaderBoardWatcher{
		service: s,
		logger:  l,
	}
}

func (lbw *LeaderBoardWatcher) Run() {
	defer close(lbw.done)
	ticker := time.NewTicker(time.Second * 3)

	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				data, err := lbw.service.GetLeader(ctx)
				if err != nil {
					lbw.logger.Err(err).Msg("error while get leader by daemon")
					return
				}
				lbw.service.NotifyAboutBet(data)
			}()
		}
	}
}
func (lbw *LeaderBoardWatcher) Terminate(ctx context.Context) error {
	lbw.logger.Info().Msg("terminating observe post leader bot watcher")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-lbw.done:
		return nil
	}
}
