package flweb

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flactions"
)

// Application is the struct that represents a Fuji Lane app
type Application struct {
	facebookClient flactions.FacebookClient
	TimeFunc       func() time.Time
}

// Start listening to requests
func (a *Application) Start() {
	a.CreateRouter().Run()
}

// CreateRouter with all the recognized paths and their handlers
func (a *Application) CreateRouter() *gin.Engine {
	r := gin.New()
	r.Use(withDiagnostics)
	r.Use(gin.Recovery())
	a.AddRoutes(r)
	return r
}

// NewApplication with the injected dependencies
func NewApplication(facebookClient flactions.FacebookClient) *Application {
	return &Application{
		facebookClient: facebookClient,
		TimeFunc:       time.Now,
	}
}
