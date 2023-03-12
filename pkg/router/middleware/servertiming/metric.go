package servertiming

import (
	orig "github.com/mitchellh/go-server-timing"
)

type Metric struct {
	*orig.Metric
}
