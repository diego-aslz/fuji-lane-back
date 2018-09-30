package fujilane

import (
	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	signInPath         = "/sign_in"
	facebookSignInPath = "/sign_in/facebook"

	accountsPath   = "/accounts"
	propertiesPath = "/properties"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	a.route(e.GET, statusPath, a.status)

	a.route(e.POST, signUpPath, a.signUp)
	a.route(e.POST, signInPath, a.signIn)
	a.route(e.POST, facebookSignInPath, a.facebookSignIn)

	a.route(e.POST, accountsPath, a.accountsCreate)
	a.route(e.POST, propertiesPath, a.propertiesCreate)
}

type ginMethod func(string, ...gin.HandlerFunc) gin.IRoutes

func (a *Application) route(method ginMethod, path string, next func(*Context)) {
	method(path, func(c *gin.Context) {
		next(&Context{context: c, now: a.timeFunc})
	})
}

func (a *Application) status(c *Context) {
	c.action = &flactions.Status{}
	performAction(c)
}

func (a *Application) signUp(c *Context) {
	c.action = &flactions.SignUp{}
	routeWithRepository(loadActionBody(validateActionBody(performAction)))(c)
}

func (a *Application) signIn(c *Context) {
	c.action = &flactions.SignIn{}
	routeWithRepository(loadActionBody(performAction))(c)
}

func (a *Application) facebookSignIn(c *Context) {
	c.action = flactions.NewFacebookSignIn(a.facebookClient)
	routeWithRepository(loadActionBody(performAction))(c)
}

func (a *Application) accountsCreate(c *Context) {
	c.action = &flactions.AccountsCreate{}
	routeWithRepository(authenticateUser(loadActionBody(performAction)))(c)
}

func (a *Application) propertiesCreate(c *Context) {
	c.action = &flactions.PropertiesCreate{}
	routeWithRepository(authenticateUser(performAction))(c)
}

func loadActionBody(next func(*Context)) func(*Context) {
	return func(c *Context) {
		if !c.parseBodyOrFail(c.action) {
			return
		}

		next(c)
	}
}

func validateActionBody(next func(*Context)) func(*Context) {
	return func(c *Context) {
		if !c.validate(c.action.(flentities.Validatable)) {
			return
		}

		next(c)
	}
}

func routeWithRepository(next func(*Context)) func(*Context) {
	return func(c *Context) {
		err := flentities.WithDatabase(flconfig.Config, func(db *gorm.DB) error {
			c.repository = &flentities.Repository{DB: db}
			next(c)
			return nil
		})

		if err != nil {
			c.ServerError(err)
		}
	}
}

func performAction(c *Context) {
	c.action.Perform(c)
}
