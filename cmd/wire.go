// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"github.com/google/wire"
	config "github.com/spf13/viper"

	"github.com/leonardochaia/vendopunkto/invoice"
	"github.com/leonardochaia/vendopunkto/monero"
	"github.com/leonardochaia/vendopunkto/server"
	"github.com/leonardochaia/vendopunkto/store/postgres"
)

// Create a new server
func NewServer() (*server.Server, error) {
	wire.Build(server.NewServer, invoice.NewManager, invoice.NewStore,
		invoice.NewHandler,
		NewDBClient, monero.CreateMoneroClient)
	return &server.Server{}, nil
}

func NewDBClient() (*postgres.Client, error) {
	var pgClient *postgres.Client
	var err error
	switch config.GetString("storage.type") {
	case "postgres":
		pgClient, err = postgres.New()
	}
	if err != nil {
		return nil, err
		// logger.Fatalw("Database Error", "error", err)
	}
	return pgClient, nil
}
