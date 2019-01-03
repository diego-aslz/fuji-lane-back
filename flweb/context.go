package flweb

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

// Context is a thin abstraction layer around gin.Context so our routes don't directly depend on it and we can
// switch web libraries with less pain if we ever need to
type Context struct {
	*gin.Context
	randSource     rand.Source
	repository     *flentities.Repository
	now            func() time.Time
	action         flactions.Action
	session        *flentities.Session
	currentUser    *flentities.User
	currentAccount *flentities.Account
}

// Diagnostics returns the Diagnostics object being used for reporting execution details
func (c *Context) Diagnostics() *fldiagnostics.Diagnostics {
	d, _ := c.Get("diagnostics")
	return d.(*fldiagnostics.Diagnostics)
}

// Now returns the current time and can be injected
func (c *Context) Now() time.Time {
	return c.now()
}

// RandomSource returns the current source for randomness
func (c *Context) RandomSource() rand.Source {
	return c.randSource
}

// Respond responds with the given status and body in JSON format
func (c *Context) Respond(status int, body interface{}) {
	c.JSON(status, body)
}

// RespondNotFound returns a default Not Found error with Not Found status code
func (c *Context) RespondNotFound() {
	c.RespondError(http.StatusNotFound, errors.New("Not Found"))
}

// RespondError creates an error response with the given error
func (c *Context) RespondError(status int, errs ...error) {
	messages := []string{}
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	c.Diagnostics().AddJSON("response_errors", messages)
	c.JSON(status, c.errorsBody(errs))
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

// parseBodyOrFail will try to parse the body as JSON and fail with BAD_REQUEST if an error is returned
func (c *Context) parseBodyOrFail(dst interface{}) bool {
	err := c.BindJSON(dst)
	if err != nil {
		c.RespondError(http.StatusBadRequest, err)
	}
	return err == nil
}

func (c *Context) getHeader(key string) string {
	values := c.Request.Header[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Repository returns the current Repository for database access
func (c *Context) Repository() *flentities.Repository {
	return c.repository
}
