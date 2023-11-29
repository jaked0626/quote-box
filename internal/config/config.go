package config

import "flag"

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
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
	flag.StringVar(&config.ApiAddr, "addr", ":3000", "HTTP network address for API")
	flag.StringVar(&config.DBDriver, "drvr", "pgx", "Database driver")
	flag.StringVar(&config.DBSource, "dsn", "", "Database source")
	flag.StringVar(&config.MigrationURL, "mgr", "", "Migration URL")
	flag.Parse()

	return
}
