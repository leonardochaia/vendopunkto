package server

import (
	"github.com/go-pg/pg"
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var ServerProviders = wire.NewSet(NewServer, NewRouter, NewPluginRouter)

func NewServer(
	router *VendoPunktoRouter,
	pluginRouter *PluginRouter,
	db *pg.DB,
	globalLogger hclog.Logger,
	pluginManager *pluginmgr.Manager) (*Server, error) {

	server := &Server{
		logger:        globalLogger.Named("server"),
		router:        router,
		db:            db,
		pluginManager: pluginManager,
		pluginRouter:  pluginRouter,
	}

	return server, nil
}
