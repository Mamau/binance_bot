package timeout

import (
	"time"
)

type Option func(*options)

type options struct {
	timeout time.Duration
}

func Timeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}
