package main

import (
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/fljobs"
)

func main() {
	flconfig.LoadConfiguration()
	fljobs.NewDefaultApplication().Start()
}
