package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server is the API web server
type Server struct {
	logger *zap.SugaredLogger
	router *VendoPunktoRouter
	server *http.Server
	db     *gorm.DB
}

func NewServer(router *VendoPunktoRouter, db *gorm.DB) (*Server, error) {

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
		Handler: *s.router,
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

	return nil

}

func (s *Server) Close() {
	s.db.Close()
}

func requestLogger(next http.Handler) http.Handler {
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
