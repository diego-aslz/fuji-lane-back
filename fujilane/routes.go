package fujilane

import "github.com/gin-gonic/gin"

const (
	statusPath = "/status"
)

// AddRoutes to a Gin Engine
func AddRoutes(e *gin.Engine) {
	e.GET(statusPath, statusRoute)
}

func statusRoute(c *gin.Context) {
	c.JSON(200, gin.H{"status": "active"})
}
