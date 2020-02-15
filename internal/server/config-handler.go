package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/conf"
)

// ConfigHandler for configurations
type ConfigHandler struct {
	logger  hclog.Logger
	runtime *conf.Runtime
}

// InternalRoutes creates a router for the internal API
func (handler *ConfigHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getConfiguration))

	return router
}

func (handler *ConfigHandler) getConfiguration(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.configHandler.getConfiguration"

	render.JSON(w, r, handler.runtime.AllSettings())
	return nil
}
