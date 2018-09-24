package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var currentSession *session

func iAmAuthenticatedWith(email string) error {
	user, err := application.usersRepository.findByEmail(email)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return fmt.Errorf("User not found: %s", email)
	}

	currentSession = newSession(user, application.timeFunc)
	currentSession.generateToken()

	return nil
}

func iSignInWith(table *gherkin.DataTable) error {
	return makePOSTRequest(signInPath, table)
}

func resetSession(_ interface{}, _ error) {
	currentSession = nil
}

func SessionContext(s *godog.Suite) {
	s.AfterScenario(resetSession)

	s.Step(`^I am authenticated with "([^"]*)"$`, iAmAuthenticatedWith)
	s.Step(`^I sign in with:$`, iSignInWith)
}
