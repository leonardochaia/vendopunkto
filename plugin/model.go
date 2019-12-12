package plugin

type PluginType string

const ActivatePluginEndpoint = "/activate"

const (
	PluginTypeWallet       PluginType = "wallet"
	PluginTypeExchangeRate PluginType = "exchange-rate"
)

type PluginInfo struct {
	Name string     `json:"name"`
	ID   string     `json:"id"`
	Type PluginType `json:"pluginType"`
}

type VendoPunktoPlugin interface {
	GetPluginInfo() (PluginInfo, error)
}
