package fujilane

import (
	"os"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/flweb"
	"github.com/rdumont/assistdog"
)

var application *flweb.Application
var router *gin.Engine
var assist *assistdog.Assist
var appTime time.Time

const timeFormat = "02 Jan 06 15:04"
const fullTimeFormat = "2006-01-02T15:04:05Z"

type fakeRandSource struct{}

func (f fakeRandSource) Int63() int64 {
	return 0
}

func (f fakeRandSource) Seed(int64) {}

func setupApplication() {
	os.Setenv("STAGE", "test")

	flconfig.LoadConfiguration()

	s3, err := flservices.NewS3()
	if err != nil {
		panic(err)
	}

	application = &flweb.Application{
		RandSource: fakeRandSource{},
		S3Service:  newFakeS3(s3),
		TimeFunc: func() time.Time {
			return appTime
		},
	}

	router = application.CreateRouter()
}

func itIsCurrently(timeExpr string) (err error) {
	appTime, err = time.Parse(timeFormat, timeExpr)
	if err != nil {
		return
	}

	return
}

func forceConfiguration(table *gherkin.DataTable) error {
	config, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	for k, v := range config {
		os.Setenv(k, v)
	}

	flconfig.LoadConfiguration()

	return nil
}

func ApplicationContext(s *godog.Suite) {
	s.BeforeSuite(setupApplication)
	s.BeforeScenario(func(_ interface{}) {
		appTime = time.Now()
	})

	s.Step(`^it is currently "([^"]*)"$`, itIsCurrently)
	s.Step(`^the following configuration:$`, forceConfiguration)
}
