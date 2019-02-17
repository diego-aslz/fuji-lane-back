package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flweb"
)

func DashboardContext(s *godog.Suite) {
	s.Step(`^I get dashboard details for:$`, performGETWithParamsStep(flweb.DashboardPath))
	s.Step(`^I list my properties\' bookings$`, performGETStep(flweb.DashboardBookingsPath))
	s.Step(`^I list my properties\' bookings for page "([^"]*)"$`, performGETStepWithPage(flweb.DashboardBookingsPath))
}
