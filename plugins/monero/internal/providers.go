package monero

import (
	"errors"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/spf13/viper"
)

type Container struct {
	Plugin       plugin.WalletPlugin
	Server       *plugin.Server
	WalletClient wallet.Client
}

var Providers = wire.NewSet(
	newMoneroWalletPlugin,
	newMoneroWalletClient,
	plugin.NewServer,
	newContainer)

func newContainer(
	plugin plugin.WalletPlugin,
	server *plugin.Server,
	walletClient wallet.Client,
) *Container {
	return &Container{
		Plugin:       plugin,
		WalletClient: walletClient,
		Server:       server,
	}
}

func newMoneroWalletPlugin(logger hclog.Logger, client wallet.Client) (plugin.WalletPlugin, error) {
	return moneroWalletPlugin{
		client: client,
		logger: logger,
	}, nil
}

func newMoneroWalletClient(globalLogger hclog.Logger) (wallet.Client, error) {
	baseURL := viper.GetString("monero.wallet_rpc_url")
	logger := globalLogger.Named("monero")

	if baseURL == "" {
		return nil, errors.New("A Monero Wallet RPC url must be provided")
	}

	client := wallet.New(wallet.Config{
		Address: baseURL + "/json_rpc",
	})

	version, err := client.GetVersion()

	if err != nil {
		logger.Warn("Failed to connect to Monero Wallet RPC", "url", baseURL)
		return nil, err
	}

	logger.Info("Monero Wallet RPC Test Success", "version", version.Version)

	return client, nil
}
