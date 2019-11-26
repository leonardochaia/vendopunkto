// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"

	"github.com/leonardochaia/vendopunkto/invoice"
	"github.com/leonardochaia/vendopunkto/monero"
	"github.com/leonardochaia/vendopunkto/server"
	"github.com/leonardochaia/vendopunkto/store"
)

// Create a new server
func NewServer(globalLogger hclog.Logger) (*server.Server, error) {
	wire.Build(
		invoice.InvoiceProviders,
		server.NewServer,
		server.NewRouter,
		store.NewDB,
		monero.CreateMoneroClient)
	return &server.Server{}, nil
}
