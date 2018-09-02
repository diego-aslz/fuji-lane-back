package fujilane

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type configuration struct {
	databaseURL         string
	facebookAppID       string
	facebookClientToken string
	stage               string
	tokenSecret         string
}

var appConfig *configuration

const projectPath = "github.com/nerde/fuji-lane-back"

// LoadConfiguration from environment
func LoadConfiguration() {
	stage := getStage()

	if err := loadEnv(stage); err != nil {
		panic(err)
	}

	appConfig = &configuration{
		databaseURL:         os.Getenv("DATABASE_URL"),
		facebookAppID:       os.Getenv("FACEBOOK_APP_ID"),
		facebookClientToken: os.Getenv("FACEBOOK_CLIENT_TOKEN"),
		stage:               stage,
		tokenSecret:         os.Getenv("TOKEN_SECRET"),
	}
}

func loadEnv(stage string) error {
	configFile := os.Getenv("GOPATH") + "/src/" + projectPath + "/config/" + stage + ".yml"

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
