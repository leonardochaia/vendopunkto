package conf

import (
	"github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	viper.SetDefault("logger.level", "info")

	// Pidfile
	viper.SetDefault("pidfile", "")

	// Profiler config
	viper.SetDefault("profiler.enabled", false)
	viper.SetDefault("profiler.host", "")
	viper.SetDefault("profiler.port", "6060")

	// Server Configuration
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.log_requests", true)
	viper.SetDefault("server.profiler_enabled", false)
	viper.SetDefault("server.profiler_path", "/debug")

	// Database Settings
	viper.SetDefault("storage.type", "postgres")
	viper.SetDefault("storage.username", "postgres")
	viper.SetDefault("storage.password", "password")
	viper.SetDefault("storage.host", "postgres")
	viper.SetDefault("storage.port", 5432)
	viper.SetDefault("storage.database", "vendopunkto")
	viper.SetDefault("storage.sslmode", "disable")
	// viper.SetDefault("storage.max_connections", 80)

	// Plugins
	// PLUGINS_ENABLED="wallet|http://localhost3333"
	viper.SetDefault("plugins.enabled", []string{})
	viper.SetDefault("plugins.server.host", "0.0.0.0")
	viper.SetDefault("plugins.server.port", "9080")
	viper.SetDefault("plugins.server.plugin_host_address", "http://localhost:9080")
	viper.SetDefault("plugins.default_exchange_rates", "")

}
