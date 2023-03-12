//go:build wireinject
// +build wireinject

package main

import (
	"binance_bot/internal/serviceprovider"

	"github.com/google/wire"

	"binance_bot/pkg/application"
)

func newApp() (*application.App, func(), error) {
	panic(wire.Build(
		serviceprovider.ProviderSet,
		createApp,
	))
}
