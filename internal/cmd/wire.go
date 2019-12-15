// This file uses wire to build all the depdendancies required

// +build wireinject

package cmd

import (
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"

	"github.com/leonardochaia/vendopunkto/internal/currency"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/internal/server"
	"github.com/leonardochaia/vendopunkto/internal/store"
	"github.com/leonardochaia/vendopunkto/internal/store/repositories"
)

// Create a new server
func NewServer(globalLogger hclog.Logger) (*server.Server, error) {
	wire.Build(
		pluginmgr.PluginProviders,
		invoice.InvoiceProviders,
		server.ServerProviders,
		currency.CurrencyProviders,
		NewHttpClient,
		repositories.Providers,
		store.NewDB,
	)
	return &server.Server{}, nil
}

func NewHttpClient() http.Client {
	return http.Client{
		Timeout: 15 * time.Second,
	}
}
