package conf

import (
	"os"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/spf13/viper"
)

// Providers for wire
var Providers = wire.NewSet(LoadStartupConfig, NewManager)

// LoadStartupConfig uses viper to unmarshal the config into struct
func LoadStartupConfig() (Startup, error) {
	var startupConf Startup
	err := viper.Unmarshal(&startupConf)
	return startupConf, err
}

// NewRuntimeConfig creates the runtime config
func NewRuntimeConfig(
	logger hclog.Logger,
	configPath string) *Runtime {
	r := &Runtime{
		Viper:  viper.New(),
		logger: logger.Named("runtimeConf"),
	}

	setRuntimeDefaults(r)

	// If a config file is found, initialize it
	if configPath != "" {
		_, err := r.InitializeConfigFile(configPath)
		if err != nil {
			errors.Errorf("Error ocurred while initializing runtime config at %s: %s", configPath, err.Error())
			os.Exit(1)
		}
	}

	return r
}

// NewManager creates a Manager instance
func NewManager(
	logger hclog.Logger,
	runtime *Runtime) Manager {
	return Manager{
		logger:  logger,
		runtime: runtime,
	}
}
