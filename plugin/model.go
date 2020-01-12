package plugin

type PluginType string

const ActivatePluginEndpoint = "/activate"

const (
	PluginTypeWallet           PluginType = "wallet"
	PluginTypeExchangeRate     PluginType = "exchange-rate"
	PluginTypeCurrencyMetadata PluginType = "currency-metadata"
)

type PluginInfo struct {
	Name string     `json:"name"`
	ID   string     `json:"id"`
	Type PluginType `json:"pluginType"`
}

func (info PluginInfo) GetAddress() string {
	return "/plugin/" + info.ID
}

type VendoPunktoPlugin interface {
	GetPluginInfo() (PluginInfo, error)
}
