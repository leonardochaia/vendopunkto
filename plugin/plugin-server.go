package plugin

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// Server is the plugin server running on every plugin
type Server struct {
	plugins []ServerPlugin
	started bool
}

// ServerPlugin wraps the plugin implementation.
// It configures the router and orchestrates the underlying calls
// to the actual implementation
type ServerPlugin interface {
	initializeRouter(router *chi.Mux)
}

// NewServer creates the plugin server which must be run by every plugin
func NewServer() *Server {
	return &Server{
		plugins: []ServerPlugin{},
	}
}

// AddPlugin adds a plugin to the server
// This must be called before the server is started
func (s *Server) AddPlugin(plugin ServerPlugin) error {
	if s.started {
		return fmt.Errorf("Plugins must be added before the server is started")
	}

	s.plugins = append(s.plugins, plugin)

	return nil
}

// Start initializes the router and starts serving on the provided net address
func (s *Server) Start(addr string) error {
	if s.started {
		return fmt.Errorf("Server has already been started")
	}
	router := chi.NewRouter()

	for _, value := range s.plugins {
		value.initializeRouter(router)
	}

	return http.ListenAndServe(addr, router)
}
