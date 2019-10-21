package main

import (
	"log"

	"github.com/friedenberg/zapstackdriver"
)

func main() {
	logger, err := zapstackdriver.NewLogger(
		zapstackdriver.FieldServiceContext{
			Service: "example",
			Version: "1-b",
		},
		[]zapstackdriver.FieldSourceReference{
			zapstackdriver.FieldSourceReference{
				Repository: "https://github.com/friedenberg/zapstackdriver",
				RevisionId: "1-b",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	logger.Error("oops test")
}
