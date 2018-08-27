package fujilane

import (
	"fmt"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/rdumont/assistdog"
)

var router *gin.Engine
var response *httptest.ResponseRecorder
var assist *assistdog.Assist

func theSystemShouldRespondWith(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	switch status {
	case "OK":
		if response.Code != 200 {
			return fmt.Errorf("Expected response to be status %d, got %d", 200, response.Code)
		}
	default:
		return fmt.Errorf("Unhandled status: %s", status)
	}

	return nil
}

func setupApplication() {
	assist = assistdog.NewDefault()
	facebookClient = &mockedFacebookClient{tokens: map[string]facebookTokenDetails{}}
	router = NewApplication(facebookClient).CreateRouter()
}

func HTTPContext(s *godog.Suite) {
	s.BeforeSuite(setupApplication)
	s.Step(`^the system should respond with "([^"]*)"$`, theSystemShouldRespondWith)
}
