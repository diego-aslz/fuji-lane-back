package flactions

import "net/http"

// Status is useful for health checks
type Status struct {
	Context
}

// Perform the action
func (a *Status) Perform() {
	a.Respond(http.StatusOK, map[string]string{"status": "active"})
}

// NewStatus returns a new Status action
func NewStatus(c Context) Action {
	return &Status{c}
}
