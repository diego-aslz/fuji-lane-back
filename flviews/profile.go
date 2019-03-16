package flviews

import (
	"github.com/nerde/fuji-lane-back/flentities"
)

type profile struct {
	Name                       *string `json:"name"`
	Email                      string  `json:"email"`
	AccountUnreadBookingsCount *int    `json:"accountUnreadBookingsCount"`
	AccountBookingsCount       *int    `json:"accountBookingsCount"`
	Picture                    *string `json:"picture"`
}

// NewProfile returns a structure representing an user's profile
func NewProfile(u *flentities.User) interface{} {
	pr := profile{
		Name:    u.Name,
		Email:   u.Email,
		Picture: u.Picture(),
	}

	if u.Account != nil {
		pr.AccountUnreadBookingsCount = &u.UnreadBookingsCount
		pr.AccountBookingsCount = &u.Account.BookingsCount
	}

	return pr
}
