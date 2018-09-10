package fujilane

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	facebookSignInPath = "/sign_in/facebook"

	propertiesPath = "/properties"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, ginAdapt(a.routeStatus))

	e.POST(signUpPath, ginAdapt(a.routeSignUp))
	e.POST(facebookSignInPath, ginAdapt(a.routeFacebookSignIn))

	e.POST(propertiesPath, ginAdapt(a.requireUser(a.routePropertiesCreate)))
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

func (a *routeContext) getHeader(key string) string {
	values := a.context.Request.Header[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (a *routeContext) set(key string, value interface{}) {
	a.context.Set(key, value)
}

func (a *Application) requireUser(next func(*routeContext)) func(*routeContext) {
	return func(c *routeContext) {
		auth := c.getHeader("Authorization")
		if auth == "" {
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		auth = strings.TrimPrefix(auth, "Bearer ")
		session, err := loadSession(auth)
		if err != nil {
			log.Printf("Unable to load session from token %s: %s\n", auth, err.Error())
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		user, err := a.usersRepository.findByEmail(session.Email)
		if err != nil {
			log.Printf("Unable to load user (email: %s): %s\n", session.Email, err.Error())
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		c.set("current-user", user)

		next(c)
	}
}

// ginAdapt wraps an application route with a function that abstracts gin.Context out of the flow so our routes can
// use the routeContext abstraction
func ginAdapt(route func(*routeContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		route(&routeContext{c})
	}
}
