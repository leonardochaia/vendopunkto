package monero

import (
	"errors"

	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func CreateMoneroClient() (wallet.Client, error) {
	baseURL := viper.GetString("monero.wallet_rpc_url")

	if baseURL == "" {
		return nil, errors.New("A Monero Wallet RPC url must be provided")
	}

	client := wallet.New(wallet.Config{
		Address: baseURL + "/json_rpc",
	})

	zap.S().Infow("Connecting to Monero Wallet RPC", "package", "monero", "url", baseURL)

	version, err := client.GetVersion()

	if err != nil {
		return nil, err
	}

	zap.S().Infow("Monero Wallet RPC", "package", "monero", "version", version)

	return client, nil
}
