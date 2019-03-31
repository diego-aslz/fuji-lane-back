package flservices

import (
	"github.com/futurenda/google-auth-id-token-verifier"
)

// GoogleAuth validates and decodes JSON Web Tokens from Google
type GoogleAuth interface {
	Verify(string) error
	Decode(string) (map[string]string, error)
}

type googleAuth struct {
	appID string
}

func (g *googleAuth) Verify(rawToken string) error {
	v := googleAuthIDTokenVerifier.Verifier{}
	return v.VerifyIDToken(rawToken, []string{g.appID})
}

func (g *googleAuth) Decode(rawToken string) (map[string]string, error) {
	claimSet, err := googleAuthIDTokenVerifier.Decode(rawToken)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"email":    claimSet.Email,
		"name":     claimSet.Name,
		"picture":  claimSet.Picture,
		"googleID": claimSet.Sub,
	}, nil
}

// NewGoogleAuth returns a new implementation of GoogleAuth
func NewGoogleAuth(appID string) GoogleAuth {
	return &googleAuth{appID: appID}
}
