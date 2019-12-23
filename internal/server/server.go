package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

// Server is the API web server
type Server struct {
	logger         hclog.Logger
	router         *VendoPunktoRouter
	internalRouter *InternalRouter
	server         *http.Server
	db             *pg.DB
	pluginManager  *pluginmgr.Manager
	startupConf    conf.Startup
}

// startInternalServer serves the internal server on the configured listener.
func (s *Server) startInternalServer() error {

	addr := net.JoinHostPort(
		s.startupConf.Server.Internal.Host,
		s.startupConf.Server.Internal.Port,
	)
	server := &http.Server{
		Addr:    addr,
		Handler: *s.internalRouter,
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
	s.logger.Info("Internal Server Listening", "address", server.Addr)

	return nil
}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	s.pluginManager.LoadPlugins()

	err := s.startInternalServer()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(
		s.startupConf.Server.Public.Host,
		s.startupConf.Server.Public.Port,
	)
	s.server = &http.Server{
		Addr:    addr,
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
	s.logger.Info("Public Server Listening", "address", s.server.Addr)

	return nil
}

// Close finalizes any open resources
func (s *Server) Close() {
	s.db.Close()
}
