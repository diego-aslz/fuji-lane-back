package flutils

import "os"

const projectPath = "github.com/nerde/fuji-lane-back"

// Root returns the app's root path
func Root() string {
	appRoot := os.Getenv("APP_ROOT")
	if appRoot != "" {
		return appRoot
	}

	return os.Getenv("GOPATH") + "/src/" + projectPath
}
