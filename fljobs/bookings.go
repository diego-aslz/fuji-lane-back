package fljobs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/raven-go"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flemail"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// BookingCreatedJob is the class name to identify this job
const BookingCreatedJob = "BookingCreated"

// EnqueueBookingCreated enqueues a BookingCreated job in the background
func (a *Application) EnqueueBookingCreated(id uint) (string, error) {
	return a.Adapter.Enqueue(BookingCreatedJob, id)
}

// newBookingCreated return a new BookingCreated job which notifies the owner about this booking
func newBookingCreated(mailer flservices.Mailer) JobFunc {
	return func(c *Context) (err error) {
		c.Add("job", BookingCreatedJob)
		c.Add("at", time.Now().Format("2006-01-02T15:04:05Z"))

		defer func() {
			if err != nil {
				c.AddError(err)
			} else if rec := recover(); rec != nil {
				if e, ok := rec.(error); ok {
					c.AddError(e)
					err = e
				}
			}

			if err != nil {
				raven.CaptureError(err, c.ToMap())
			}

			log.Println(c)
		}()

		id, _ := c.Args[0].(json.Number).Int64()
		c.Add("booking_id", fmt.Sprint(id))

		err = flentities.WithRepository(func(r *flentities.Repository) error {
			booking := &flentities.Booking{}
			if err := r.Preload("Unit.Property").Preload("User").Find(booking, id).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					c.AddQuoted("message", "Booking is not found, ignoring")
					return nil
				}

				return err
			}
			c.AddJSON("booking", booking)

			owner := &flentities.User{}
			if err := r.Where("account_id = ?", booking.Unit.Property.AccountID).Find(owner).Error; err != nil {
				return err
			}

			c.AddJSON("owner", owner)

			mail, err := flemail.BookingCreated(booking, owner)
			if err != nil {
				return err
			}

			if err := mailer.Send(mail); err != nil {
				return err
			}

			c.AddQuoted("message", "Email sent successfully")

			return nil
		})

		return
	}
}
