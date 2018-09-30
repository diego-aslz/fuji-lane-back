package flactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nerde/fuji-lane-back/flconfig"
)

// FacebookTokenDetails contains information for a given token
type FacebookTokenDetails struct {
	AppID       string   `json:"app_id"`
	Application string   `json:"application"`
	ExpiresAt   int      `json:"expires_at"`
	IsValid     bool     `json:"is_valid"`
	Scopes      []string `json:"scopes"`
	UserID      string   `json:"user_id"`
}

// FacebookClient handles debugging authentication tokens
type FacebookClient interface {
	Debug(string) (FacebookTokenDetails, error)
}

type facebook struct {
	appID  string
	client FacebookClient
}

func (f *facebook) validate(token string, userID string) error {
	details, err := f.client.Debug(token)
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
	return &facebook{
		client: client,
		appID:  flconfig.Config.FacebookAppID,
	}
}

type facebookTokenDetailsResponse struct {
	FacebookTokenDetails `json:"data"`
}

// FacebookHTTPClient is the HTTP implementation of FacebookClient.
type FacebookHTTPClient struct {
	appToken string
}

// Debug makes HTTP calls to Facebook to verify if the given token is invalid
func (c *FacebookHTTPClient) Debug(token string) (FacebookTokenDetails, error) {
	url := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s", token, c.appToken)
	resp, err := http.Get(url)
	if err != nil {
		return FacebookTokenDetails{}, err
	}
	defer resp.Body.Close()

	facebookTokenDetailsResponse := &facebookTokenDetailsResponse{}

	err = json.NewDecoder(resp.Body).Decode(facebookTokenDetailsResponse)
	if err != nil {
		return FacebookTokenDetails{}, err
	}

	return facebookTokenDetailsResponse.FacebookTokenDetails, nil
}

// NewFacebookHTTPClient creates a facebookClient that makes HTTP requests to validate tokens against Facebook API
func NewFacebookHTTPClient() FacebookClient {
	return &FacebookHTTPClient{flconfig.Config.FacebookClientToken}
}
