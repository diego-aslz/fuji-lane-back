package main

import (
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flconfig"
	fujilane "github.com/nerde/fuji-lane-back/fujilane"
)

func main() {
	flconfig.LoadConfiguration()
	fujilane.NewApplication(flactions.NewFacebookHTTPClient()).Start()
}
