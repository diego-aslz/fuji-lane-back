package fujilane

import (
	"github.com/nerde/fuji-lane-back/flweb"

	"github.com/DATA-DOG/godog"
)

func requestStatus() error {
	return performGET(flweb.StatusPath)
}

func StatusContext(s *godog.Suite) {
	s.Step(`^I request a status check$`, requestStatus)
}
