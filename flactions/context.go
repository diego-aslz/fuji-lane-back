package flactions

import (
	"time"

	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

// Context provides the methods for actions to get information on the request and provide a response
type Context interface {
	Diagnostics() *fldiagnostics.Diagnostics
	Now() time.Time

	CurrentUser() *flentities.User

	Repository() *flentities.Repository

	// Respond sets the HTTP response status and body with the given parameters
	Respond(status int, body interface{})

	// RespondError sets the HTTP response status with the given parameter and generates an error response. The client
	// will be sent the error message
	RespondError(status int, err error)

	// ServerError handles internal errors that should not be exposed to the client. It sends status 500 and a generic
	// error message to the client
	ServerError(err error)
}
