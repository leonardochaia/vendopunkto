package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/spf13/viper"
)

// Server is the API web server
type Server struct {
	logger        hclog.Logger
	router        *VendoPunktoRouter
	pluginRouter  *PluginRouter
	server        *http.Server
	db            *pg.DB
	pluginManager *pluginmgr.Manager
}

// StartPluginServer serves the plugin server on the configured listener.
func (s *Server) startPluginServer() error {
	s.pluginManager.LoadPlugins()

	server := &http.Server{
		Addr:    net.JoinHostPort(viper.GetString("plugins.server.host"), viper.GetString("plugins.server.port")),
		Handler: *s.pluginRouter,
	}

	// Listen
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", server.Addr, err)
	}

	go func() {
		if err = server.Serve(listener); err != nil {
			s.logger.Error("API Listen error", "error", err, "address", server.Addr)
			os.Exit(1)
		}
	}()
	s.logger.Info("Plugin Server Listening", "address", server.Addr)

	return nil
}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	err := s.startPluginServer()
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:    net.JoinHostPort(viper.GetString("server.host"), viper.GetString("server.port")),
		Handler: *s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			s.logger.Error("API Listen error", "error", err, "address", s.server.Addr)
			os.Exit(1)
		}
	}()
	s.logger.Info("API Server Listening", "address", s.server.Addr)

	return nil
}

func (s *Server) Close() {
	s.db.Close()
}
