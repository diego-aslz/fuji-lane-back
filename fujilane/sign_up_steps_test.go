package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func emailSignUp(table *gherkin.DataTable) error {
	return makePOSTRequest(signUpPath, table)
}

func SignUpTestsContext(s *godog.Suite) {
	s.Step(`^the following user signs up with his email:$`, emailSignUp)
}
