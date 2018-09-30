package main

import (
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flweb"
)

func main() {
	flconfig.LoadConfiguration()
	flweb.NewApplication(flactions.NewFacebookHTTPClient()).Start()
}
