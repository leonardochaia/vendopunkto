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
	manager, err := currency.NewManager(db, globalLogger2)
	if err != nil {
		return nil, err
	}
	pluginmgrManager := pluginmgr.NewManager(globalLogger2, manager)
	invoiceManager, err := invoice.NewManager(db, pluginmgrManager, globalLogger2)
	if err != nil {
		return nil, err
	}
	handler := invoice.NewHandler(invoiceManager, globalLogger2)
	vendoPunktoRouter, err := server.NewRouter(handler, globalLogger2)
	if err != nil {
		return nil, err
	}
	pluginRouter, err := server.NewPluginRouter(handler, globalLogger2)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.NewServer(vendoPunktoRouter, pluginRouter, db, globalLogger2, pluginmgrManager)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}
