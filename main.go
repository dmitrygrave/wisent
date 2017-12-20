package main

import (
	"github.com/dmitrygrave/wisent/utils/config"
	"github.com/dmitrygrave/wisent/utils/logging"
)

func main() {
	config.InitConfig("config/config.prod.json")

	logging.Infof("Currently on env: %s with config file %s", config.Env(), config.LogFilename())
}
