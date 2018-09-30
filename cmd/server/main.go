package main

import (
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/flweb"
)

func main() {
	flconfig.LoadConfiguration()
	flweb.NewApplication(flservices.NewFacebookHTTPClient()).Start()
}
