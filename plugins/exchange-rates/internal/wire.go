// This file uses wire to build all the depdendancies required

// +build wireinject

package rates

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
)

func NewContainer(globalLogger hclog.Logger) (*Container, error) {
	wire.Build(
		Providers,
	)
	return &Container{}, nil
}
