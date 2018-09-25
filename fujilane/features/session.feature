Feature: Session

  Scenario: Accessing a protected resource without an authentication token
    When I add a new property
    Then the system should respond with "UNAUTHORIZED"
    And we should have no properties

  Scenario: Accessing a protected resource with an authentication token for an invalid user
    Given it is currently "01 Jun 18 08:00"
    And the following session:
      | Email      | diego@selzlein.com   |
      | IssuedAt   | 2018-06-01T08:00:00Z |
      | RenewAfter | 2018-06-05T08:00:00Z |
      | ExpiresAt  | 2018-06-08T08:00:00Z |
    When I add a new property
    Then the system should respond with "UNAUTHORIZED" and the following errors:
      | You need to sign in |
    And we should have no properties

  Scenario: Accessing a protected resource with an expired session
    Given it is currently "01 Jun 18 08:00"
    And the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And the following session:
      | Email      | diego@selzlein.com   |
      | IssuedAt   | 2018-05-21T08:00:00Z |
      | RenewAfter | 2018-05-25T08:00:00Z |
      | ExpiresAt  | 2018-05-28T08:00:00Z |
    When I add a new property
    Then the system should respond with "UNAUTHORIZED" and the following errors:
      | Your session expired |
    And we should have no properties
