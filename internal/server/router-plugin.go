package server

import (
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

type PluginRouter interface {
	chi.Router
}

func NewPluginRouter(
	invoice *invoice.Handler,
	globalLogger hclog.Logger,
	db *pg.DB,
) (*PluginRouter, error) {

	var router PluginRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("plugin-server")
	setupMiddlewares(router, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, db))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoice.InternalRoutes())
	})
	return &router, nil
}
