Feature: Facebook Signing in

  Scenario: Signing in for the first time
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email              | Name           | FacebookID |
      | diego@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK"

  Scenario: Signing in for a second time updating attributes
    Given the following users:
      | Email              | Name                 | FacebookID |
      | diego@selzlein.com | Diego Aguir Selzlein | 123        |
    And Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-222   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-222          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email              | Name           | FacebookID |
      | diego@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK"

  Scenario: Signing in with unrecognized token
    When the following user signs in via Facebook:
      | accessToken | unrecognized-token |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have no users
    And the system should respond with "UNAUTHORIZED"

  Scenario: Signing in with invalid token
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | false   | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have no users
    And the system should respond with "UNAUTHORIZED"
