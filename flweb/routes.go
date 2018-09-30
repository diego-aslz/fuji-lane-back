package flweb

import (
	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
)

const (
	// StatusPath path for health check
	StatusPath = "/status"

	// SignUpPath for email sign up
	SignUpPath = "/sign_up"
	// SignInPath for email sign in
	SignInPath = "/sign_in"
	// FacebookSignInPath for Facebook sign in
	FacebookSignInPath = "/sign_in/facebook"

	// AccountsPath for accounts management
	AccountsPath = "/accounts"
	// PropertiesPath for properties management
	PropertiesPath = "/properties"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	a.route(e.GET, StatusPath, a.status)

	a.route(e.POST, SignUpPath, a.signUp)
	a.route(e.POST, SignInPath, a.signIn)
	a.route(e.POST, FacebookSignInPath, a.facebookSignIn)

	a.route(e.POST, AccountsPath, a.accountsCreate)
	a.route(e.POST, PropertiesPath, a.propertiesCreate)
}

type ginMethod func(string, ...gin.HandlerFunc) gin.IRoutes

func (a *Application) route(method ginMethod, path string, next func(*Context)) {
	method(path, func(c *gin.Context) {
		next(&Context{context: c, now: a.TimeFunc})
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
