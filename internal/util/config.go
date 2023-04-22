package util

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

// LoadConfig: loads configuration variables
func LoadConfig(path string, errorLog *log.Logger) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		errorLog.Output(2, trace)
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		errorLog.Output(2, trace)
		return
	}
	return config
}
