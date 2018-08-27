package fujilane

import (
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
)

func iRequestAStatusCheck() error {
	response = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/status", nil)

	if err == nil {
		router.ServeHTTP(response, req)
	}

	return err
}

func StatusContext(s *godog.Suite) {
	s.Step(`^I request a status check$`, iRequestAStatusCheck)
}
