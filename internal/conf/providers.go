package conf

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Providers for wire
var Providers = wire.NewSet(LoadStartupConfig)

// LoadStartupConfig uses viper to unmarshal the config into struct
func LoadStartupConfig() (Startup, error) {
	var startupConf Startup
	err := viper.Unmarshal(&startupConf)
	return startupConf, err
}
