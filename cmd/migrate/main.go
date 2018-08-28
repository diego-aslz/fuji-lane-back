package main

import (
	"log"

	"github.com/nerde/fuji-lane-back/fujilane"
)

func main() {
	if err := fujilane.Migrate(); err != nil {
		log.Fatal(err.Error())
	}
}
