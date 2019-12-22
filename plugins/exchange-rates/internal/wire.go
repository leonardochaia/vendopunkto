// This file uses wire to build all the depdendancies required

// +build wireinject

package rates

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
)

func NewContainer(globalLogger hclog.Logger) (*Container, error) {
	wire.Build(
		Providers,
		clients.Providers,
	)
	return &Container{}, nil
}
