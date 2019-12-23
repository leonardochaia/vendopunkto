package cmd

import (
	"fmt"
	"os"
	"strings"

	_ "net/http/pprof" // Import for pprof

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leonardochaia/vendopunkto/internal/conf"
)

var (

	// Config and global logger
	configFile string
	pidFile    string

	// The Root Cli Handler
	rootCmd = &cobra.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Create Pid File
			pidFile = viper.GetString("pidfile")
			if pidFile != "" {
				file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
				defer file.Close()
				_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Remove Pid file
			if pidFile != "" {
				os.Remove(pidFile)
			}
		},
	}
)

// Execute starts the program
func Execute() {
	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

// This is the main initializer handling cli, config and log
func init() {
	// Initialize configuration
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Sets up the config file, environment etc
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	// If a config file is found, initialize it
	if configFile != "" {
		created, err := conf.InitializeConfigFile(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error ocurred while initializing config: %s ERROR: %s\n", configFile, err.Error())
			os.Exit(1)
		}

		if created {
			fmt.Fprintf(os.Stdout, "New config file created at: %s\n", configFile)
			os.Exit(0)
		}
	}
}
