// Package plugin contains utilities for building a VendoPunkto plugin in Go
// A plugin is basically a binary which runs an HTTP server that VendoPunkto
// consumes. This package aims to help doing so.
package plugin

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/errors"
)

// Server is the plugin server running on every plugin
type Server struct {
	plugins        []ServerPlugin
	pluginInfos    []PluginInfo
	started        bool
	internalClient clients.InternalClient
	Logger         hclog.Logger
	http           clients.HTTP
}

// ServerPlugin wraps the plugin implementation.
// It configures the router and orchestrates the underlying calls
// to the actual implementation
type ServerPlugin interface {
	initializeRouter(router *chi.Mux) error
	GetPluginImpl() (VendoPunktoPlugin, error)
}

// NewServer creates the plugin server which must be run by every plugin
func NewServer(
	logger hclog.Logger,
	http clients.HTTP) *Server {
	return &Server{
		plugins:     []ServerPlugin{},
		pluginInfos: []PluginInfo{},
		started:     false,
		Logger:      logger,
		http:        http,
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

	s.pluginInfos = []PluginInfo{}

	s.initializeAllPlugins(router)

	router.Post(ActivatePluginEndpoint, errors.WrapHandler(s.activatePluginHandler))

	return http.ListenAndServe(addr, router)
}

// GetInternalClient returns the clients.InternalClient that you can use to
// consume the VendoPunkto internal API
func (s *Server) GetInternalClient() (clients.InternalClient, error) {

	if s.internalClient == nil {
		err := fmt.Errorf("Server has not been activated or activation has failed")
		s.Logger.Error("Server had no internal client", "error", err)
		return nil, err
	}

	return s.internalClient, nil
}

// initializeAllPlugins loops through registered plugins and
// initializes their router and stores references to implementations
func (s *Server) initializeAllPlugins(router chi.Router) error {
	for _, value := range s.plugins {

		impl, err := value.GetPluginImpl()
		if err != nil {
			return err
		}

		info, err := impl.GetPluginInfo()

		if err != nil {
			return err
		}

		pRouter := chi.NewRouter()
		router.Mount(info.GetAddress(), pRouter)
		err = value.initializeRouter(pRouter)

		if err != nil {
			return err
		}

		s.pluginInfos = append(s.pluginInfos, info)

		s.Logger.Info("Routed plugin",
			"id", info.ID,
			"name", info.Name,
			"address", info.GetAddress(),
		)
	}

	return nil
}

func (s *Server) activatePluginHandler(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "plugin.base.activate"

	var params = new(ActivatePluginParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	client, err := clients.NewInternalClient(params.HostAddress, s.http)

	if err != nil {
		s.Logger.Error("Failed to create internal client", "error", err)
	}

	s.Logger.Info("Activating plugin", "host", params.HostAddress)
	s.internalClient = client

	render.JSON(w, r, s.pluginInfos)
	return nil
}
