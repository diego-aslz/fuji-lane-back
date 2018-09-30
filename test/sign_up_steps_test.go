package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flweb"
)

func emailSignUp(table *gherkin.DataTable) error {
	return performPOSTWithTable(flweb.SignUpPath, table)
}

func SignUpTestsContext(s *godog.Suite) {
	s.Step(`^the following user signs up with his email:$`, emailSignUp)
}
