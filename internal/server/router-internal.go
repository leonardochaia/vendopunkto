package server

import (
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/currency"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

// InternalRouter is the router for the internal api
type InternalRouter interface {
	chi.Router
}

// NewInternalRouter creates the router
func NewInternalRouter(
	invoice *invoice.Handler,
	globalLogger hclog.Logger,
	currencies currency.Handler,
	db *pg.DB,
	startupConf conf.Startup,
) (*InternalRouter, error) {

	var router InternalRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("internal-server")
	setupMiddlewares(router, startupConf, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, db))

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/invoices", invoice.InternalRoutes())
		r.Mount("/currencies", currencies.InternalRoutes())
	})
	return &router, nil
}
