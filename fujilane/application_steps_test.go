package fujilane

import (
	"log"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rdumont/assistdog"
)

var application *Application
var router *gin.Engine
var assist *assistdog.Assist

const timeFormat = "02 Jan 06 15:04"

func setupApplication() {
	assist = assistdog.NewDefault()
	facebookClient = &mockedFacebookClient{tokens: map[string]facebookTokenDetails{}}
	application = NewApplication(facebookClient)
	router = application.CreateRouter()
}

func cleanup(_ interface{}, _ error) {
	err := withDatabase(func(db *gorm.DB) error {
		return db.Unscoped().Delete(User{}).Error
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
func itIsCurrently(timeExpr string) error {
	t, err := time.Parse(timeFormat, timeExpr)
	if err != nil {
		return err
	}

	application.timeFunc = func() time.Time {
		return t
	}

	return nil
}

func ApplicationContext(s *godog.Suite) {
	s.BeforeSuite(setupApplication)
	s.AfterScenario(cleanup)

	s.Step(`^it is currently "([^"]*)"$`, itIsCurrently)
}
