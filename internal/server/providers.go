package server

import (
	"github.com/go-pg/pg"
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/internal/pluginwallet"
)

// ServerProviders wire providers for the server package
var ServerProviders = wire.NewSet(NewServer, NewRouter, NewInternalRouter)

// NewServer creates the server
func NewServer(
	router *VendoPunktoRouter,
	internalRouter *InternalRouter,
	db *pg.DB,
	globalLogger hclog.Logger,
	pluginManager *pluginmgr.Manager,
	walletPoller *pluginwallet.WalletPoller,
	invoiceTopic invoice.Topic,
	startupConf conf.Startup) (*Server, error) {

	server := &Server{
		logger:         globalLogger.Named("server"),
		router:         router,
		db:             db,
		pluginManager:  pluginManager,
		internalRouter: internalRouter,
		startupConf:    startupConf,
		walletPoller:   walletPoller,
		invoiceTopic:   invoiceTopic,
	}

	return server, nil
}
