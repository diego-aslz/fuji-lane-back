package fljobs

import (
	"fmt"

	"github.com/nerde/fuji-lane-back/flservices"
)

// Application handles background jobs
type Application struct {
	Adapter
	Mailer flservices.Mailer
}

// AddJobs adds all supported jobs
func (a *Application) AddJobs() {
	a.Adapter.Add(BookingCreatedJob, newBookingCreated(a.Mailer))
}

// Start the jobs application
func (a *Application) Start() {
	a.AddJobs()

	fmt.Println("Starting jobs application")
	if err := a.Adapter.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}

// NewDefaultApplication returns a new jobs application with default dependencies
func NewDefaultApplication() *Application {
	return &Application{Mailer: flservices.NewSMTPMailer(), Adapter: NewWorkersAdapter()}
}
