package server

import (
	"github.com/go-chi/chi"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
)

type PluginRouter interface {
	chi.Router
}

func NewPluginRouter(
	invoice *invoice.Handler,
	globalLogger hclog.Logger) (*PluginRouter, error) {

	var router PluginRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("plugin-server")
	setupMiddlewares(router, logger)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoice.InternalRoutes())
	})
	return &router, nil
}
