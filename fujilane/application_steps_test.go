package fujilane

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		if response.Code != http.StatusOK {
			return fmt.Errorf("Expected response to be status %d, got %d", http.StatusOK, response.Code)
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

func cleanup(_ interface{}, _ error) {
	err := withDatabase(func(db *gorm.DB) error {
		return db.Unscoped().Delete(User{}).Error
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func ApplicationContext(s *godog.Suite) {
	s.BeforeSuite(setupApplication)
	s.AfterScenario(cleanup)
	s.Step(`^the system should respond with "([^"]*)"$`, theSystemShouldRespondWith)
}
