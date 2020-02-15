package conf

import (
	"strings"
	"time"

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

// Runtime is the config that is used on runtime and can be changed after
// the application has started
type Runtime struct {
	*viper.Viper
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
