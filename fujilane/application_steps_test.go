package fujilane

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rdumont/assistdog"
	"github.com/rdumont/assistdog/defaults"
)

var application *Application
var router *gin.Engine
var assist *assistdog.Assist

const timeFormat = "02 Jan 06 15:04"

func setupApplication() {
	os.Setenv("STAGE", "test")

	LoadConfiguration()

	assist = assistdog.NewDefault()
	assist.RegisterComparer(time.Time{}, timeComparer)

	facebookClient = &mockedFacebookClient{tokens: map[string]facebookTokenDetails{}}
	application = NewApplication(facebookClient)

	router = application.CreateRouter()
}

func timeComparer(raw string, rawActual interface{}) error {
	at, ok := rawActual.(time.Time)
	if !ok {
		return fmt.Errorf("%v is not time.Time", rawActual)
	}

	et, err := defaults.ParseTime(raw)
	if err != nil {
		return err
	}

	expected := et.(time.Time).UTC()
	actual := at.UTC()
	if expected != actual {
		return fmt.Errorf("Expected %v, but got %v", expected, actual)
	}

	return nil
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
