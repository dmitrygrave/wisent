package main

import (
	"github.com/dmitrygrave/wisent/utils/logging"
)

func main() {
	logging.Error("This is a Error message")

	logging.Fatalf("This is a fatal message %s", "Hello world")
}
