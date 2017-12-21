package main

import (
	"github.com/dmitrygrave/wisent/utils/config"
	"github.com/dmitrygrave/wisent/utils/logging"
)

func main() {
	config.InitConfig("config/config.dev.json")
	logging.Init(config.Env())

	logging.Infof("Currently on env: %s with config file %s", config.Env(), config.LogFilename())
}
