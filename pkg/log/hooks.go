package log

import (
	"time"

	"github.com/rs/zerolog"
)

type timeHook struct{}

func (ts timeHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Time(zerolog.TimestampFieldName, time.Now())
}
