package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flweb"
)

var facebookClient *mockedFacebookClient

type mockedFacebookClient struct {
	tokens map[string]flactions.FacebookTokenDetails
}

func (c *mockedFacebookClient) mock(token string, details flactions.FacebookTokenDetails) {
	c.tokens[token] = details
}

func (c *mockedFacebookClient) Debug(token string) (flactions.FacebookTokenDetails, error) {
	details, ok := c.tokens[token]

	if !ok {
		return flactions.FacebookTokenDetails{}, fmt.Errorf("Could not find details for token %s", token)
	}

	return details, nil
}

func facebookRecognizesTheFollowingTokens(table *gherkin.DataTable) error {
	detailRows, err := assist.ParseSlice(table)
	if err != nil {
		return err
	}

	for _, row := range detailRows {
		facebookClient.mock(row["accessToken"], flactions.FacebookTokenDetails{
			AppID:   row["AppID"],
			IsValid: row["IsValid"] == "true",
			UserID:  row["UserID"],
		})
	}

	return nil
}

func theFollowingUserSignsInViaFacebook(table *gherkin.DataTable) error {
	return performPOSTWithTable(flweb.FacebookSignInPath, table)
}

func SignInContext(s *godog.Suite) {
	s.Step(`^Facebook recognizes the following tokens:$`, facebookRecognizesTheFollowingTokens)
	s.Step(`^the following user signs in via Facebook:$`, theFollowingUserSignsInViaFacebook)
}
