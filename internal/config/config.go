package config

import "flag"

// Config stores all configuration of the application.
type Config struct {
	Addr         string
	ApiAddr      string
	DBDriver     string
	DBSource     string
	MigrationURL string
}

// LoadConfig: loads configuration variables
func LoadConfig() (config Config) {
	flag.StringVar(&config.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.ApiAddr, "apiaddr", ":3000", "HTTP network address for API")
	flag.StringVar(&config.DBDriver, "drvr", "pgx", "Database driver")
	flag.StringVar(&config.DBSource, "dsn", "", "Database source")
	flag.StringVar(&config.MigrationURL, "mgr", "", "Migration URL")
	flag.Parse()

	return
}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
// type Config struct {
// 	DBDriver     string `mapstructure:"DB_DRIVER"`
// 	DBSource     string `mapstructure:"DB_SOURCE"`
// 	MigrationURL string `mapstructure:"MIGRATION_URL"`
// }

// // loadConfig: loads configuration variables
// func (app *application) loadConfig(path string) (config Config) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env")

// 	viper.AutomaticEnv()

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
// 		app.errorLog.Output(2, trace)
// 		return
// 	}
// 	err = viper.Unmarshal(&config)
// 	if err != nil {
// 		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
// 		app.errorLog.Output(2, trace)
// 		return
// 	}
// 	return config
// }
