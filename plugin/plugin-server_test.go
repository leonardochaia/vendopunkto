package plugin

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/testutils"
)

type testWalletPlugin struct {
	symbol         string
	name           string
	qrCodeTemplate string
}

func (p testWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	return p.symbol + "-test", nil
}

func (p testWalletPlugin) GetWalletInfo() (WalletPluginInfo, error) {
	return WalletPluginInfo{
		Currency: WalletPluginCurrency{
			Name:           p.name,
			Symbol:         p.symbol,
			QRCodeTemplate: p.qrCodeTemplate,
		},
	}, nil
}

func (p testWalletPlugin) GetPluginInfo() (PluginInfo, error) {
	return PluginInfo{
		Name: p.name + " Test Wallet",
		ID:   p.symbol + "-test-wallet",
		Type: PluginTypeWallet,
	}, nil
}

func TestServerPluginAddition(t *testing.T) {

	server := &Server{
		Logger: hclog.Default(),
	}
	pluginImpl := testWalletPlugin{
		name:   "Monero",
		symbol: "XMR",
	}
	serverPlugin := BuildWalletPlugin(pluginImpl)

	err := server.AddPlugin(serverPlugin)
	testutils.Ok(t, err)

	testutils.Equals(t, server.started, false)

	testutils.Equals(t, server.plugins[0], serverPlugin)

	router := chi.NewRouter()
	err = server.initializeAllPlugins(router)
	testutils.Ok(t, err)

	pInfo, err := pluginImpl.GetPluginInfo()
	testutils.Ok(t, err)

	testutils.Equals(t, pInfo, server.pluginInfos[0])

	r := router.Routes()
	testutils.Equals(t, r[0].Pattern, pInfo.GetAddress()+"/*")
}
