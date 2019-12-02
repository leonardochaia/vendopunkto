package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// Handler is the Wallet Plugin HTTP API implementation
// which runs on the plugin server and will be called
// by vendopunkto
// It exposes the wallet interface provided by the
// plugin developer into HTTP endpoints
type Handler struct {
	wallet       WalletPlugin
	serverPlugin ServerPlugin
}

type CoinWalletAddressResponse struct {
	Address string `json:"address"`
}

type CoinWalletAddressParams struct {
	InvoiceID string `json:"invoiceId"`
}

type ActivatePluginParams struct {
	HostAddress string `json:"hostAddress"`
}

func NewWalletHandler(plugin WalletPlugin, serverPlugin ServerPlugin) *chi.Mux {
	router := chi.NewRouter()

	handler := &Handler{
		wallet:       plugin,
		serverPlugin: serverPlugin,
	}

	router.Post(GenerateAddressWalletEndpoint, errors.WrapHandler(handler.generateAddress))
	router.Post(ActivatePluginEndpoint, errors.WrapHandler(handler.activatePlugin))

	return router
}

func (handler *Handler) generateAddress(w http.ResponseWriter, r *http.Request) *errors.APIError {
	var params = new(CoinWalletAddressParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	res, err := handler.wallet.GenerateNewAddress(params.InvoiceID)
	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, &CoinWalletAddressResponse{
		Address: res,
	})

	return nil
}

func (handler *Handler) activatePlugin(w http.ResponseWriter, r *http.Request) *errors.APIError {
	var params = new(ActivatePluginParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	err := handler.serverPlugin.initializePlugin(params.HostAddress)
	if err != nil {
		return errors.InternalServerError(err)
	}

	res, err := handler.wallet.GetPluginInfo()

	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, res)
	return nil
}
