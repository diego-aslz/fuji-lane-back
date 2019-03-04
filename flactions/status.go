package flactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status is useful for health checks
type Status struct {
	Context
}

// Perform the action
func (a *Status) Perform() {
	a.Respond(http.StatusOK, gin.H{"status": "active"})
}

// NewStatus returns a new Status action
func NewStatus(c Context) Action {
	return &Status{c}
}
