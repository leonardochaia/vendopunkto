package server

import (
	"github.com/go-chi/chi"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

// InternalRouter is the router for the internal api
type InternalRouter interface {
	chi.Router
}

// NewInternalRouter creates the router
func NewInternalRouter(
	invoice *InvoiceHandler,
	currencies *CurrencyHandler,
	config *ConfigHandler,
	plugin *PluginHandler,
	globalLogger hclog.Logger,
	txBuilder store.TransactionBuilder,
	startupConf conf.Startup,
) (*InternalRouter, error) {

	var router InternalRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("internal-server")
	setupMiddlewares(router, startupConf, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, txBuilder))

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/invoices", invoice.InternalRoutes())
		r.Mount("/currencies", currencies.InternalRoutes())
		r.Mount("/config", config.InternalRoutes())
		r.Mount("/plugins", plugin.InternalRoutes())
	})

	serveSPA(router, "/", "spa/dist/admin")

	return &router, nil
}
