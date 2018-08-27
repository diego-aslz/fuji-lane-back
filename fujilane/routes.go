package fujilane

import (
	"github.com/gin-gonic/gin"
)

const (
	statusPath         = "/status"
	facebookSignInPath = "/sign_in/facebook"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, a.statusRoute)
	e.POST(facebookSignInPath, a.facebookSignInRoute)
}

func (a *Application) statusRoute(c *gin.Context) {
	c.JSON(200, gin.H{"status": "active"})
}

type facebookSignInBody struct {
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
}

func (a *Application) facebookSignInRoute(c *gin.Context) {
	body := &facebookSignInBody{}
	err := c.BindJSON(body)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = a.facebook.validate(body.AccessToken, body.ID)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, gin.H{"status": "active"})
}
