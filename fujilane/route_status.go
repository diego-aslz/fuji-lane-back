package fujilane

import "github.com/gin-gonic/gin"

func (a *Application) routeStatus(c *gin.Context) {
	c.JSON(200, gin.H{"status": "active"})
}
