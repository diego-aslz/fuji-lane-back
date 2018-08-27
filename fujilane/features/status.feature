Feature: Status check

  Scenario: Checking system status
    When I request a status check
    Then the system should respond with "OK"
