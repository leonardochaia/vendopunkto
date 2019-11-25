// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"

	"github.com/leonardochaia/vendopunkto/invoice"
	"github.com/leonardochaia/vendopunkto/monero"
	"github.com/leonardochaia/vendopunkto/server"
	"github.com/leonardochaia/vendopunkto/store"
)

// Create a new server
func NewServer() (*server.Server, error) {
	wire.Build(
		server.NewServer,
		server.NewRouter,
		invoice.NewManager,
		invoice.NewHandler,
		store.NewDB,
		monero.CreateMoneroClient)
	return &server.Server{}, nil
}
