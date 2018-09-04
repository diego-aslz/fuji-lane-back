package fujilane

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	facebookSignInPath = "/sign_in/facebook"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, ginAdapt(a.routeStatus))

	e.POST(signUpPath, ginAdapt(a.routeSignUp))
	e.POST(facebookSignInPath, ginAdapt(a.routeFacebookSignIn))
}

// routeContext is a thin abstraction layer around gin.Context so our routes don't directly depend on it and we can
// switch web libraries with less pain if we ever need to
type routeContext struct {
	context *gin.Context
}

// respond responds with the given status and body in JSON format
func (a *routeContext) respond(status int, body interface{}) {
	a.context.JSON(status, body)
}

func (a *routeContext) errorsBody(errs []error) map[string]interface{} {
	messages := []string{}
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return map[string]interface{}{"errors": messages}
}

func (a *routeContext) fail(status int, err error) {
	a.context.AbortWithError(status, err)
}

func (a *routeContext) parseBodyAndValidate(dst Validatable) bool {
	return a.parseBodyOrFail(dst) && a.validate(dst)
}

func (a *routeContext) validate(v Validatable) bool {
	errs := v.Validate()
	if len(errs) > 0 {
		a.respond(http.StatusUnprocessableEntity, a.errorsBody(errs))
		return false
	}

	return true
}

// parseBodyOrFail will try to parse the body as JSON and fail with BAD_REQUEST if an error is returned
func (a *routeContext) parseBodyOrFail(dst interface{}) bool {
	err := a.context.BindJSON(dst)
	if err != nil {
		a.fail(http.StatusBadRequest, err)
	}
	return err == nil
}

// ginAdapt wraps an application route with a function that abstracts gin.Context out of the flow so our routes can
// use the routeContext abstraction
func ginAdapt(route func(*routeContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		route(&routeContext{c})
	}
}
