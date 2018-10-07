package main

import (
	"log"

	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flentities"
)

func main() {
	flconfig.LoadConfiguration()

	if flconfig.Config.Stage != "test" && flconfig.Config.Stage != "development" {
		log.Fatalf("Command \"reset\" should not be run in %s stage!", flconfig.Config.Stage)
		return
	}

	if err := flentities.Reset(); err != nil {
		log.Fatal(err.Error())
	}
}
