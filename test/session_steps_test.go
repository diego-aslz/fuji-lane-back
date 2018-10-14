package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

var currentSession *flentities.Session

func iAmAuthenticatedWith(email string) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		user, err := r.FindUserByEmail(email)
		if err != nil {
			return err
		}
		if user.ID == 0 {
			return fmt.Errorf("User not found: %s", email)
		}

		currentSession = flentities.NewSession(user, application.TimeFunc)
		currentSession.GenerateToken()

		return nil
	})
}

func theFollowingSession(table *gherkin.DataTable) error {
	s, err := assist.CreateInstance(new(flentities.Session), table)

	if err != nil {
		return err
	}

	currentSession = s.(*flentities.Session)
	currentSession.Secret = flconfig.Config.TokenSecret
	currentSession.GenerateToken()

	return nil
}

func resetSession(_ interface{}, _ error) {
	currentSession = nil
}

func SessionContext(s *godog.Suite) {
	s.AfterScenario(resetSession)

	s.Step(`^I am authenticated with "([^"]*)"$`, iAmAuthenticatedWith)
	s.Step(`^I sign in with:$`, postTableStep(flweb.SignInPath))
	s.Step(`^the following session:$`, theFollowingSession)
}
