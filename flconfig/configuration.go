package flconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	raven "github.com/getsentry/raven-go"
	"github.com/nerde/fuji-lane-back/fujilane"
	yaml "gopkg.in/yaml.v2"
)

// SMTP represents a SMTP configuration
type SMTP struct {
	Auth        string
	DefaultFrom string
	Host        string
	Port        int
	User        string
	Password    string
}

// Configuration contains global system configuration details
type Configuration struct {
	AppURL                   string
	AWSRegion                string
	AWSBucket                string
	DatabaseLogs             bool
	DatabaseURL              string
	FacebookAppID            string
	FacebookClientToken      string
	GoogleAuthAppID          string
	JobsStatsPort            int
	MaxImageSizeMB           int
	RedisURL                 string
	Stage                    string
	TokenSecret              string
	SendgridKey              string
	SendgridNewsletterListID string
	SMTP                     SMTP
}

// Config is the current global configuration being used
var Config *Configuration

// LoadConfiguration from environment
func LoadConfiguration() {
	var err error
	stage := getStage()

	if err = LoadEnv(stage); err != nil {
		panic(err)
	}

	maxImageSize := getIntVar("MAX_IMAGE_SIZE_MB")
	if maxImageSize == 0 {
		maxImageSize = 20
	}

	Config = &Configuration{
		AppURL:                   os.Getenv("APP_URL"),
		AWSRegion:                os.Getenv("AWS_REGION"),
		AWSBucket:                os.Getenv("AWS_BUCKET"),
		DatabaseLogs:             os.Getenv("DATABASE_LOGS") == "true",
		DatabaseURL:              os.Getenv("DATABASE_URL"),
		FacebookAppID:            os.Getenv("FACEBOOK_APP_ID"),
		FacebookClientToken:      os.Getenv("FACEBOOK_CLIENT_TOKEN"),
		GoogleAuthAppID:          os.Getenv("GOOGLE_AUTH_APP_ID"),
		JobsStatsPort:            getIntVar("JOBS_STATS_PORT"),
		MaxImageSizeMB:           maxImageSize,
		RedisURL:                 os.Getenv("REDIS_URL"),
		Stage:                    stage,
		TokenSecret:              os.Getenv("TOKEN_SECRET"),
		SendgridKey:              os.Getenv("SENDGRID_KEY"),
		SendgridNewsletterListID: os.Getenv("SENDGRID_NEWSLETTER_LIST_ID"),
		SMTP: SMTP{
			Auth:        os.Getenv("SMTP_AUTH"),
			DefaultFrom: os.Getenv("SMTP_FROM"),
			Host:        os.Getenv("SMTP_HOST"),
			Port:        getIntVar("SMTP_PORT"),
			User:        os.Getenv("SMTP_USER"),
			Password:    os.Getenv("SMTP_PASSWORD"),
		},
	}
}

func getIntVar(name string) int {
	value := os.Getenv(name)

	if value == "" {
		return 0
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Unable to parse int var %s with value %s (%s)\n", name, value, err.Error())
	}

	return i
}

// LoadEnv loads environment variables from the YAML configuration file for the current stage. If not present, it
// does nothing
func LoadEnv(stage string) error {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.SetRelease(os.Getenv("HEROKU_RELEASE_VERSION"))
	raven.SetEnvironment(os.Getenv("SENTRY_ENVIRONMENT"))

	configFile := fujilane.Root() + "/flconfig/" + stage + ".yml"

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No configuration file found at %s, skipping\n", configFile)
			return nil
		}

		return err
	}

	config := map[string]string{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}

	for varName, value := range config {
		os.Setenv(varName, value)
	}

	return nil
}

func getStage() string {
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "development"
	}
	return stage
}
