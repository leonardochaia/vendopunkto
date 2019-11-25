package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/invoice"
)

// Server is the API web server
type Server struct {
	logger *zap.SugaredLogger
	router chi.Router
	server *http.Server
	db     *gorm.DB
}

func NewServer(invoices *invoice.Handler, db *gorm.DB) (*Server, error) {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Log Requests
	if config.GetBool("server.log_requests") {
		router.Use(RequestLogger)
	}

	// CORS Config
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/invoices", invoices.Routes())
	})

	server := &Server{
		logger: zap.S().With("package", "server"),
		router: router,
		db:     db,
	}

	return server, nil
}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(config.GetString("server.host"), config.GetString("server.port")),
		Handler: s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			s.logger.Fatalw("API Listen error", "error", err, "address", s.server.Addr)
		}
	}()
	s.logger.Infow("API Listening", "address", s.server.Addr)

	// Enable profiler
	if config.GetBool("server.profiler_enabled") && config.GetString("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", config.GetString("server.profiler_path"))
		s.router.Mount(config.GetString("server.profiler_path"), middleware.Profiler())
	}

	return nil

}

func (s *Server) Close() {
	s.db.Close()
}

// RenderOrErrInternal will render whatever you pass it (assuming it has Renderer) or prints an internal error
func RenderOrErrInternal(w http.ResponseWriter, r *http.Request, d render.Renderer) {
	if err := render.Render(w, r, d); err != nil {
		render.Render(w, r, errors.ErrInternal(err))
		return
	}
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		fields := []zapcore.Field{
			zap.Int("status", ww.Status()),
			zap.Duration("took", latency),
			zap.String("remote", r.RemoteAddr),
			zap.String("request", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("package", "server.request"),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request-id", requestID))
		}
		zap.L().Info("API Request", fields...)
	})
}
