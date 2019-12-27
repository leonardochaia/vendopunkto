package conf

import (
	"github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	viper.SetDefault("logger.level", "info")

	// Pidfile
	viper.SetDefault("pidfile", "")

	// Public Server Configuration
	viper.SetDefault("server.public.host", "0.0.0.0")
	viper.SetDefault("server.public.port", "8080")
	viper.SetDefault("server.log_requests", true)
	viper.SetDefault("server.profiler_enabled", false)
	viper.SetDefault("server.profiler_path", "/debug")

	// Internal Server Configuration
	viper.SetDefault("server.internal.host", "0.0.0.0")
	viper.SetDefault("server.internal.port", "9080")

	// Database Settings
	viper.SetDefault("storage.type", "postgres")
	viper.SetDefault("storage.username", "postgres")
	viper.SetDefault("storage.password", "password")
	viper.SetDefault("storage.host", "postgres")
	viper.SetDefault("storage.port", 5432)
	viper.SetDefault("storage.database", "vendopunkto")
	viper.SetDefault("storage.sslmode", "disable")

	// Plugins
	// string separated list of plugin URLs to activate on startup
	viper.SetDefault("plugins.enabled", []string{})

	// this is the address that VendoPunkto will provide to plugins
	// for them to reach back
	viper.SetDefault("plugins.default_exchange_rates", "")

	// How often to poll wallets for new transfers
	viper.SetDefault("plugins.wallet_poll_interval", "10s")

}
