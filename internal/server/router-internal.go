package server

import (
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
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
	db *pg.DB,
) (*InternalRouter, error) {

	var router InternalRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("internal-server")
	setupMiddlewares(router, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, db))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoice.InternalRoutes())
	})
	return &router, nil
}
