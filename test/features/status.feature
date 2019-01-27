Feature: Status check

  Scenario: Checking system status
    When I request a status check
    Then I should receive an "OK" response with the following body:
      | status | active |
