package flconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Configuration contains global system configuration details
type Configuration struct {
	AWSRegion           string
	AWSBucket           string
	DatabaseLogs        bool
	DatabaseURL         string
	FacebookAppID       string
	FacebookClientToken string
	MaxImageSizeMB      int
	Stage               string
	TokenSecret         string
}

// Config is the current global configuration being used
var Config *Configuration

const projectPath = "github.com/nerde/fuji-lane-back"

// LoadConfiguration from environment
func LoadConfiguration() {
	var err error
	stage := getStage()

	if err = LoadEnv(stage); err != nil {
		panic(err)
	}

	maxImageSize := 0
	maxImageSizeVar := os.Getenv("MAX_IMAGE_SIZE_MB")
	if maxImageSizeVar != "" {
		maxImageSize, err = strconv.Atoi(maxImageSizeVar)
		if err != nil {
			fmt.Printf("Unable to load MAX_IMAGE_SIZE_MB (%s), falling back to default.\n", err.Error())
		}
	}
	if maxImageSize == 0 {
		maxImageSize = 20
	}

	Config = &Configuration{
		AWSRegion:           os.Getenv("AWS_REGION"),
		AWSBucket:           os.Getenv("AWS_BUCKET"),
		DatabaseLogs:        os.Getenv("DATABASE_LOGS") == "true",
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		FacebookAppID:       os.Getenv("FACEBOOK_APP_ID"),
		FacebookClientToken: os.Getenv("FACEBOOK_CLIENT_TOKEN"),
		MaxImageSizeMB:      maxImageSize,
		Stage:               stage,
		TokenSecret:         os.Getenv("TOKEN_SECRET"),
	}
}

// LoadEnv loads environment variables from the YAML configuration file for the current stage. If not present, it
// does nothing
func LoadEnv(stage string) error {
	configFile := os.Getenv("GOPATH") + "/src/" + projectPath + "/flconfig/" + stage + ".yml"

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
