package main

import (
	fujilane "github.com/nerde/fuji-lane-back/fujilane"
)

func main() {
	fujilane.LoadConfiguration()
	fujilane.NewApplication(fujilane.NewFacebookHTTPClient()).Start()
}
