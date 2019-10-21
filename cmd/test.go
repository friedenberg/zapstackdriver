package main

import (
	"github.com/friedenberg/zapstackdriver"
)

func main() {
	loggerConfig, _ := zapstackdriver.NewConfig()
	logger, _ := loggerConfig.Build()

	logger.Error("oops test")
}
