package core

import (
	"log"
	"os"
)

func ErrorsHandler(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
