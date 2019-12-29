package development

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/leonardochaia/vendopunkto/unit"
	"github.com/rs/xid"
)

type fakeAddress struct {
	address      string
	creationTime time.Time
}

type fakeWalletPlugin struct {
	symbol         string
	name           string
	qrCodeTemplate string
	addresses      []fakeAddress
}

func (p *fakeWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	h := sha256.Sum256([]byte(p.symbol + "-fake-" + xid.New().String()))
	addr := fakeAddress{
		address:      fmt.Sprintf("%x", h),
		creationTime: time.Now(),
	}
	p.addresses = append(p.addresses, addr)
	return addr.address, nil
}

func (p *fakeWalletPlugin) GetWalletInfo() (plugin.WalletPluginInfo, error) {
	return plugin.WalletPluginInfo{
		Currency: plugin.WalletPluginCurrency{
			Name:           p.name,
			Symbol:         p.symbol,
			QRCodeTemplate: p.qrCodeTemplate,
		},
	}, nil
}

func (p *fakeWalletPlugin) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: p.name + " Fake Wallet",
		ID:   p.symbol + "-fake-wallet",
		Type: plugin.PluginTypeWallet,
	}, nil
}

func (p *fakeWalletPlugin) GetIncomingTransfers(params plugin.WalletPluginIncomingTransferParams) (
	[]plugin.WalletPluginIncomingTransferResult, error) {

	output := []plugin.WalletPluginIncomingTransferResult{}

	now := time.Now()
	for _, addr := range p.addresses {
		confTime := addr.creationTime.Add(30 * time.Second)
		confirmations := uint64(0)
		if confTime.Before(now) {
			confirmations = uint64(now.Unix()) - uint64(addr.creationTime.Unix())
		}

		memPoolTime := addr.creationTime.Add(10 * time.Second)
		height := uint64(memPoolTime.Unix())

		if memPoolTime.Before(now) {
			h := sha256.Sum256([]byte(fmt.Sprintf("fake-tx-hash-%d", height)))
			tx := plugin.WalletPluginIncomingTransferResult{
				TxHash:        fmt.Sprintf("%x", h),
				Address:       addr.address,
				BlockHeight:   height,
				Amount:        unit.NewFromFloat(1),
				Confirmations: confirmations,
			}

			output = append(output, tx)
		}
	}

	return output, nil
}
