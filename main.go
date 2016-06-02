package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/aristanetworks/openconfigbeat/beater"
)

func main() {
	err := beat.Run("openconfigbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
