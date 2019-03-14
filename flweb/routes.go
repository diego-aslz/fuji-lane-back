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
	// BookingsPath for accessing user bookings
	BookingsPath = "/bookings"
	// CitiesPath for listing cities
	CitiesPath = "/cities"
	// CountriesPath for listing countries
	CountriesPath = "/countries"
	// DashboardPath for getting dashboard details
	DashboardPath = "/dashboard"
	// DashboardBookingsPath for accessing my properties' bookings
	DashboardBookingsPath = "/dashboard/bookings"
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
	// NewsletterSubscribePath to subscribe to our newsletter
	NewsletterSubscribePath = "/newsletter/subscribe"
	// ProfilePath to publish a unit
	ProfilePath = "/profile"
	// PropertiesPath to access properties
	PropertiesPath = "/properties"
	// PropertiesPublishPath to publish a property
	PropertiesPublishPath = "/properties/:id/publish"
	// PropertiesSitemapPath to access properties
	PropertiesSitemapPath = "/properties_sitemap.xml"
	// PropertiesUnpublishPath to unpublish a property
	PropertiesUnpublishPath = "/properties/:id/unpublish"
	// PropertyPath to access a specific property
	PropertyPath = "/properties/:id"
	// SearchPath to search for listings
	SearchPath = "/search"
	// UnitPath to access a specific unit
	UnitPath = "/units/:id"
	// UnitsPath to access units
	UnitsPath = "/units"
	// UnitsPublishPath to publish a unit
	UnitsPublishPath = "/units/:id/publish"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	a.route(e.GET, StatusPath, flactions.NewStatus)

	a.route(e.POST, NewsletterSubscribePath, a.newsletterSubscribe, parseBody)

	a.route(e.GET, AmenityTypesPath, flactions.NewAmenityTypesList)
	a.route(e.GET, BookingsPath, flactions.NewBookingsList, withRepository, loadSession, requireUser)
	a.route(e.POST, BookingsPath, a.bookingsCreate, withRepository, loadSession, requireUser, parseBody)
	a.route(e.GET, CitiesPath, flactions.NewCitiesList, withRepository)
	a.route(e.GET, CountriesPath, flactions.NewCountriesList, withRepository)

	a.route(e.GET, DashboardPath, flactions.NewDashboard, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.GET, DashboardBookingsPath, flactions.NewDashboardBookings, withRepository, loadSession, requireUser,
		requireAccount)

	a.route(e.GET, SearchPath, flactions.NewSearch, withRepository)

	a.routeAuthentication(e)
	a.routeImages(e)
	a.routeProperties(e)
	a.routeUnits(e)
}

func (a *Application) routeAuthentication(e *gin.Engine) {
	a.route(e.POST, AccountsPath, flactions.NewAccountsCreate, withRepository, loadSession, requireUser, parseBody)
	a.route(e.GET, ProfilePath, flactions.NewProfileShow, withRepository, loadSession, requireUser)
	a.route(e.PUT, ProfilePath, flactions.NewProfileUpdate, withRepository, loadSession, requireUser, parseBody)

	a.route(e.GET, RenewSessionPath, flactions.NewSessionsRenew, withRepository, loadSession, requireUser)
	a.route(e.POST, FacebookSignInPath, a.facebookSignIn, withRepository, parseBody)
	a.route(e.POST, SignInPath, flactions.NewSignIn, withRepository, parseBody)
	a.route(e.POST, SignUpPath, flactions.NewSignUp, withRepository, parseBody)
}

func (a *Application) routeImages(e *gin.Engine) {
	a.route(e.DELETE, ImagePath, a.imagesDestroy, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.POST, ImagesPath, a.imagesCreate, withRepository, loadSession, requireUser, parseBody, requireAccount)
	a.route(e.POST, ImagesSortPath, flactions.NewImagesSort, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.PUT, ImagesUploadedPath, flactions.NewImagesUploaded, withRepository, loadSession, requireUser, requireAccount)
}

func (a *Application) routeUnits(e *gin.Engine) {
	a.route(e.GET, UnitPath, flactions.NewUnitsShow, withRepository, loadSession, requireUser)
	a.route(e.POST, UnitsPath, flactions.NewUnitsCreate, withRepository, loadSession, requireUser, parseBody)
	a.route(e.PUT, UnitPath, flactions.NewUnitsUpdate, withRepository, loadSession, requireUser, parseBody)
	a.route(e.PUT, UnitsPublishPath, flactions.NewUnitsPublish, withRepository, loadSession, requireUser, requireAccount)
}

func (a *Application) routeProperties(e *gin.Engine) {
	a.route(e.GET, ListingPath, flactions.NewListingsShow, withRepository, loadSession)

	a.route(e.GET, PropertyPath, flactions.NewPropertiesShow, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.GET, PropertiesPath, flactions.NewPropertiesList, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.POST, PropertiesPath, flactions.NewPropertiesCreate, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.PUT, PropertyPath, flactions.NewPropertiesUpdate, withRepository, loadSession, requireUser, requireAccount, parseBody)
	a.route(e.PUT, PropertiesPublishPath, flactions.NewPropertiesPublish, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.PUT, PropertiesUnpublishPath, flactions.NewPropertiesUnpublish, withRepository, loadSession, requireUser, requireAccount)
	a.route(e.GET, PropertiesSitemapPath, flactions.NewPropertiesSitemap, withRepository)
}

type ginMethod func(string, ...gin.HandlerFunc) gin.IRoutes

func (a *Application) route(method ginMethod, path string, actionProvider flactions.Provider,
	middleware ...func(contextFunc) contextFunc) {

	next := combineMiddleware(middleware...)

	method(path, func(c *gin.Context) {
		ctx := a.newContext(c)
		withAction(actionProvider(ctx), next)(ctx)
	})
}

func (a *Application) newContext(c *gin.Context) *Context {
	return &Context{Context: c, now: a.TimeFunc, randSource: a.RandSource}
}

func (a *Application) facebookSignIn(c flactions.Context) flactions.Action {
	return flactions.NewFacebookSignIn(a.FacebookClient, c)
}

func (a *Application) imagesCreate(c flactions.Context) flactions.Action {
	return flactions.NewImagesCreate(a.S3Service, c)
}

func (a *Application) imagesDestroy(c flactions.Context) flactions.Action {
	return flactions.NewImagesDestroy(a.S3Service, c)
}

func (a *Application) bookingsCreate(c flactions.Context) flactions.Action {
	return flactions.NewBookingsCreate(c, a.Mailer)
}

func (a *Application) newsletterSubscribe(c flactions.Context) flactions.Action {
	return flactions.NewNewsletterSubscribe(c, a.Sendgrid)
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
