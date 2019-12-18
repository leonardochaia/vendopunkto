package plugin

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/leonardochaia/vendopunkto/unit"
)

type WalletPluginCurrency struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	QRCodeTemplate string `json:"qrCodeTemplate"`
}

type WalletPluginInfo struct {
	Currency WalletPluginCurrency `json:"currency"`
}

// WalletPlugin must be implemented for a Coin to be supported by vendopunkto
type WalletPlugin interface {
	VendoPunktoPlugin
	GetWalletInfo() (WalletPluginInfo, error)
	GenerateNewAddress(invoiceID string) (string, error)
}

// walletServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type walletServerPlugin struct {
	Impl WalletPlugin
}

func BuildWalletPlugin(impl WalletPlugin) ServerPlugin {
	return &walletServerPlugin{
		Impl: impl,
	}
}

func (serverPlugin *walletServerPlugin) initializeRouter(router *chi.Mux) error {
	handler := NewWalletHandler(serverPlugin.Impl, serverPlugin)

	router.Mount(WalletMainEndpoint, handler)
	return nil
}

func (serverPlugin *walletServerPlugin) GetPluginImpl() (VendoPunktoPlugin, error) {
	return serverPlugin.Impl, nil
}

const (
	WalletMainEndpoint            = "/vp/wallet"
	GenerateAddressWalletEndpoint = "/address"
	WalletInfoEndpoint            = "/info"
)

func (info WalletPluginInfo) BuildQRCode(
	address string,
	amount unit.AtomicUnit) (string, error) {

	if info.Currency.QRCodeTemplate == "" {
		// default to bip21
		info.Currency.QRCodeTemplate = fmt.Sprintf(
			"%s:{{.Address}}?amount={{.AmountFormatted}}",
			strings.ToLower(info.Currency.Name))
	}

	data := struct {
		Address         string
		Amount          unit.AtomicUnit
		AmountFormatted string
	}{
		Address:         address,
		Amount:          amount,
		AmountFormatted: amount.Formatted(),
	}

	t, err := template.New("bip21").Parse(info.Currency.QRCodeTemplate)
	if err != nil {
		return "", err
	}

	var qr bytes.Buffer
	err = t.Execute(&qr, data)
	if err != nil {
		return "", err
	}

	return qr.String(), nil
}
