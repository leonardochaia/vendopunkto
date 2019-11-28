package cmd

import (
	"os"

	"github.com/spf13/cobra"

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

			// Create the server (uses wire DI)
			s, err := NewServer(globalLogger)
			if err != nil {
				globalLogger.Error("Could not create server",
					"error", err,
				)
				os.Exit(1)
			}

			err = s.ListenAndServe()

			if err != nil {
				globalLogger.Error("Could not start server",
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
