package flweb

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/gin-gonic/gin"
)

func withDiagnostics(c *gin.Context) {
	start := time.Now()

	contextDiagnostics := &fldiagnostics.Diagnostics{}
	c.Set("diagnostics", contextDiagnostics)

	c.Next()

	end := time.Now()

	diagnostics := (&fldiagnostics.Diagnostics{}).
		Add("at", end.Format("2006-01-02T15:04:05Z")).
		Add("status", strconv.Itoa(c.Writer.Status())).
		Add("duration", end.Sub(start).String()).
		Add("ip", c.ClientIP()).
		Add("method", c.Request.Method).
		AddQuoted("url", c.Request.URL.String()).
		Concat(contextDiagnostics)

	log.Println(diagnostics)
}

type contextFunc func(*Context)

func parseBody(next contextFunc) contextFunc {
	return func(c *Context) {
		if !c.parseBodyOrFail(c.action) {
			return
		}

		c.Diagnostics().AddJSON("body", c.action)

		if v, ok := c.action.(flentities.Validatable); ok {
			errs := v.Validate()
			if len(errs) > 0 {
				c.Diagnostics().AddJSON("validation_errors", c.errorsBody(errs))
				c.Respond(http.StatusUnprocessableEntity, c.errorsBody(errs))
				return
			}
		}

		next(c)
	}
}

func withRepository(next contextFunc) contextFunc {
	return func(c *Context) {
		err := flentities.WithRepository(func(r *flentities.Repository) error {
			c.repository = r
			next(c)
			return nil
		})

		if err != nil {
			c.ServerError(err)
		}
	}
}

func withAction(a flactions.Action, next func(c *Context)) func(c *Context) {
	return func(c *Context) {
		c.action = a
		c.Diagnostics().Add("action", reflect.TypeOf(a).Elem().Name())

		next(c)
	}
}

func performAction(c *Context) {
	c.action.Perform()
}

func requireAccount(next contextFunc) contextFunc {
	return func(c *Context) {
		account := c.CurrentAccount()

		if account == nil {
			c.RespondError(http.StatusPreconditionRequired, errors.New("You need a company account to perform this action"))
			return
		}

		next(c)
	}
}
