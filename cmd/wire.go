// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"

	"github.com/leonardochaia/vendopunkto/invoice"
	"github.com/leonardochaia/vendopunkto/pluginmgr"
	"github.com/leonardochaia/vendopunkto/server"
	"github.com/leonardochaia/vendopunkto/store"
)

// Create a new server
func NewServer(globalLogger hclog.Logger) (*server.Server, error) {
	wire.Build(
		pluginmgr.NewManager,
		invoice.InvoiceProviders,
		server.NewServer,
		server.NewRouter,
		store.NewDB,
	)
	return &server.Server{}, nil
}
