package main

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	monero "github.com/leonardochaia/vendopunkto/plugins/monero/internal"
	"github.com/spf13/viper"
)

func main() {

	// Sets up the config file, environment etc
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "monero-plugin",
		Output: os.Stdout,
		Level:  hclog.LevelFromString(strings.ToUpper(viper.GetString("logger.level"))),
	})

	monero, err := monero.NewMoneroWalletPlugin(logger)

	if err != nil {
		panic(err.Error())
	}

	server := plugin.NewServer()

	err = server.AddPlugin(plugin.BuildWalletPlugin(monero))
	if err != nil {
		panic(err.Error())
	}

	err = server.Start("0.0.0.0:3333")

	if err != nil {
		panic(err.Error())
	}
}
