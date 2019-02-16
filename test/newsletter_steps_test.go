package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/flweb"
)

func requestNewsletterSubscribe(table *gherkin.DataTable) error {
	body, err := assist.CreateInstance(new(flservices.SendgridContact), table)
	if err != nil {
		return err
	}

	return performPOST(flweb.NewsletterSubscribePath, body)
}

func assertNewsletterSubscriptions(table *gherkin.DataTable) error {
	expectedCount := len(table.Rows) - 1
	count := len(sendgrid.subscriptions)

	if expectedCount != count {
		return fmt.Errorf("Expected %d subscriptions, got %d", expectedCount, count)
	}

	return assist.CompareToSlice(sendgrid.subscriptions, table)
}

func NewsletterContext(s *godog.Suite) {
	s.Step(`^I subscribe to newsletter with:$`, requestNewsletterSubscribe)
	s.Step(`^I should have the following newsletter subscriptions:$`, assertNewsletterSubscriptions)
}
