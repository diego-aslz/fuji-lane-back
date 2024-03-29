package flweb

import (
	"math/rand"
	"time"

	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/nerde/fuji-lane-back/fljobs"

	"github.com/nerde/fuji-lane-back/fujilane"

	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flservices"
)

// Application is the struct that represents a Fuji Lane app
type Application struct {
	FacebookClient flservices.FacebookClient
	TimeFunc       func() time.Time
	RandSource     rand.Source
	Mailer         flservices.Mailer
	S3Service      flservices.S3Service
	Sendgrid       flservices.Sendgrid
	GoogleAuth     flservices.GoogleAuth
	Jobs           *fljobs.Application
}

// Start listening to requests
func (a *Application) Start() {
	a.CreateRouter().Run()
}

// CreateRouter with all the recognized paths and their handlers
func (a *Application) CreateRouter() *gin.Engine {
	r := gin.New()
	r.Use(withDiagnostics)
	r.Use(recovery)
	a.AddRoutes(r)

	return r
}

// NewDefaultApplication returns a new web application
func NewDefaultApplication() (*Application, error) {
	s3, err := flservices.NewS3()
	if err != nil {
		return nil, err
	}

	return &Application{
		FacebookClient: flservices.NewFacebookHTTPClient(),
		GoogleAuth:     flservices.NewGoogleAuth(flconfig.Config.GoogleAuthAppID),
		TimeFunc:       time.Now,
		RandSource:     fujilane.NewRandomSource(),
		S3Service:      s3,
		Sendgrid:       flservices.NewSendgridAPI(),
		Mailer:         flservices.NewSMTPMailer(),
		Jobs:           fljobs.NewDefaultApplication(),
	}, nil
}
