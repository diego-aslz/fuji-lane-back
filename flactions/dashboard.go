package flactions

import (
	"errors"
	"net/http"
	"time"

	"github.com/nerde/fuji-lane-back/flreports"
)

// Dashboard lists user properties
type Dashboard struct {
	Context
}

// Perform executes the action
func (a *Dashboard) Perform() {
	user := a.CurrentUser()

	rawSince := a.Query("since")
	rawUntil := a.Query("until")

	a.Diagnostics().Add("since", rawSince).Add("until", rawUntil)

	if rawSince == "" || rawUntil == "" {
		a.RespondError(http.StatusBadRequest, errors.New("Parameters 'since' and 'until' are required"))
		return
	}

	since, err := time.Parse(time.RFC3339, rawSince)
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondError(http.StatusBadRequest, errors.New("Invalid 'since' parameter"))
		return
	}

	var until time.Time
	until, err = time.Parse(time.RFC3339, rawUntil)
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondError(http.StatusBadRequest, errors.New("Invalid 'until' parameter"))
		return
	}

	report, err := flreports.NewDashboard(a.Repository(), *user.AccountID, since, until)
	if err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, report)
}

// NewDashboard returns a new Dashboard action
func NewDashboard(c Context) Action {
	return &Dashboard{c}
}
