package flweb

import (
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

func recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				NotifyError(err, c)
			}

			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	c.Next()
}

// NotifyError sends an error to our bug tracking system for debugging
func NotifyError(err error, c *gin.Context) {
	extra := []raven.Interface{}
	var tags map[string]string

	if c != nil {
		extra = append(extra, raven.NewHttp(c.Request))

		if user, ok := c.Get("user"); ok {
			if user, ok := user.(*flentities.User); ok {
				ruser := &raven.User{
					ID:    fmt.Sprint(user.ID),
					Email: user.Email,
				}

				if user.Name != nil {
					ruser.Username = *user.Name
				}

				extra = append(extra, ruser)
			}
		}

		if diagnostics, ok := c.Get("diagnostics"); ok {
			if diagnostics, ok := diagnostics.(*fldiagnostics.Diagnostics); ok {
				tags = diagnostics.ToMap()
				diagnostics.AddError(err)
			}
		}
	}

	raven.CaptureError(err, tags, extra...)
}
