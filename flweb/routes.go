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
	// ListingPath to get listing details
	ListingPath = "/listings/:id"
	// ProfilePath to publish a unit
	ProfilePath = "/profile"
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

	a.route(e.GET, RenewSessionPath, a.renewSession, withRepository, loadSession, requireUser)
	a.route(e.POST, FacebookSignInPath, a.facebookSignIn, withRepository, loadActionBody)
	a.route(e.POST, SignInPath, a.signIn, withRepository, loadActionBody)
	a.route(e.POST, SignUpPath, a.signUp, withRepository, loadActionBody)

	a.route(e.GET, CitiesPath, a.citiesList, withRepository)
	a.route(e.GET, CountriesPath, a.countriesList, withRepository)

	a.route(e.POST, AccountsPath, a.accountsCreate, withRepository, loadSession, requireUser, loadActionBody)

	a.route(e.GET, AmenityTypesPath, a.amenityTypesList)

	a.route(e.GET, PropertyPath, a.propertiesShow, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.GET, PropertiesPath, a.propertiesList, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.POST, PropertiesPath, a.propertiesCreate, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.PUT, PropertyPath, a.propertiesUpdate, withRepository, loadSession, requireUser, requireAccount, loadActionBody)
	a.route(e.PUT, PropertiesPublishPath, a.propertiesPublish, withRepository, loadSession, requireUser, requireAccount)

	a.route(e.DELETE, ImagePath, a.imagesDestroy, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.POST, ImagesPath, a.imagesCreate, withRepository, loadSession, requireUser, loadActionBody, requireAccount)
	a.route(e.POST, ImagesSortPath, a.imagesSort, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.PUT, ImagesUploadedPath, a.imagesUploaded, withRepository, loadSession, requireUser, requireAccount)

	a.route(e.GET, ListingPath, a.listingsShow, withRepository, loadSession)

	a.route(e.PUT, ProfilePath, a.profileUpdate, withRepository, loadSession, requireUser, loadActionBody)

	a.route(e.GET, UnitPath, a.unitsShow, withRepository, loadSession, requireUser)
	a.route(e.POST, UnitsPath, a.unitsCreate, withRepository, loadSession, requireUser, loadActionBody)
	a.route(e.PUT, UnitPath, a.unitsUpdate, withRepository, loadSession, requireUser, loadActionBody)
	a.route(e.PUT, UnitsPublishPath, a.unitsPublish, withRepository, loadSession, requireUser, requireAccount)

	a.route(e.GET, DashboardPath, a.dashboard, withRepository, loadSession, requireUser, requireAccount)
}

type ginMethod func(string, ...gin.HandlerFunc) gin.IRoutes

func (a *Application) route(method ginMethod, path string, actionProvider func() flactions.Action,
	middleware ...func(contextFunc) contextFunc) {

	next := combineMiddleware(middleware...)

	method(path, func(c *gin.Context) {
		withAction(actionProvider(), next)(a.newContext(c))
	})
}

func (a *Application) newContext(c *gin.Context) *Context {
	return &Context{Context: c, now: a.TimeFunc, randSource: a.RandSource}
}

func (a *Application) status() flactions.Action {
	return &flactions.Status{}
}

func (a *Application) signUp() flactions.Action {
	return &flactions.SignUp{}
}

func (a *Application) signIn() flactions.Action {
	return &flactions.SignIn{}
}

func (a *Application) facebookSignIn() flactions.Action {
	return flactions.NewFacebookSignIn(a.facebookClient)
}

func (a *Application) renewSession() flactions.Action {
	return &flactions.RenewSession{}
}

func (a *Application) accountsCreate() flactions.Action {
	return &flactions.AccountsCreate{}
}

func (a *Application) countriesList() flactions.Action {
	return &flactions.CountriesList{}
}

func (a *Application) amenityTypesList() flactions.Action {
	return &flactions.AmenityTypesList{}
}

func (a *Application) citiesList() flactions.Action {
	return &flactions.CitiesList{}
}

func (a *Application) profileUpdate() flactions.Action {
	return &flactions.ProfileUpdate{}
}

func (a *Application) propertiesCreate() flactions.Action {
	return &flactions.PropertiesCreate{}
}

func (a *Application) propertiesUpdate() flactions.Action {
	return &flactions.PropertiesUpdate{}
}

func (a *Application) propertiesPublish() flactions.Action {
	return &flactions.PropertiesPublish{}
}

func (a *Application) propertiesShow() flactions.Action {
	return &flactions.PropertiesShow{}
}

func (a *Application) listingsShow() flactions.Action {
	return &flactions.ListingsShow{}
}

func (a *Application) propertiesList() flactions.Action {
	return &flactions.PropertiesList{}
}

func (a *Application) imagesCreate() flactions.Action {
	return flactions.NewImagesCreate(a.S3Service)
}

func (a *Application) unitsCreate() flactions.Action {
	return &flactions.UnitsCreate{}
}

func (a *Application) unitsUpdate() flactions.Action {
	return &flactions.UnitsUpdate{}
}

func (a *Application) unitsPublish() flactions.Action {
	return &flactions.UnitsPublish{}
}

func (a *Application) unitsShow() flactions.Action {
	return &flactions.UnitsShow{}
}

func (a *Application) imagesSort() flactions.Action {
	return &flactions.ImagesSort{}
}

func (a *Application) imagesUploaded() flactions.Action {
	return &flactions.ImagesUploaded{}
}

func (a *Application) imagesDestroy() flactions.Action {
	return flactions.NewImagesDestroy(a.S3Service)
}

func (a *Application) dashboard() flactions.Action {
	return &flactions.Dashboard{}
}

func combineMiddleware(middleware ...func(contextFunc) contextFunc) contextFunc {
	l := len(middleware)

	if l == 0 {
		return performAction
	}

	if l == 1 {
		return middleware[0](performAction)
	}

	return middleware[0](combineMiddleware(middleware[1:l]...))
}
