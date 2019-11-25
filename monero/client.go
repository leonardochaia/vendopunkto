package monero

import (
	"errors"

	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
* Creates the wallet.Client using the config
 */
func CreateMoneroClient() (wallet.Client, error) {
	baseUrl := viper.GetString("monero.wallet_rpc_url")

	if baseUrl == "" {
		return nil, errors.New("A Monero Wallet RPC url must be provided")
	}

	client := wallet.New(wallet.Config{
		Address: baseUrl + "/json_rpc",
	})

	zap.S().Infow("Connecting to Monero Wallet RPC", "url", baseUrl)

	version, err := client.GetVersion()

	if err != nil {
		return nil, err
	}

	zap.S().Infow("Monero Wallet RPC", "version", version)

	return client, nil
}
