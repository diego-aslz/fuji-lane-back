package flweb

import (
	"log"
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
		Add("path", c.Request.URL.Path).
		Concat(contextDiagnostics)

	log.Println(diagnostics)
}

func loadActionBody(next func(*Context)) func(*Context) {
	return func(c *Context) {
		if !c.parseBodyOrFail(c.action) {
			return
		}

		c.Diagnostics().AddJSON("body", c.action)

		if validatable, ok := c.action.(flentities.Validatable); ok {
			if !c.validate(validatable) {
				return
			}
		}

		next(c)
	}
}

func withRepository(next func(*Context)) func(*Context) {
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
	c.action.Perform(c)
}
