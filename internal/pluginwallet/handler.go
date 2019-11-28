package pluginwallet

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
)

type Handler struct {
	logger hclog.Logger
}

func (handler *Handler) confirmPayment(w http.ResponseWriter, r *http.Request) *errors.APIError {

	render.JSON(w, r, "TODO")
	return nil
}
