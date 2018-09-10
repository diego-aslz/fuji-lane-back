package fujilane

import (
	"errors"
	"fmt"
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

func loadSession(tokenStr string) (*session, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(appConfig.tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token is invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &session{
			Email:      claims["Email"].(string),
			IssuedAt:   time.Unix(int64(claims["IssuedAt"].(float64)), 0),
			RenewAfter: time.Unix(int64(claims["RenewAfter"].(float64)), 0),
			ExpiresAt:  time.Unix(int64(claims["ExpiresAt"].(float64)), 0),
			secret:     appConfig.tokenSecret,
		}, nil
	}

	return nil, errors.New("Could not cast claims")
}
