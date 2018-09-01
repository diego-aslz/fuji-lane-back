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
	e.GET(statusPath, a.routeStatus)
	e.POST(facebookSignInPath, a.routeFacebookSignIn)
}

type facebookSignInBody struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
