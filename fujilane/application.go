package fujilane

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Application is the struct that represents a Fuji Lane app
type Application struct {
	facebook             *facebook
	usersRepository      *usersRepository
	propertiesRepository *propertiesRepository
	timeFunc             func() time.Time
}

// Start listening to requests
func (a *Application) Start() {
	a.CreateRouter().Run()
}

// CreateRouter with all the recognized paths and their handlers
func (a *Application) CreateRouter() *gin.Engine {
	r := gin.New()
	r.Use(a.ginLogger)
	r.Use(gin.Recovery())
	a.AddRoutes(r)
	return r
}

func (a *Application) ginLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	end := time.Now()
	duration := end.Sub(start)
	method := c.Request.Method
	statusCode := c.Writer.Status()

	log.Printf("at=%v status=%d duration=%v ip=%s method=%s path=%s %s\n", end.Format("2006-01-02T15:04:05Z"),
		statusCode, duration, c.ClientIP(), method, path, c.GetString("log-details"))
}

// NewApplication with the injected dependencies
func NewApplication(facebookClient FacebookClient) *Application {
	return &Application{
		facebook:             newFacebook(facebookClient),
		usersRepository:      &usersRepository{},
		propertiesRepository: &propertiesRepository{},
		timeFunc:             time.Now,
	}
}
