package flweb

import (
	"log"
	"strconv"
	"time"

	"github.com/nerde/fuji-lane-back/fldiagnostics"

	"github.com/gin-gonic/gin"
)

func (a *Application) logMiddleware(c *gin.Context) {
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
