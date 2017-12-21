package main

import (
	"fmt"
	"os"

	"github.com/dmitrygrave/wisent/utils/config"
	"github.com/dmitrygrave/wisent/utils/logging"
	"github.com/dmitrygrave/wisent/utils/signals"
)

func main() {
	configFile, isSet := os.LookupEnv("WISENTCONFIG")

	if isSet != true {
		// No configuration file present
		// TODO: (maybe) provide a default config
		fmt.Fprintln(os.Stderr, "WISENTCONFIG is not set! No configuration found. Exiting...")
		os.Exit(1)
	}

	config.InitConfig(configFile)
	logging.InitLogging(config.Env())

	signals.HandleInterrupts()

	logging.Infof("Currently on env: %s with config file %s%s", config.Env(), config.LogDirectory, config.LogFilename())
}
