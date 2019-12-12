// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"

	"github.com/leonardochaia/vendopunkto/internal/currency"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/internal/server"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

// Create a new server
func NewServer(globalLogger hclog.Logger) (*server.Server, error) {
	wire.Build(
		pluginmgr.PluginProviders,
		invoice.InvoiceProviders,
		server.ServerProviders,
		currency.CurrencyProviders,
		store.NewDB,
	)
	return &server.Server{}, nil
}
