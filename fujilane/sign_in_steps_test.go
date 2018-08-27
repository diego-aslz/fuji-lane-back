package fujilane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var facebookClient *mockedFacebookClient

type mockedFacebookClient struct {
	tokens map[string]facebookTokenDetails
}

func (c *mockedFacebookClient) mock(token string, details facebookTokenDetails) {
	c.tokens[token] = details
}

func (c *mockedFacebookClient) debug(token string) (facebookTokenDetails, error) {
	details, ok := c.tokens[token]

	if !ok {
		return facebookTokenDetails{}, fmt.Errorf("Could not find details for token %s", token)
	}

	return details, nil
}

func facebookRecognizesTheFollowingTokens(table *gherkin.DataTable) error {
	detailRows, err := assist.ParseSlice(table)
	if err != nil {
		return err
	}

	for _, row := range detailRows {
		facebookClient.mock(row["accessToken"], facebookTokenDetails{
			AppID:   row["AppID"],
			IsValid: row["IsValid"] == "true",
			UserID:  row["UserID"],
		})
	}

	return nil
}

func theFollowingUserSignsInViaFacebook(table *gherkin.DataTable) error {
	response = httptest.NewRecorder()

	credentials, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	json, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "/sign_in/facebook", strings.NewReader(string(json)))

	if err != nil {
		return err
	}

	router.ServeHTTP(response, req)

	return nil
}

func weShouldHaveTheFollowingUsers(table *gherkin.DataTable) error {
	return godog.ErrPending
}

func SignInContext(s *godog.Suite) {
	s.Step(`^Facebook recognizes the following tokens:$`, facebookRecognizesTheFollowingTokens)
	s.Step(`^the following user signs in via Facebook:$`, theFollowingUserSignsInViaFacebook)
	s.Step(`^we should have the following users:$`, weShouldHaveTheFollowingUsers)
}
