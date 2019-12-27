package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
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

func (p testWalletPlugin) GetIncomingTransfers(params WalletPluginIncomingTransferParams) ([]WalletPluginIncomingTransferResult, error) {
	return nil, nil
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

func TestActivationHandler(t *testing.T) {
	expected := PluginInfo{
		Name: "Test",
		ID:   "test",
	}
	server := &Server{
		Logger:      hclog.Default(),
		pluginInfos: []PluginInfo{expected},
	}

	req, err := http.NewRequest("POST", "/activate", bytes.NewBuffer([]byte{}))
	testutils.Ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(errors.WrapHandler(server.activatePluginHandler))

	handler.ServeHTTP(rr, req)

	var result []PluginInfo

	err = json.NewDecoder(rr.Body).Decode(&result)
	testutils.Ok(t, err)

	testutils.Assert(t, result[0].ID == expected.ID,
		"Expected plugin info to be returned")
	testutils.Assert(t, rr.Code == 200, "Expected request code to be 200")
}
