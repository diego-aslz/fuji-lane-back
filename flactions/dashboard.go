package flactions

import (
	"errors"
	"net/http"
	"time"

	"github.com/nerde/fuji-lane-back/flreports"
)

// Dashboard lists user properties
type Dashboard struct{}

// Perform executes the action
func (a *Dashboard) Perform(c Context) {
	user := c.CurrentUser()

	rawSince := c.Query("since")
	rawUntil := c.Query("until")

	c.Diagnostics().Add("since", rawSince).Add("until", rawUntil)

	if rawSince == "" || rawUntil == "" {
		c.RespondError(http.StatusBadRequest, errors.New("Parameters 'since' and 'until' are required"))
		return
	}

	since, err := time.Parse(time.RFC3339, rawSince)
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondError(http.StatusBadRequest, errors.New("Invalid 'since' parameter"))
		return
	}

	var until time.Time
	until, err = time.Parse(time.RFC3339, rawUntil)
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondError(http.StatusBadRequest, errors.New("Invalid 'until' parameter"))
		return
	}

	report, err := flreports.NewDashboard(c.Repository(), *user.AccountID, since, until)
	if err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, report)
}
