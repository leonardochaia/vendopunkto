package pluginwallet

import "github.com/go-chi/chi"

import "github.com/leonardochaia/vendopunkto/errors"

type Router interface {
	chi.Router
}

func NewRouter(handler *Handler) *Router {

	var router Router
	router = chi.NewRouter()

	router.Post("/payments/confirm", errors.WrapHandler(handler.confirmPayment))

	return &router
}
