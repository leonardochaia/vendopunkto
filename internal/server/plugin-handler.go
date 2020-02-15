package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
)

// PluginHandler for plugins
type PluginHandler struct {
	logger  hclog.Logger
	plugins vendopunkto.PluginManager
}

// InternalRoutes creates a router for the internal API
func (handler *PluginHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getPlugins))

	return router
}

func (handler *PluginHandler) getPlugins(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.pluginHandler.getPlugins"

	plugins, err := handler.plugins.GetAllPlugins()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, plugins)
	return nil
}
