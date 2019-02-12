package main

import (
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flweb"
)

func main() {
	flconfig.LoadConfiguration()
	app, err := flweb.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Start()
}
