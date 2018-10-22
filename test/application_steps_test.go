package fujilane

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/flweb"
	"github.com/rdumont/assistdog"
	"github.com/rdumont/assistdog/defaults"
)

var application *flweb.Application
var router *gin.Engine
var assist *assistdog.Assist
var appTime time.Time

const timeFormat = "02 Jan 06 15:04"

type fakeRandSource struct{}

func (f fakeRandSource) Int63() int64 {
	return 0
}

func (f fakeRandSource) Seed(int64) {
}

func setupApplication() {
	os.Setenv("STAGE", "test")

	flconfig.LoadConfiguration()

	assist = assistdog.NewDefault()
	assist.RegisterComparer(time.Time{}, timeComparer)
	assist.RegisterComparer(true, boolComparer)
	assist.RegisterComparer(uint(0), uintComparer)
	assist.RegisterComparer(refStr("a"), strPtrComparer)
	assist.RegisterComparer(refInt(1), intPtrComparer)
	assist.RegisterComparer(refUint(1), uintPtrComparer)
	assist.RegisterParser(uint(0), uintParser)
	assist.RegisterParser(true, boolParser)
	assist.RegisterParser(refStr("a"), strPtrParser)
	assist.RegisterParser(refInt(1), intPtrParser)

	facebookClient = &mockedFacebookClient{tokens: map[string]flservices.FacebookTokenDetails{}}
	application = flweb.NewApplication(facebookClient)
	application.RandSource = fakeRandSource{}
	application.S3Service = newFakeS3(application.S3Service)

	application.TimeFunc = func() time.Time {
		return appTime
	}

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

func uintComparer(raw string, rawActual interface{}) error {
	rawInt, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}

	if uint(rawInt) == rawActual.(uint) {
		return nil
	}

	return fmt.Errorf("Expected %d, but got %d", rawInt, rawActual)
}

func boolComparer(raw string, rawActual interface{}) error {
	actual := fmt.Sprint(rawActual)
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func strPtrComparer(raw string, rawActual interface{}) error {
	actual := derefStr(rawActual.(*string))
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func intPtrComparer(raw string, rawActual interface{}) error {
	actual := strconv.Itoa(derefInt(rawActual.(*int)))
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func uintPtrComparer(raw string, rawActual interface{}) error {
	actual := fmt.Sprint(derefUint(rawActual.(*uint)))
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func uintParser(raw string) (interface{}, error) {
	i, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return uint(i), nil
}

func boolParser(raw string) (interface{}, error) {
	if raw != "true" && raw != "false" {
		return nil, fmt.Errorf("Don't know how to parse \"%s\" to bool", raw)
	}

	return raw == "true", nil
}

func strPtrParser(raw string) (interface{}, error) {
	if raw == "" {
		return nil, nil
	}

	return &raw, nil
}

func intPtrParser(raw string) (interface{}, error) {
	if raw == "" {
		return nil, nil
	}

	i, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return &i, nil
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

	if value, ok := config["MAX_IMAGE_SIZE_MB"]; ok {
		if flconfig.Config.MaxImageSizeMB, err = strconv.Atoi(value); err != nil {
			return err
		}
	}

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
