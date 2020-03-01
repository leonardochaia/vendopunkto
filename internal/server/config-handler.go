package server

import (
	"encoding/json"
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
	manager conf.Manager
}

// InternalRoutes creates a router for the internal API
func (handler *ConfigHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getConfiguration))
	router.Post("/", errors.WrapHandler(handler.updateConfiguration))

	return router
}

func (handler *ConfigHandler) getConfiguration(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.configHandler.getConfiguration"

	render.JSON(w, r, handler.runtime.AllSettings())
	return nil
}

func (handler *ConfigHandler) updateConfiguration(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.configHandler.updateConfiguration"

	// we don't know the type of value. it could be a string or a []string
	var params = make(map[string]json.RawMessage)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	for key, rawMessage := range params {
		if !handler.runtime.IsKnownKey(key) {
			err := errors.Errorf("Provided key is not a known runtime config key: %s", key)
			return errors.E(op, errors.Parameters, err)
		}

		var arrValue []string
		// try to unmarshal into a []string first
		if err := json.Unmarshal(rawMessage, &arrValue); err != nil {
			// if it's not an array, try with string.
			var strValue string
			if err := json.Unmarshal(rawMessage, &strValue); err != nil {
				err = errors.Errorf("Attempted to parse JSON as string or []string and failed. Key: %s, error: %s",
					key, err)
				return errors.E(op, errors.Parameters, err)
			}

			// save the string value
			err = handler.manager.SaveConfiguration(key, strValue)
			if err != nil {
				return errors.E(op, errors.Parameters, err)
			}
			continue
		}

		// save the array value
		err := handler.manager.SaveConfiguration(key, arrValue)
		if err != nil {
			return errors.E(op, errors.Parameters, err)
		}
		// TODO: validate that the provided values make sense.
	}

	render.JSON(w, r, handler.runtime.AllSettings())

	return nil
}
