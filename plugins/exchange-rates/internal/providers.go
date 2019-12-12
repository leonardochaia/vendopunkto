package rates

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	gecko "github.com/superoo7/go-gecko/v3"
)

type Container struct {
	Plugin      plugin.ExchangeRatesPlugin
	Server      *plugin.Server
	RatesClient *gecko.Client
}

var Providers = wire.NewSet(
	newGeckoExchangeRatesPlugin,
	newGeckoClient,
	plugin.NewServer,
	newContainer)

func newContainer(
	plugin plugin.ExchangeRatesPlugin,
	server *plugin.Server,
	ratesClient *gecko.Client,
) *Container {
	return &Container{
		Plugin:      plugin,
		RatesClient: ratesClient,
		Server:      server,
	}
}

func newGeckoExchangeRatesPlugin(logger hclog.Logger, client *gecko.Client) (plugin.ExchangeRatesPlugin, error) {
	return geckoExchangeRatesPlugin{
		client: client,
	}, nil
}

func newGeckoClient(globalLogger hclog.Logger) (*gecko.Client, error) {
	logger := globalLogger.Named("monero")

	client := gecko.NewClient(nil)

	result, err := client.Ping()

	if err != nil {
		logger.Warn("Failed to connect to Gecko API", "error", err)
		return nil, err
	}

	logger.Info("Gecko API Test Success", "ping", result.GeckoSays)

	return client, nil
}
