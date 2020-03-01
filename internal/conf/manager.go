package conf

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
)

// Manager is in charge of handling config logics
type Manager struct {
	logger  hclog.Logger
	runtime *Runtime
}

// SaveConfiguration updates conf.Runtime with the new key/value pair.
// it will also save the file to disk.
func (m Manager) SaveConfiguration(key string, newValue interface{}) error {
	const op errors.Op = "conf.manager.saveConfiguration"

	logger := m.logger.With("key", key, "value", newValue)
	logger.Info("Updating configuration")
	m.runtime.Set(key, newValue)

	// TODO: validate that the provided values make sense.

	if key == PricingCurrenciesKey {
		pricingCurrencies := m.runtime.GetPricingCurrencies()
		defaultCurrency := m.runtime.GetDefaultPricingCurrency()
		found := false
		for _, c := range pricingCurrencies {
			if c == defaultCurrency {
				found = true
				continue
			}
		}

		if !found {
			logger.Info("The "+DefaultPricingCurrencyKey+" is no longer present in "+PricingCurrenciesKey+
				". Falling back to first pricing currency", "newDefault", pricingCurrencies[0])
			m.runtime.Set(DefaultPricingCurrencyKey, pricingCurrencies[0])
		}
	}

	err := m.runtime.WriteConfig()
	if err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	return nil
}
