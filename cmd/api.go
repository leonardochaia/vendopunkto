package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/leonardochaia/vendopunkto/conf"
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
			s, err := NewServer()
			if err != nil {
				logger.Fatalw("Could not create server",
					"error", err,
				)
			}
			err = s.ListenAndServe()
			if err != nil {
				logger.Fatalw("Could not start server",
					"error", err,
				)
			}

			defer s.Close()

			<-conf.Stop.Chan() // Wait until StopChan
			conf.Stop.Wait()   // Wait until everyone cleans up
			zap.L().Sync()     // Flush the logger

		},
	}
)
