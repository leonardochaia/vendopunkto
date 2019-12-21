package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// WalletPluginHandler is the Wallet Plugin HTTP API implementation
// which runs on the plugin server and will be called
// by vendopunkto
// It exposes the wallet interface provided by the
// plugin developer into HTTP endpoints
type WalletPluginHandler struct {
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

	handler := &WalletPluginHandler{
		wallet:       plugin,
		serverPlugin: serverPlugin,
	}

	router.Post(GenerateAddressWalletEndpoint, errors.WrapHandler(handler.generateAddress))
	router.Get(WalletInfoEndpoint, errors.WrapHandler(handler.getWalletInfo))

	return router
}

func (handler *WalletPluginHandler) generateAddress(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "plugin.wallet.generateAddress"
	var params = new(CoinWalletAddressParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	res, err := handler.wallet.GenerateNewAddress(params.InvoiceID)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, &CoinWalletAddressResponse{
		Address: res,
	})

	return nil
}

func (handler *WalletPluginHandler) getWalletInfo(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "plugin.wallet.getWalletInfo"
	res, err := handler.wallet.GetWalletInfo()

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, res)
	return nil
}
