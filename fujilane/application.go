package fujilane

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Application is the struct that represents a Fuji Lane app
type Application struct {
	facebook        *facebook
	usersRepository *usersRepository
	timeFunc        func() time.Time
}

// Start listening to requests
func (a *Application) Start() {
	a.CreateRouter().Run()
}

// CreateRouter with all the recognized paths and their handlers
func (a *Application) CreateRouter() *gin.Engine {
	r := gin.Default()
	a.AddRoutes(r)
	return r
}

// NewApplication with the injected dependencies
func NewApplication(facebookClient FacebookClient) *Application {
	return &Application{
		facebook:        newFacebook(facebookClient),
		usersRepository: &usersRepository{},
		timeFunc:        time.Now,
	}
}
