package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dmitrygrave/wisent/utils/config"
	"github.com/dmitrygrave/wisent/utils/logging"
	"github.com/dmitrygrave/wisent/utils/signals"
)

func runServer() {
	// Server loop
	conf := config.Web()
	server := &http.Server{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))

	if err != nil {
		logging.Errorf("Error announcing on host: %s port %d", conf.Host, conf.Port)
		return
	}

	logging.Infof("Listening on %s", listener.Addr())

	signals.AppendInterrupt(func() {
		server.Shutdown(context.Background())
	})

	err = server.Serve(listener)

	if err == nil {
		logging.Panic("Serve returned nil! This should never happen!")
	}

	if err == http.ErrServerClosed {
		return
	}

	// TODO: Check errors
}

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

	runServer()
}
