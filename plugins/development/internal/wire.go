// This file uses wire to build all the dependencies required

// +build wireinject

package development

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
