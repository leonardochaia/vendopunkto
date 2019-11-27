package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// Handler is the Wallet Plugin HTTP API implementation
// which runs on the plugin server
type Handler struct {
	wallet WalletPlugin
}

type CoinWalletAddressResponse struct {
	Address string `json:"address"`
}

type CoinWalletAddressParams struct {
	InvoiceID string `json:"invoiceId"`
}

func NewWalletHandler(plugin WalletPlugin) *chi.Mux {
	router := chi.NewRouter()

	handler := &Handler{
		wallet: plugin,
	}

	router.Post(GenerateAddressWalletEndpoint, errors.WrapHandler(handler.generateWalletHandler))
	router.Get(PluginInfoEndpoint, errors.WrapHandler(handler.getPluginInfo))

	return router
}

func (handler *Handler) generateWalletHandler(w http.ResponseWriter, r *http.Request) *errors.APIError {
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

func (handler *Handler) getPluginInfo(w http.ResponseWriter, r *http.Request) *errors.APIError {
	res, err := handler.wallet.GetPluginInfo()

	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, res)
	return nil
}
