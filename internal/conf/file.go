package conf

import (
	"os"

	"github.com/spf13/viper"
)

// InitializeConfigFile will read a file from the provided path
// if it's found, it will read the config with viper.
// If not, it will create and initialize the file with defaults
func InitializeConfigFile(path string) (bool, error) {
	viper.SetConfigFile(path)

	if _, e := os.Stat(path); os.IsNotExist(e) {
		// create the file
		emptyFile, err := os.Create(path)
		if err != nil {
			return false, err
		}

		defer emptyFile.Close()

		err = viper.WriteConfig()
		return true, err
	}

	err := viper.ReadInConfig()
	return false, err
}
