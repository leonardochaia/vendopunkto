package rates

import (
	"strings"

	"github.com/spf13/viper"
)

func init() {

	// Sets up the config file, environment etc
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	viper.SetDefault("plugin.server.host", "0.0.0.0")
	viper.SetDefault("plugin.server.port", "4201")
}
