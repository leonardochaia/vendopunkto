package monero

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("monero.server.host", "0.0.0.0")
	viper.SetDefault("monero.server.port", "4300")

	viper.SetDefault("plugin.server.host", "0.0.0.0")
	viper.SetDefault("plugin.server.port", "4200")
}
