Feature: Signing in

  Scenario: Signing in for the first time with Facebook
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123      |
      | id          | 123            |
      | name        | Diego Selzlein |
    Then we should have the following users:
      | email | name           | facebookId |
      |       | Diego Selzlein | 123        |
    And the system should respond with "OK"
