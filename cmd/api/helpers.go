package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError: writes an error message and stack trace to errorLog,
// then sends 500 Internal Server Error response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError: sends specific status code and corresponding description to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound: convenience wrapper around clientError that sends a 404 Not Found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// badRequest: convenience wrapper around clientError that sends a 400 Bad Request response
func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
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
