package fujilane

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
)

var response *httptest.ResponseRecorder

func theSystemShouldRespondWith(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	switch status {
	case "OK":
		return assertResponseStatus(http.StatusOK)
	case "UNAUTHORIZED":
		return assertResponseStatus(http.StatusUnauthorized)
	default:
		return fmt.Errorf("Unhandled status: %s", status)
	}
}

func assertResponseStatus(expected int) error {
	if response.Code != expected {
		return fmt.Errorf("Expected response to be status %d, got %d", expected, response.Code)
	}
	return nil
}

func HTTPContext(s *godog.Suite) {
	s.Step(`^the system should respond with "([^"]*)"$`, theSystemShouldRespondWith)
}
