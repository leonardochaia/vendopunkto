package plugin

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
)

// WalletPluginCurrency provides metadata for WalletPlugin currencies
type WalletPluginCurrency struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	QRCodeTemplate string `json:"qrCodeTemplate"`
}

// WalletPluginInfo provides metadata for the WalletPlugin
type WalletPluginInfo struct {
	Currency WalletPluginCurrency `json:"currency"`
}

// WalletPluginIncomingTransferParams are the filters for returning transfers
type WalletPluginIncomingTransferParams struct {
	MinBlockHeight uint64 `json:"minBlockHeight"`
}

// WalletPluginIncomingTransferResult is a transfer representation
type WalletPluginIncomingTransferResult struct {
	TxHash        string          `json:"txHash"`
	Address       string          `json:"address"`
	BlockHeight   uint64          `json:"blockHeight"`
	Confirmations uint64          `json:"confirmations"`
	Amount        decimal.Decimal `json:"amount"`
}

// WalletPlugin must be implemented for a currency to be supported by VendoPunkto
type WalletPlugin interface {
	VendoPunktoPlugin
	// GetWalletInfo returns information about the currency of this wallet
	GetWalletInfo() (WalletPluginInfo, error)
	// GenerateNewAddress will be called everytime a new payment for this
	// currency is requested. It must generate a new addres everytime. This is
	// by design, given that received payments will be matched by address.
	GenerateNewAddress(invoiceID string) (string, error)
	// GetIncomingTransfers returns the list of input transfers, confirmed and
	// on the mempool. Filters should be applied
	GetIncomingTransfers(params WalletPluginIncomingTransferParams) ([]WalletPluginIncomingTransferResult, error)
}

// walletServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type walletServerPlugin struct {
	Impl WalletPlugin
}

// BuildWalletPlugin needs to be called by implementors with their WalletPlugin
// implementation.
// Afterwards, you need to call server.AddPlugin with the resulting ServerPlugin
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
	// WalletMainEndpoint is root wallet plugin path
	WalletMainEndpoint = "/vp/wallet"
	// GenerateAddressWalletEndpoint the suffix for address generation
	GenerateAddressWalletEndpoint = "/address"
	// GetIncomingTransfersWalletEndpoint returns incoming transfers
	GetIncomingTransfersWalletEndpoint = "/incoming-transfers"
	// WalletInfoEndpoint the suffix for info
	WalletInfoEndpoint = "/info"
)

// BuildQRCode generates the string for a QR code based on
// the template. If the WalletPluginInfo has no template, BIP21 will be used
func (info WalletPluginInfo) BuildQRCode(
	address string,
	amount decimal.Decimal) (string, error) {

	if info.Currency.QRCodeTemplate == "" {
		// default to bip21
		info.Currency.QRCodeTemplate = fmt.Sprintf(
			"%s:{{.Address}}?amount={{.AmountFormatted}}",
			strings.ToLower(info.Currency.Name))
	}

	data := struct {
		Address         string
		Amount          decimal.Decimal
		AmountFormatted string
	}{
		Address:         address,
		Amount:          amount,
		AmountFormatted: amount.String(),
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
