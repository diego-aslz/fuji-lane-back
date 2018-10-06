package flentities

import (
	"errors"
	"fmt"
	"time"

	"github.com/nerde/fuji-lane-back/fldiagnostics"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nerde/fuji-lane-back/flconfig"
)

// Session contains authentication information
type Session struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	RenewAfter time.Time `json:"renew_after"`
	User       *User     `json:"user"`
	Account    *Account  `json:"account"`
	Secret     string    `json:"-"`
}

// GenerateToken generates a signed token for this session
func (s *Session) GenerateToken() (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, s.claims())
	s.Token, err = token.SignedString([]byte(s.Secret))
	return
}

func (s *Session) claims() jwt.MapClaims {
	return jwt.MapClaims{
		"Email":      s.Email,
		"IssuedAt":   s.IssuedAt.Unix(),
		"ExpiresAt":  s.ExpiresAt.Unix(),
		"RenewAfter": s.RenewAfter.Unix(),
	}
}

// FilterSensitiveInformation hides the Token
func (s Session) FilterSensitiveInformation() fldiagnostics.SensitivePayload {
	s.Token = "[FILTERED]"
	return s
}

// NewSession returns a new Session for the given user and calculates expiration time with the given time function
func NewSession(user *User, timeFunc func() time.Time) *Session {
	now := timeFunc()
	return &Session{
		Email:      user.Email,
		IssuedAt:   now,
		RenewAfter: now.Add(4 * 24 * time.Hour),
		ExpiresAt:  now.Add(7 * 24 * time.Hour),
		User:       user,
		Account:    user.Account,
		Secret:     flconfig.Config.TokenSecret,
	}
}

// LoadSession returns a Session loaded from the given token
func LoadSession(tokenStr string) (*Session, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(flconfig.Config.TokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token is invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &Session{
			Email:      claims["Email"].(string),
			IssuedAt:   time.Unix(int64(claims["IssuedAt"].(float64)), 0),
			RenewAfter: time.Unix(int64(claims["RenewAfter"].(float64)), 0),
			ExpiresAt:  time.Unix(int64(claims["ExpiresAt"].(float64)), 0),
			Secret:     flconfig.Config.TokenSecret,
		}, nil
	}

	return nil, errors.New("Could not cast claims")
}
