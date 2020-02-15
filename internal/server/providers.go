package server

import (
	"github.com/go-pg/pg"
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/pluginwallet"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

// Providers wire providers for the server package
var Providers = wire.NewSet(
	NewServer,
	NewRouter,
	NewInternalRouter,
	NewInvoiceHandler,
	NewCurrencyHandler,
	NewConfigHandler,
	NewPluginHandler,
)

// NewServer creates the server
func NewServer(
	router *VendoPunktoRouter,
	internalRouter *InternalRouter,
	db *pg.DB,
	txBuilder store.TransactionBuilder,
	globalLogger hclog.Logger,
	pluginManager vendopunkto.PluginManager,
	walletPoller *pluginwallet.WalletPoller,
	invoiceTopic vendopunkto.InvoiceTopic,
	startupConf conf.Startup) (*Server, error) {

	server := &Server{
		logger:         globalLogger.Named("server"),
		router:         router,
		db:             db,
		txBuilder:      txBuilder,
		pluginManager:  pluginManager,
		internalRouter: internalRouter,
		startupConf:    startupConf,
		walletPoller:   walletPoller,
		invoiceTopic:   invoiceTopic,
	}

	return server, nil
}

// NewInvoiceHandler the handler for invoices
func NewInvoiceHandler(
	manager vendopunkto.InvoiceManager,
	globalLogger hclog.Logger,
	pluginMgr vendopunkto.PluginManager,
	topic vendopunkto.InvoiceTopic) *InvoiceHandler {
	return &InvoiceHandler{
		logger:    globalLogger.Named("invoice-handler"),
		manager:   manager,
		pluginMgr: pluginMgr,
		topic:     topic,
	}
}

// NewCurrencyHandler creates the currency handler
func NewCurrencyHandler(manager vendopunkto.PluginManager,
	currencyRepo vendopunkto.CurrencyRepository,
	runtime *conf.Runtime,
	globalLogger hclog.Logger) (*CurrencyHandler, error) {
	return &CurrencyHandler{
		logger:       globalLogger.Named("currency-handler"),
		plugins:      manager,
		currencyRepo: currencyRepo,
		runtimeConf:  runtime,
	}, nil
}

// NewConfigHandler creates the config handler
func NewConfigHandler(runtime *conf.Runtime,
	globalLogger hclog.Logger) (*ConfigHandler, error) {
	return &ConfigHandler{
		logger:  globalLogger.Named("currency-handler"),
		runtime: runtime,
	}, nil
}

// NewPluginHandler creates the plugin handler
func NewPluginHandler(plugins vendopunkto.PluginManager,
	globalLogger hclog.Logger) (*PluginHandler, error) {
	return &PluginHandler{
		logger:  globalLogger.Named("currency-handler"),
		plugins: plugins,
	}, nil
}
