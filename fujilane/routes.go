package fujilane

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	statusPath = "/status"

	signUpPath         = "/sign_up"
	signInPath         = "/sign_in"
	facebookSignInPath = "/sign_in/facebook"

	propertiesPath = "/properties"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, ginAdapt(a.routeStatus))

	e.POST(signUpPath, ginAdapt(a.routeSignUp))
	e.POST(signInPath, ginAdapt(a.routeSignIn))
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
	a.addLogError(err)
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

func (a *routeContext) addLog(key, value string) {
	logs := a.context.GetString("log-details")
	if len(logs) > 0 {
		logs += " "
	}
	a.context.Set("log-details", logs+key+"="+value)
}

func (a *routeContext) addLogQuoted(key, value string) {
	a.addLog(key, "\""+value+"\"")
}

func (a *routeContext) addLogError(err error) {
	a.addLogQuoted("error", err.Error())
}

func (a *routeContext) addLogJSON(key string, value interface{}) {
	jsonObj, err := json.Marshal(value)
	if err == nil {
		a.addLog(key, string(jsonObj))
	}
}

func (a *routeContext) currentUser() *User {
	v, _ := a.context.Get("current-user")
	return v.(*User)
}

func (a *Application) requireUser(next func(*routeContext)) func(*routeContext) {
	return func(c *routeContext) {
		auth := c.getHeader("Authorization")
		if auth == "" {
			c.addLogQuoted("reason", "Missing authentication token")
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		auth = strings.TrimPrefix(auth, "Bearer ")
		session, err := loadSession(auth)
		if err != nil {
			c.addLogJSON("token", auth)
			c.addLogQuoted("reason", "Unable to load session from token")
			c.addLogError(err)
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		user, err := a.usersRepository.findByEmail(session.Email)
		if err != nil {
			safeSession := *session
			safeSession.Token = "[FILTERED]"
			c.addLogJSON("session", safeSession)
			c.addLogQuoted("reason", "Unable to load user")
			c.addLogError(err)
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		if user == nil || user.ID == 0 {
			safeSession := *session
			safeSession.Token = "[FILTERED]"
			c.addLogJSON("session", safeSession)
			c.addLogQuoted("reason", "User not found")
			c.fail(http.StatusUnauthorized, errors.New("You need to sign in"))
			return
		}

		c.addLog("user", user.Email)
		c.addLog("user_id", fmt.Sprint(user.ID))
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
