package server

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var ServerProviders = wire.NewSet(NewServer, NewRouter)

func NewServer(
	router *VendoPunktoRouter,
	db *gorm.DB,
	globalLogger hclog.Logger,
	pluginManager *pluginmgr.Manager) (*Server, error) {

	server := &Server{
		logger:        globalLogger.Named("server"),
		router:        router,
		db:            db,
		pluginManager: pluginManager,
	}

	return server, nil
}
