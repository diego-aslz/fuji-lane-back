package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
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

	since, err := flentities.ParseDate(rawSince)
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondError(http.StatusBadRequest, errors.New("Invalid 'since' parameter"))
		return
	}

	var until flentities.Date
	until, err = flentities.ParseDate(rawUntil)
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
