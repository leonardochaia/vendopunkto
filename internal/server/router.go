package server

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

// VendoPunktoRouter is the public routerplugin
type VendoPunktoRouter interface {
	chi.Router
}

// NewRouter Creates the chi Router and configures global paths
func NewRouter(
	invoices *invoice.Handler,
	globalLogger hclog.Logger,
	txBuilder store.TransactionBuilder,
	startupConf conf.Startup,
) (*VendoPunktoRouter, error) {

	var router VendoPunktoRouter
	router = chi.NewRouter()

	logger := globalLogger.Named("api-server")

	// add basic middlewares
	setupMiddlewares(router, startupConf, logger)

	// tx per request
	router.Use(store.NewTxPerRequestMiddleware(globalLogger, txBuilder))

	// versions
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/info", GetVersion())
		r.Mount("/invoices", invoices.Routes())
	})

	serveSPA(router, "/", "spa/dist/vendopunkto")

	return &router, nil
}

func setupMiddlewares(
	router chi.Router,
	startupConf conf.Startup,
	logger hclog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.GetHead)

	// Log Requests
	if startupConf.Server.LogRequests {
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
	if startupConf.Server.ProfilerEnabled && startupConf.Server.ProfilerPath != "" {
		logger.Debug("Profiler enabled on API", "path", startupConf.Server.ProfilerPath)
		router.Mount(startupConf.Server.ProfilerPath, middleware.Profiler())
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

// serveSPA is serving static files
func serveSPA(r chi.Router, public string, static string) {

	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", 301).ServeHTTP)
		public += "/"
	}

	r.Get(public+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}
