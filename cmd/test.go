package main

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/friedenberg/zapstackdriver"
)

func main() {
	loggerConfig, _ := zapstackdriver.NewConfig()
	logger, _ := loggerConfig.Build()

	logger.Error("oops")
	fmt.Println(string(debug.Stack()))
}

func test() (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(0)
}
