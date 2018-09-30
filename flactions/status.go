package flactions

import "net/http"

// Status is useful for health checks
type Status struct{}

// Perform the action
func (a *Status) Perform(c Context) {
	c.Respond(http.StatusOK, map[string]string{"status": "active"})
}
