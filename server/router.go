package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/invoice"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
)

type VendoPunktoRouter interface {
	chi.Router
}

// Creates the chi Router and configures global paths
func NewRouter(invoices *invoice.Handler) (*VendoPunktoRouter, error) {

	var router VendoPunktoRouter
	router = chi.NewRouter()
	setupMiddlewares(router)

	router.Get("/info", GetVersion())

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoices.Routes())
	})
	return &router, nil
}

func setupMiddlewares(router VendoPunktoRouter) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Log Requests
	if config.GetBool("server.log_requests") {
		router.Use(RequestLogger)
	}

	// CORS Config
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead, http.MethodOptions,
			http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodDelete, http.MethodPatch,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// Enable profiler
	if config.GetBool("server.profiler_enabled") && config.GetString("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", config.GetString("server.profiler_path"))
		router.Mount(config.GetString("server.profiler_path"), middleware.Profiler())
	}
}
