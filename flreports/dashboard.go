package flreports

import (
	"time"

	"github.com/nerde/fuji-lane-back/flentities"
)

// Dashboard is used to show owners an overview of their bookings and visits
type Dashboard struct {
	Totals map[string]int64 `json:"totals"`
	Daily  []DashboardDaily `json:"daily"`
}

// DashboardDaily holds totals for one day
type DashboardDaily struct {
	Date          string `json:"date"`
	BookingsCount int64  `json:"bookingsCount"`
	VisitsCount   int64  `json:"visitsCount"`
}

type dateTotal struct {
	Date  time.Time
	Total int64
}

// NewDashboard returns a new DashboardReport with values for the given parameters
func NewDashboard(r *flentities.Repository, accountID uint, since, until time.Time) (*Dashboard, error) {
	report := &Dashboard{
		Totals: map[string]int64{
			"searches":  0,
			"visits":    0,
			"requests":  0,
			"favorites": 0,
		},
		Daily: []DashboardDaily{},
	}
	dateLayout := "2006-01-02"

	bookingTotals := []dateTotal{}
	err := r.Model(&flentities.Booking{}).Joins("JOIN units ON bookings.unit_id = units.id").
		Joins("JOIN properties ON units.property_id = properties.id").Where("properties.account_id = ?", accountID).
		Where("bookings.created_at >= ?", since).Where("bookings.created_at < ?", until).
		Select("date(bookings.created_at) date, count(*) total").Group("date(bookings.created_at)").
		Scan(&bookingTotals).Error

	if err != nil {
		return nil, err
	}

	s := since
	for s.Before(until) {
		formattedDate := s.Format(dateLayout)

		total := dateTotal{}
		for _, t := range bookingTotals {
			if t.Date.Format(dateLayout) == formattedDate {
				total = t
				break
			}
		}

		report.Daily = append(report.Daily, DashboardDaily{
			Date:          formattedDate,
			BookingsCount: total.Total,
		})

		report.Totals["requests"] += total.Total

		s = s.Add(24 * time.Hour)
	}

	return report, nil
}
