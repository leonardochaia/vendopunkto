package server

import (
	"bytes"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/currency"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
	"github.com/spf13/viper"
)

type VendoPunktoRouter interface {
	chi.Router
}

// Creates the chi Router and configures global paths
func NewRouter(
	invoices *invoice.Handler,
	currencies currency.Handler,
	globalLogger hclog.Logger,
	db *pg.DB,
) (*VendoPunktoRouter, error) {

	var router VendoPunktoRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("api-server")

	// add basic middlewares
	setupMiddlewares(router, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, db))

	// global routes
	router.Get("/info", GetVersion())

	// versions
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoices.Routes())
		r.Mount("/currencies", currencies.Routes())
	})

	return &router, nil
}

func setupMiddlewares(router chi.Router, logger hclog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.GetHead)

	// Log Requests
	if viper.GetBool("server.log_requests") {
		router.Use(newRequestLogger(logger))
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
	if viper.GetBool("server.profiler_enabled") && viper.GetString("server.profiler_path") != "" {
		logger.Debug("Profiler enabled on API", "path", viper.GetString("server.profiler_path"))
		router.Mount(viper.GetString("server.profiler_path"), middleware.Profiler())
	}
}
func newRequestLogger(parentLogger hclog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var requestID string
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				requestID = reqID.(string)
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			responseBuffer := new(bytes.Buffer)
			ww.Tee(responseBuffer)
			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			logger := parentLogger.Named("request").With(
				"status", ww.Status(),
				"duration", latency.String(),
				"remote", r.RemoteAddr,
				"url", r.RequestURI,
				"method", r.Method,
				"requestID", requestID)

			if ww.Status() >= 500 {
				logger.Error("Request errored",
					"response", responseBuffer.String())
			} else if ww.Status() >= 400 {
				logger.Warn("Bad request",
					"response", responseBuffer.String())
			} else if ww.Status() >= 200 {
				logger.Info("Request succees")
			}
		})
	}

}
