// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/currency"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/internal/server"
	"github.com/leonardochaia/vendopunkto/internal/store"
	"github.com/leonardochaia/vendopunkto/internal/store/repositories"
	"net/http"
	"time"
)

import (
	_ "net/http/pprof"
)

// Injectors from wire.go:

func NewServer(globalLogger2 hclog.Logger) (*server.Server, error) {
	db, err := store.NewDB()
	if err != nil {
		return nil, err
	}
	invoiceRepository := repositories.NewPostgresInvoiceRepository(db)
	client := NewHttpClient()
	manager := pluginmgr.NewManager(globalLogger2, client)
	invoiceManager, err := invoice.NewManager(invoiceRepository, manager, globalLogger2)
	if err != nil {
		return nil, err
	}
	handler := invoice.NewHandler(invoiceManager, globalLogger2)
	currencyHandler, err := currency.NewHandler(manager, globalLogger2)
	if err != nil {
		return nil, err
	}
	vendoPunktoRouter, err := server.NewRouter(handler, currencyHandler, globalLogger2)
	if err != nil {
		return nil, err
	}
	pluginRouter, err := server.NewPluginRouter(handler, globalLogger2)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.NewServer(vendoPunktoRouter, pluginRouter, db, globalLogger2, manager)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

// wire.go:

func NewHttpClient() http.Client {
	return http.Client{
		Timeout: 15 * time.Second,
	}
}
