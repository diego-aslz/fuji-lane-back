package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flweb"
)

func SignUpTestsContext(s *godog.Suite) {
	s.Step(`^the following user signs up with his email:$`, postTableStep(flweb.SignUpPath))
}
