package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// LogConfig contains the logging configuration options for the application
type LogConfig struct {
	Directory  string `json:"directory"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
}

// WebConfig contains the web server configuration options for the application
type WebConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Config contains all the user defined configuration variables for the
// application
type Config struct {
	// The environment the application is running in
	Env string `json:"env"`
	// Logging configuration
	Log LogConfig `json:"log"`
	// Web Server configuration
	Web WebConfig `json:"web"`
}

// AppConfig holds the configuration for the application
var AppConfig *Config

// InitConfig takes a json filename and unmarshals it into the Config struct
func InitConfig(file string) {
	raw, err := ioutil.ReadFile(file)

	if err != nil {
		// log to printf because logging is unavailable at this time
		fmt.Fprintf(os.Stderr, "Unable to open configuration file: %s", file)
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &AppConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading configuration file")
		os.Exit(1)
	}
}

// Env returns the current environment of the application
func Env() string {
	return AppConfig.Env
}

// Log returns the logging configuration options
func Log() *LogConfig {
	return &AppConfig.Log
}

// Web returns the web server configuration optiosn
func Web() *WebConfig {
	return &AppConfig.Web
}
