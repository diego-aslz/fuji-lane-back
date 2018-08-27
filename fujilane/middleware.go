package fujilane

import "github.com/gin-gonic/gin"

// AddMiddleware configuration to a Gin Engine
func AddMiddleware(e *gin.Engine) {
	e.Use(gin.Recovery())
}
