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
	// PropertyPath to show property details
	PropertyPath = "/properties/:id"
	// PropertiesImagesPath for obtaining pre-signed upload URLs
	PropertiesImagesPath = "/properties/:id/images"
	// PropertiesImagesUploadedPath for marking an image as uploaded
	PropertiesImagesUploadedPath = "/properties/:property_id/images/:id/uploaded"
	// ImagePath for accessing a specific image
	ImagePath = "/images/:id"
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
	a.route(e.GET, PropertyPath, a.propertiesShow)
	a.route(e.POST, PropertiesImagesPath, a.propertiesImagesCreate)
	a.route(e.PUT, PropertiesImagesUploadedPath, a.propertiesImagesUploaded)
	a.route(e.DELETE, ImagePath, a.imagesDestroy)
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
				performAction)))(c)
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
				requireAccount(
					performAction))))(c)
}

func (a *Application) propertiesShow(c *Context) {
	withAction(&flactions.PropertiesShow{},
		withRepository(
			authenticateUser(
				performAction)))(c)
}

func (a *Application) propertiesImagesCreate(c *Context) {
	withAction(flactions.NewPropertiesImagesCreate(a.S3Service),
		withRepository(
			authenticateUser(
				loadActionBody(
					requireAccount(
						performAction)))))(c)
}

func (a *Application) propertiesImagesUploaded(c *Context) {
	withAction(&flactions.PropertiesImagesUploaded{},
		withRepository(
			authenticateUser(
				performAction)))(c)
}

func (a *Application) imagesDestroy(c *Context) {
	withAction(flactions.NewImagesDestroy(a.S3Service),
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}
