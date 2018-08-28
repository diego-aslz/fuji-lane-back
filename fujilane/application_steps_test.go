package fujilane

import (
	"log"

	"github.com/DATA-DOG/godog"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rdumont/assistdog"
)

var router *gin.Engine
var assist *assistdog.Assist

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
}
