package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/fujilane"
)

var router *gin.Engine
var response *httptest.ResponseRecorder

func iRequestAStatusCheck() error {
	response = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/status", nil)

	if err == nil {
		router.ServeHTTP(response, req)
	}

	return err
}

func theSystemShouldRespondMeWith(status string) error {
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

func FeatureContext(s *godog.Suite) {
	s.BeforeSuite(func() {
		router = gin.Default()
		fujilane.AddRoutes(router)
	})
	s.Step(`^I request a status check$`, iRequestAStatusCheck)
	s.Step(`^the system should respond me with "([^"]*)"$`, theSystemShouldRespondMeWith)
}
