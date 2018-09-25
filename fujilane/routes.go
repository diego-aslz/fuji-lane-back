package fujilane

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	e.GET(statusPath, ginAdapt(a.routeStatus))

	e.POST(signUpPath, ginAdapt(a.routeSignUp))
	e.POST(signInPath, ginAdapt(a.routeSignIn))
	e.POST(facebookSignInPath, ginAdapt(a.routeFacebookSignIn))

	e.POST(accountsPath, ginAdapt(a.authenticateUser(a.routeAccountsCreate)))
	e.POST(propertiesPath, ginAdapt(a.authenticateUser(a.routePropertiesCreate)))
}

// routeContext is a thin abstraction layer around gin.Context so our routes don't directly depend on it and we can
// switch web libraries with less pain if we ever need to
type routeContext struct {
	context *gin.Context
}

// respond responds with the given status and body in JSON format
func (c *routeContext) respond(status int, body interface{}) {
	c.context.JSON(status, body)
}

func (c *routeContext) respondError(status int, err error) {
	c.addLogQuoted("response_error", err.Error())
	c.context.JSON(status, c.errorsBody([]error{err}))
}

func (c *routeContext) errorsBody(errs []error) map[string]interface{} {
	messages := []string{}
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return map[string]interface{}{"errors": messages}
}

func (c *routeContext) fatal(err error) {
	c.addLogError(err)
	c.respondError(http.StatusInternalServerError, errors.New("Sorry, something went wrong"))
}

func (c *routeContext) parseBodyAndValidate(dst Validatable) bool {
	return c.parseBodyOrFail(dst) && c.validate(dst)
}

func (c *routeContext) validate(v Validatable) bool {
	errs := v.Validate()
	if len(errs) > 0 {
		c.respond(http.StatusUnprocessableEntity, c.errorsBody(errs))
		return false
	}

	return true
}

// parseBodyOrFail will try to parse the body as JSON and fail with BAD_REQUEST if an error is returned
func (c *routeContext) parseBodyOrFail(dst interface{}) bool {
	err := c.context.BindJSON(dst)
	if err != nil {
		c.respondError(http.StatusBadRequest, err)
	}
	return err == nil
}

func (c *routeContext) getHeader(key string) string {
	values := c.context.Request.Header[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *routeContext) set(key string, value interface{}) {
	c.context.Set(key, value)
}

// ginAdapt wraps an application route with a function that abstracts gin.Context out of the flow so our routes can
// use the routeContext abstraction
func ginAdapt(route func(*routeContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		route(&routeContext{c})
	}
}
