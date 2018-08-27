package main

import (
	fujilane "github.com/nerde/fuji-lane-back/fujilane"
)

func main() {
	app := fujilane.NewApplication(fujilane.NewFacebookHTTPClient())
	app.Start()
}
