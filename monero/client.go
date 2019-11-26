package monero

import (
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/spf13/viper"
)

func CreateMoneroClient(globalLogger hclog.Logger) (wallet.Client, error) {
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
