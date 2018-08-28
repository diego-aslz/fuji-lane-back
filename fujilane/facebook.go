package fujilane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type facebookTokenDetails struct {
	AppID       string
	Application string
	ExpiresAt   int
	IsValid     bool
	Scopes      []string
	UserID      string
}

type facebookTokenDetailsResponse struct {
	FacebookTokenDetails facebookTokenDetails `json:"data"`
}

// FacebookClient handles debugging authentication tokens
type FacebookClient interface {
	debug(string) (facebookTokenDetails, error)
}

type facebookHTTPClient struct {
	appToken string
}

func (c *facebookHTTPClient) debug(token string) (facebookTokenDetails, error) {
	url := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s", token, c.appToken)
	resp, err := http.Get(url)
	if err != nil {
		return facebookTokenDetails{}, err
	}
	defer resp.Body.Close()

	facebookTokenDetailsResponse := &facebookTokenDetailsResponse{}

	err = json.NewDecoder(resp.Body).Decode(facebookTokenDetailsResponse)
	if err != nil {
		return facebookTokenDetails{}, err
	}

	return facebookTokenDetailsResponse.FacebookTokenDetails, nil
}

type facebook struct {
	appID  string
	client FacebookClient
}

func (f *facebook) validate(token string, userID string) error {
	details, err := f.client.debug(token)
	if err != nil {
		return err
	}

	if !details.IsValid {
		return fmt.Errorf("Token %s is not valid according to Facebook", token)
	}

	if details.AppID != f.appID {
		return fmt.Errorf("Token %s belongs to app %s, expected %s", token, details.AppID, f.appID)
	}

	if details.UserID != userID {
		return fmt.Errorf("Token %s belongs to user %s, expected %s", token, details.UserID, userID)
	}

	return nil
}

func newFacebook(client FacebookClient) *facebook {
	f := &facebook{
		client: client,
		appID:  os.Getenv("FACEBOOK_APP_ID"),
	}

	return f
}

// NewFacebookHTTPClient creates a facebookClient that makes HTTP requests to validate tokens against Facebook API
func NewFacebookHTTPClient() FacebookClient {
	return &facebookHTTPClient{os.Getenv("FACEBOOK_CLIENT_TOKEN")}
}
