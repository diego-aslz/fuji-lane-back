package fujilane

import (
	"github.com/gin-gonic/gin"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	facebookSignInPath = "/sign_in/facebook"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, a.routeStatus)

	e.POST(signUpPath, a.routeSignUp)
	e.POST(facebookSignInPath, a.routeFacebookSignIn)
}
