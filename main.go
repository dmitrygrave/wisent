package main

import (
	"github.com/dmitrygrave/wisent/utils/logging"
)

func main() {
	logging.Warn("This is a info message")

	logging.Fatalf("This is a fatal message %s", "Hello world")
}
