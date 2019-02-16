Feature: Newsletter

  Scenario: Subscribing to the newsletter
    When I subscribe to newsletter with:
      | Email | diegoselzlein@fujilane.com |
    Then I should have the following newsletter subscriptions:
      | Email                      |
      | diegoselzlein@fujilane.com |
