package flviews

import (
	"time"

	"github.com/nerde/fuji-lane-back/flentities"
)

type session struct {
	Token      string          `json:"token"`
	IssuedAt   time.Time       `json:"issuedAt"`
	ExpiresAt  time.Time       `json:"expiresAt"`
	RenewAfter time.Time       `json:"renewAfter"`
	User       interface{}     `json:"user"`
	Account    *sessionAccount `json:"account"`
}

type sessionAccount struct {
	Name          string  `json:"name"`
	Phone         *string `json:"phone"`
	CountryID     *uint   `json:"countryID"`
	BookingsCount int     `json:"bookingsCount"`
}

// NewSession returns a representation of session details
func NewSession(s *flentities.Session) interface{} {
	se := session{
		Token:      s.Token,
		IssuedAt:   s.IssuedAt,
		ExpiresAt:  s.ExpiresAt,
		RenewAfter: s.RenewAfter,
		User:       NewProfile(s.User),
	}

	if s.Account != nil {
		se.Account = &sessionAccount{
			Name:          s.Account.Name,
			Phone:         s.Account.Phone,
			CountryID:     s.Account.CountryID,
			BookingsCount: s.Account.BookingsCount,
		}
	}

	return se
}
