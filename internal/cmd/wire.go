// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"

	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/internal/pluginwallet"
	"github.com/leonardochaia/vendopunkto/internal/server"
	"github.com/leonardochaia/vendopunkto/internal/store"
	"github.com/leonardochaia/vendopunkto/internal/store/repositories"
)

// Create a new server
func NewServer(globalLogger hclog.Logger) (*server.Server, error) {
	wire.Build(
		conf.Providers,
		pluginmgr.Providers,
		invoice.Providers,
		server.Providers,
		clients.Providers,
		repositories.Providers,
		store.Providers,
		pluginwallet.NewPoller,
	)
	return &server.Server{}, nil
}
