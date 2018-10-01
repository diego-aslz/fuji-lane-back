package main

import (
	"log"

	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flentities"
)

func main() {
	flconfig.LoadConfiguration()
	if err := flentities.Seed(); err != nil {
		log.Fatal(err.Error())
	}
}
