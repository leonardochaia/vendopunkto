package conf

import (
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

// The keys for the Runtime config
const (
	// PluginHostsKey string separated list of plugin URLs
	PluginHostsKey = "plugin_hosts"

	ExchangeRatesPluginKey    = "exchange_rates_plugin"
	CurrencyMetadataPluginKey = "currency_metadata_plugin"

	// WalletPollingIntervalKey determines how often to poll wallets for new transfers
	WalletPollingIntervalKey = "wallet_poller_interval"

	PricingCurrenciesKey      = "pricing_currencies"
	DefaultPricingCurrencyKey = "default_pricing_currency"
	PaymentMethodsKey         = "payment_methods"
)

// Runtime is the config that is used on runtime and can be changed after
// the application has started
type Runtime struct {
	*viper.Viper
	logger hclog.Logger
}

func setRuntimeDefaults(r *Runtime) {
	r.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	r.AutomaticEnv()                                   // Automatically use environment variables where available
	r.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environment variables use underscores instead of periods

	r.SetDefault(PluginHostsKey, []string{})
	r.SetDefault(ExchangeRatesPluginKey, "")
	r.SetDefault(CurrencyMetadataPluginKey, "")

	r.SetDefault(WalletPollingIntervalKey, "10s")

	r.SetDefault(PricingCurrenciesKey, []string{"usd", "btc"})
	r.SetDefault(DefaultPricingCurrencyKey, "usd")

	r.SetDefault(PaymentMethodsKey, []string{"btc", "xmr"})
}

func (r *Runtime) GetPluginHosts() []string {
	return r.GetStringSlice(PluginHostsKey)
}

func (r *Runtime) GetExchangeRatesPlugin() string {
	return r.GetString(ExchangeRatesPluginKey)
}

func (r *Runtime) GetCurrencyMetadataPlugin() string {
	return r.GetString(CurrencyMetadataPluginKey)
}

func (r *Runtime) GetWalletPollingInterval() time.Duration {
	return r.GetDuration(WalletPollingIntervalKey)
}

func (r *Runtime) GetPricingCurrencies() []string {
	return r.GetStringSlice(PricingCurrenciesKey)
}

func (r *Runtime) GetPaymentMethods() []string {
	return r.GetStringSlice(PaymentMethodsKey)
}

func (r *Runtime) GetDefaultPricingCurrency() string {
	return r.GetString(DefaultPricingCurrencyKey)
}

// InitializeConfigFile will read a file from the provided path
// if it's found, it will read the config with viper.
// If not, it will create and initialize the file with defaults
func (r *Runtime) InitializeConfigFile(path string) (bool, error) {
	r.SetConfigFile(path)

	logger := r.logger.With("path", path)

	if _, e := os.Stat(path); os.IsNotExist(e) {

		logger.Info("Could not find runtime config. Creating file with defaults.")

		// create the file
		emptyFile, err := os.Create(path)
		if err != nil {
			return false, err
		}

		defer emptyFile.Close()

		err = r.WriteConfig()
		return true, err
	}

	logger.Info("Reading runtime config file")

	err := r.ReadInConfig()
	return false, err
}
