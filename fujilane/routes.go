package fujilane

import (
	"errors"
	"net/http"
	"time"

	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/fldiagnostics"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	signInPath         = "/sign_in"
	facebookSignInPath = "/sign_in/facebook"

	accountsPath   = "/accounts"
	propertiesPath = "/properties"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, a.ginAdapt(routeAction(flactions.NewStatus)))

	e.POST(signUpPath, a.ginAdapt(routeWithRepository(routeActionWithValidatableBody(flactions.NewSignUp))))
	e.POST(signInPath, a.ginAdapt(routeWithRepository(routeActionWithBody(flactions.NewSignIn))))
	e.POST(facebookSignInPath, a.ginAdapt(routeWithRepository(routeActionWithBody(a.createFacebookSignIn))))

	e.POST(accountsPath,
		a.ginAdapt(routeWithRepository(authenticateUser(routeActionWithBody(flactions.NewAccountsCreate)))))
	e.POST(propertiesPath,
		a.ginAdapt(routeWithRepository(authenticateUser(routeAction(flactions.NewPropertiesCreate)))))
}

// ginAdapt wraps an application route with a function that abstracts gin.Context out of the flow so our routes can
// use the routeContext abstraction
func (a *Application) ginAdapt(route func(*Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		route(&Context{
			context: c,
			now:     a.timeFunc,
		})
	}
}

func (a *Application) createFacebookSignIn() flactions.Action {
	return flactions.NewFacebookSignIn(a.facebookClient)
}

// Context is a thin abstraction layer around gin.Context so our routes don't directly depend on it and we can
// switch web libraries with less pain if we ever need to
type Context struct {
	context    *gin.Context
	repository *flentities.Repository
	now        func() time.Time
}

// Diagnostics returns the Diagnostics object being used for reporting execution details
func (c *Context) Diagnostics() *fldiagnostics.Diagnostics {
	d, _ := c.context.Get("diagnostics")
	return d.(*fldiagnostics.Diagnostics)
}

// Now returns the current time and can be injected
func (c *Context) Now() time.Time {
	return c.now()
}

// Respond responds with the given status and body in JSON format
func (c *Context) Respond(status int, body interface{}) {
	c.context.JSON(status, body)
}

// RespondError creates an error response with the given error
func (c *Context) RespondError(status int, err error) {
	c.Diagnostics().AddQuoted("response_error", err.Error())
	c.context.JSON(status, c.errorsBody([]error{err}))
}

func (c *Context) errorsBody(errs []error) map[string]interface{} {
	messages := []string{}
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return map[string]interface{}{"errors": messages}
}

// ServerError adds the error to Diagnostics and responds with 500 status and a generic error message
func (c *Context) ServerError(err error) {
	c.Diagnostics().AddError(err)
	c.RespondError(http.StatusInternalServerError, errors.New("Sorry, something went wrong"))
}

func (c *Context) parseBodyAndValidate(dst flentities.Validatable) bool {
	return c.parseBodyOrFail(dst) && c.validate(dst)
}

func (c *Context) validate(v flentities.Validatable) bool {
	errs := v.Validate()
	if len(errs) > 0 {
		c.Respond(http.StatusUnprocessableEntity, c.errorsBody(errs))
		return false
	}

	return true
}

// parseBodyOrFail will try to parse the body as JSON and fail with BAD_REQUEST if an error is returned
func (c *Context) parseBodyOrFail(dst interface{}) bool {
	err := c.context.BindJSON(dst)
	if err != nil {
		c.RespondError(http.StatusBadRequest, err)
	}
	return err == nil
}

func (c *Context) getHeader(key string) string {
	values := c.context.Request.Header[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *Context) set(key string, value interface{}) {
	c.context.Set(key, value)
}

// Repository returns the current Repository for database access
func (c *Context) Repository() *flentities.Repository {
	return c.repository
}

func routeActionWithBody(creator flactions.ActionCreator) func(*Context) {
	return func(c *Context) {
		action := creator()
		if !c.parseBodyOrFail(action) {
			return
		}

		action.Perform(c)
	}
}

func routeActionWithValidatableBody(creator flactions.ActionCreator) func(*Context) {
	return func(c *Context) {
		action := creator()
		if !c.parseBodyAndValidate(action.(flentities.Validatable)) {
			return
		}

		action.Perform(c)
	}
}

func routeAction(creator flactions.ActionCreator) func(*Context) {
	return func(c *Context) {
		creator().Perform(c)
	}
}

func routeWithRepository(next func(*Context)) func(*Context) {
	return func(c *Context) {
		err := flentities.WithDatabase(flconfig.Config, func(db *gorm.DB) error {
			c.repository = &flentities.Repository{DB: db}
			next(c)
			return nil
		})

		if err != nil {
			c.ServerError(err)
		}
	}
}
