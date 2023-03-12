package application

import "context"

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type Daemon interface {
	Run()
	Terminate(ctx context.Context) error
}
