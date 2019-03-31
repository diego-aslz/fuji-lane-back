package fujilane

import (
	"net/http"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flweb"
)

type fakeGoogleAuth struct {
	claims map[string]string
}

func (g *fakeGoogleAuth) Verify(rawToken string) error {
	return nil
}

func (g *fakeGoogleAuth) Decode(_ string) (map[string]string, error) {
	return g.claims, nil
}

func performGoogleSignIn(table *gherkin.DataTable) error {
	var body map[string]string

	if len(table.Rows) > 0 {
		var err error
		body, err = assist.ParseMap(table)
		if err != nil {
			return err
		}
	}

	application.GoogleAuth = &fakeGoogleAuth{body}

	return performPOST(flweb.GoogleSignInPath, nil, func(req *http.Request) {
		req.Header["Authorization"] = []string{"123"}
	})
}

func GoogleContext(s *godog.Suite) {
	s.Step(`^the following user signs in via Google:$`, performGoogleSignIn)
}
