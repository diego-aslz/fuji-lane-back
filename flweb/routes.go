package flweb

import (
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flactions"
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
	// CountriesPath for listing countries
	CountriesPath = "/countries"
	// PropertiesPath for properties management
	PropertiesPath = "/properties"
	// PropertiesShowPath to show property details
	PropertiesShowPath = "/properties/:id"
	// PropertiesImagesNewPath for obtaining pre-signed upload URLs
	PropertiesImagesNewPath = "/properties/:id/images/new"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	a.route(e.GET, StatusPath, a.status)

	a.route(e.POST, SignUpPath, a.signUp)
	a.route(e.POST, SignInPath, a.signIn)
	a.route(e.POST, FacebookSignInPath, a.facebookSignIn)

	a.route(e.POST, AccountsPath, a.accountsCreate)
	a.route(e.GET, CountriesPath, a.countriesList)
	a.route(e.POST, PropertiesPath, a.propertiesCreate)
	a.route(e.GET, PropertiesShowPath, a.propertiesShow)
	a.route(e.GET, PropertiesImagesNewPath, a.propertiesImagesNew)
}

type ginMethod func(string, ...gin.HandlerFunc) gin.IRoutes

func (a *Application) route(method ginMethod, path string, next func(*Context)) {
	method(path, func(c *gin.Context) {
		next(&Context{Context: c, now: a.TimeFunc, randSource: a.RandSource})
	})
}

func (a *Application) status(c *Context) {
	withAction(&flactions.Status{}, performAction)(c)
}

func (a *Application) signUp(c *Context) {
	withAction(&flactions.SignUp{},
		withRepository(
			loadActionBody(
				validateActionBody(
					performAction))))(c)
}

func (a *Application) signIn(c *Context) {
	withAction(&flactions.SignIn{},
		withRepository(
			loadActionBody(
				performAction)))(c)
}

func (a *Application) facebookSignIn(c *Context) {
	withAction(flactions.NewFacebookSignIn(a.facebookClient),
		withRepository(
			loadActionBody(
				performAction)))(c)
}

func (a *Application) accountsCreate(c *Context) {
	withAction(&flactions.AccountsCreate{},
		withRepository(
			authenticateUser(
				loadActionBody(
					performAction))))(c)
}

func (a *Application) countriesList(c *Context) {
	withAction(&flactions.CountriesList{},
		withRepository(
			performAction))(c)
}

func (a *Application) propertiesCreate(c *Context) {
	withAction(&flactions.PropertiesCreate{},
		withRepository(
			authenticateUser(
				performAction)))(c)
}

func (a *Application) propertiesShow(c *Context) {
	withAction(&flactions.PropertiesShow{},
		withRepository(
			authenticateUser(
				performAction)))(c)
}

func (a *Application) propertiesImagesNew(c *Context) {
	withAction(flactions.NewPropertiesImagesNew(),
		withRepository(
			authenticateUser(
				performAction)))(c)
}
