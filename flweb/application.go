package flweb

import (
	"math/rand"
	"time"

	"github.com/nerde/fuji-lane-back/flutils"

	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flservices"
)

// Application is the struct that represents a Fuji Lane app
type Application struct {
	facebookClient flservices.FacebookClient
	TimeFunc       func() time.Time
	RandSource     rand.Source
	flservices.S3Service
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
func NewApplication(facebookClient flservices.FacebookClient) *Application {
	s3, err := flservices.NewS3()
	if err != nil {
		panic(err)
	}

	return &Application{
		facebookClient: facebookClient,
		TimeFunc:       time.Now,
		RandSource:     flutils.NewRandomSource(),
		S3Service:      s3,
	}
}
