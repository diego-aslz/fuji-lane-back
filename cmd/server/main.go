package main

import (
	"github.com/gin-gonic/gin"
	fujilane "github.com/nerde/fuji-lane-back/fujilane"
)

func main() {
	r := gin.Default()
	fujilane.AddMiddleware(r)
	fujilane.AddRoutes(r)
	r.Run()
}
