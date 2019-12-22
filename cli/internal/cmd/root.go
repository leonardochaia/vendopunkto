package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/conf"
)

var (

	// Config and global logger
	configFile              string
	pidFile                 string
	vendoPunktoHost         string
	vendoPunktoPublicPort   uint
	vendoPunktoInternalPort uint
	publicClient            clients.PublicClient
	internalClient          clients.InternalClient
	httpClient              = clients.NewHTTPClient()

	// The Root Cli Handler
	rootCmd = &cobra.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if vendoPunktoHost == "" {
				return errors.Str("Empty host provided")
			}

			if vendoPunktoPublicPort == 0 {
				return errors.Str("Empty public port provided")
			}

			if vendoPunktoInternalPort == 0 {
				return errors.Str("Empty internal port provided")
			}

			var err error
			publicClient, err = clients.NewPublicClient(
				fmt.Sprintf("%s:%d", vendoPunktoHost, vendoPunktoPublicPort),
				httpClient)
			if err != nil {
				return err
			}

			internalClient, err = clients.NewInternalClient(
				fmt.Sprintf("%s:%d", vendoPunktoHost, vendoPunktoInternalPort),
				httpClient)

			return err
		},
	}
)

// Execute starts the program
func Execute() {
	// Run the program
	rootCmd.Execute()
}

// This is the main initializer handling cli, config and log
func init() {

	// Initialize configuration
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")

	rootCmd.PersistentFlags().StringVarP(&vendoPunktoHost, "host", "H", "http://localhost",
		"The VendoPunkto host URL, without port, by default http://localhost")

	rootCmd.PersistentFlags().UintVar(&vendoPunktoPublicPort, "public-port", 8080,
		"The VendoPunkto public port, by default 8080")

	rootCmd.PersistentFlags().UintVar(&vendoPunktoInternalPort, "internal-port", 9080,
		"The VendoPunkto internal port, by default 9080")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Sets up the config file, environment etc
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	// If a config file is found, read it in.
	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read config file: %s ERROR: %s\n", configFile, err.Error())
			os.Exit(1)
		}

	}
}
