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
	randSource rand.Source
	repository *flentities.Repository
	now        func() time.Time
	action     flactions.Action
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
func (c *Context) RespondError(status int, err error) {
	c.Diagnostics().AddQuoted("response_error", err.Error())
	c.JSON(status, c.errorsBody([]error{err}))
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

func (c *Context) validate(v flentities.Validatable) bool {
	errs := v.Validate()
	if len(errs) > 0 {
		c.Diagnostics().AddJSON("validation_errors", c.errorsBody(errs))
		c.Respond(http.StatusUnprocessableEntity, c.errorsBody(errs))
		return false
	}

	return true
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

func (c *Context) set(key string, value interface{}) {
	c.Set(key, value)
}

// Repository returns the current Repository for database access
func (c *Context) Repository() *flentities.Repository {
	return c.repository
}
