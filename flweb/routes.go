package flweb

import (
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flactions"
)

const (
	// StatusPath path for health check
	StatusPath = "/status"

	// FacebookSignInPath for Facebook sign in
	FacebookSignInPath = "/sign_in/facebook"
	// RenewSessionPath use used to get a new session and a new token
	RenewSessionPath = "/renew_session"
	// SignInPath for email sign in
	SignInPath = "/sign_in"
	// SignUpPath for email sign up
	SignUpPath = "/sign_up"

	// AccountsPath for accounts management
	AccountsPath = "/accounts"
	// AmenityTypesPath for listing amenity types
	AmenityTypesPath = "/amenity_types/:target"
	// CitiesPath for listing cities
	CitiesPath = "/cities"
	// CountriesPath for listing countries
	CountriesPath = "/countries"
	// DashboardPath for getting dashboard details
	DashboardPath = "/dashboard"
	// ImagePath to access a specific image
	ImagePath = "/images/:id"
	// ImagesSortPath to access images
	ImagesSortPath = "/images/sort"
	// ImagesPath to access images
	ImagesPath = "/images"
	// ImagesUploadedPath for marking an image as uploaded
	ImagesUploadedPath = "/images/:id/uploaded"
	// PropertiesPath to access properties
	PropertiesPath = "/properties"
	// PropertiesPublishPath to publish a property
	PropertiesPublishPath = "/properties/:id/publish"
	// PropertyPath to access a specific property
	PropertyPath = "/properties/:id"
	// UnitPath to access a specific unit
	UnitPath = "/units/:id"
	// UnitsPath to access units
	UnitsPath = "/units"
	// UnitsPublishPath to publish a unit
	UnitsPublishPath = "/units/:id/publish"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	a.route(e.GET, StatusPath, a.status)

	a.route(e.GET, RenewSessionPath, a.renewSession)
	a.route(e.POST, FacebookSignInPath, a.facebookSignIn)
	a.route(e.POST, SignInPath, a.signIn)
	a.route(e.POST, SignUpPath, a.signUp)

	a.route(e.GET, CitiesPath, a.citiesList)
	a.route(e.GET, CountriesPath, a.countriesList)

	a.route(e.POST, AccountsPath, a.accountsCreate)

	a.route(e.GET, AmenityTypesPath, a.amenityTypesList)

	a.route(e.GET, PropertyPath, a.propertiesShow)
	a.route(e.GET, PropertiesPath, a.propertiesList)
	a.route(e.POST, PropertiesPath, a.propertiesCreate)
	a.route(e.PUT, PropertyPath, a.propertiesUpdate)
	a.route(e.PUT, PropertiesPublishPath, a.propertiesPublish)

	a.route(e.DELETE, ImagePath, a.imagesDestroy)
	a.route(e.POST, ImagesPath, a.imagesCreate)
	a.route(e.POST, ImagesSortPath, a.imagesSort)
	a.route(e.PUT, ImagesUploadedPath, a.imagesUploaded)

	a.route(e.GET, UnitPath, a.unitsShow)
	a.route(e.POST, UnitsPath, a.unitsCreate)
	a.route(e.PUT, UnitPath, a.unitsUpdate)
	a.route(e.PUT, UnitsPublishPath, a.unitsPublish)

	a.route(e.GET, DashboardPath, a.dashboard)
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

func (a *Application) renewSession(c *Context) {
	withAction(&flactions.RenewSession{},
		withRepository(
			authenticateUser(
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

func (a *Application) amenityTypesList(c *Context) {
	withAction(&flactions.AmenityTypesList{}, performAction)(c)
}

func (a *Application) citiesList(c *Context) {
	withAction(&flactions.CitiesList{},
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

func (a *Application) propertiesUpdate(c *Context) {
	withAction(&flactions.PropertiesUpdate{},
		withRepository(
			authenticateUser(
				requireAccount(
					loadActionBody(
						performAction)))))(c)
}

func (a *Application) propertiesPublish(c *Context) {
	withAction(&flactions.PropertiesPublish{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) propertiesShow(c *Context) {
	withAction(&flactions.PropertiesShow{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) propertiesList(c *Context) {
	withAction(&flactions.PropertiesList{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) imagesCreate(c *Context) {
	withAction(flactions.NewImagesCreate(a.S3Service),
		withRepository(
			authenticateUser(
				loadActionBody(
					requireAccount(
						performAction)))))(c)
}

func (a *Application) unitsCreate(c *Context) {
	withAction(&flactions.UnitsCreate{},
		withRepository(
			authenticateUser(
				loadActionBody(
					performAction))))(c)
}

func (a *Application) unitsUpdate(c *Context) {
	withAction(&flactions.UnitsUpdate{},
		withRepository(
			authenticateUser(
				loadActionBody(
					performAction))))(c)
}

func (a *Application) unitsPublish(c *Context) {
	withAction(&flactions.UnitsPublish{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) unitsShow(c *Context) {
	withAction(&flactions.UnitsShow{},
		withRepository(
			authenticateUser(
				performAction)))(c)
}

func (a *Application) imagesSort(c *Context) {
	withAction(&flactions.ImagesSort{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) imagesUploaded(c *Context) {
	withAction(&flactions.ImagesUploaded{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) imagesDestroy(c *Context) {
	withAction(flactions.NewImagesDestroy(a.S3Service),
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}

func (a *Application) dashboard(c *Context) {
	withAction(&flactions.Dashboard{},
		withRepository(
			authenticateUser(
				requireAccount(
					performAction))))(c)
}
