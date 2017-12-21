package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config contains all the user defined configuration variables for the
// application
type Config struct {
	// The environment the application is running in
	Env string `json:"env"`
	// Logging configuration
	Log struct {
		Directory  string `json:"directory"`
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
	} `json:"log"`
}

// AppConfig holds the configuration for the application
var AppConfig *Config

// InitConfig takes a json filename and unmarshals it into the Config struct
func InitConfig(file string) {
	raw, readErr := ioutil.ReadFile(file)

	if readErr != nil {
		// log to printf because logging is unavailable at this time
		fmt.Fprintf(os.Stderr, "Unable to open configuration file: %s", file)
		os.Exit(1)
	}

	unmarshalErr := json.Unmarshal(raw, &AppConfig)

	if unmarshalErr != nil {
		fmt.Fprint(os.Stderr, "Error reading configuration file")
		os.Exit(1)
	}
}

// Env returns the current environment of the application
func Env() string {
	return AppConfig.Env
}

// LogDirectory returns the directory where the log file is located
func LogDirectory() string {
	return AppConfig.Log.Directory
}

// LogFilename returns the name of the log file
func LogFilename() string {
	return AppConfig.Log.Filename
}

// LogMaxSize returns the maximum size of the log file in megabytes
func LogMaxSize() int {
	return AppConfig.Log.MaxSize
}

// LogMaxBackups returns the number of logfile backups which should be made
func LogMaxBackups() int {
	return AppConfig.Log.MaxBackups
}

// LogMaxAge returns the maximum time to keep log file backups in days
func LogMaxAge() int {
	return AppConfig.Log.MaxAge
}
