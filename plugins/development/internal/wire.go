// This file uses wire to build all the depdendancies required

// +build wireinject

package development

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
)

func NewContainer(globalLogger hclog.Logger) (*Container, error) {
	wire.Build(
		clients.Providers,
		Providers,
	)
	return &Container{}, nil
}
