package cmd

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leonardochaia/vendopunkto/internal/conf"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cobra.Command, args []string) { // Initialize the databse

			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "vendopunkto",
				Output: os.Stdout,
				Level:  hclog.LevelFromString(strings.ToLower(viper.GetString("logger.level"))),
				Color:  hclog.AutoColor,
			})

			// Create the server (uses wire DI)
			s, err := NewServer(logger)
			if err != nil {
				logger.Error("Could not create server",
					"error", err,
				)
				os.Exit(1)
			}

			err = s.ListenAndServe()

			if err != nil {
				logger.Error("Could not start server",
					"error", err,
				)
				os.Exit(1)
			}

			defer s.Close()

			<-conf.Stop.Chan() // Wait until StopChan
			conf.Stop.Wait()   // Wait until everyone cleans up
		},
	}
)
