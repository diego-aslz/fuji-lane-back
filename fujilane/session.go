package fujilane

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type session struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	RenewAfter time.Time `json:"renew_after"`
	secret     string
}

func newSession(user *User, timeFunc func() time.Time) *session {
	now := timeFunc()
	return &session{
		Email:      user.Email,
		IssuedAt:   now,
		RenewAfter: now.Add(4 * 24 * time.Hour),
		ExpiresAt:  now.Add(7 * 24 * time.Hour),
		secret:     appConfig.tokenSecret,
	}
}

func (s *session) generateToken() (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, s.claims())
	s.Token, err = token.SignedString([]byte(s.secret))
	return
}

func (s *session) claims() jwt.MapClaims {
	return jwt.MapClaims{
		"Email":      s.Email,
		"IssuedAt":   s.IssuedAt.Unix(),
		"ExpiresAt":  s.ExpiresAt.Unix(),
		"RenewAfter": s.RenewAfter.Unix(),
	}
}
