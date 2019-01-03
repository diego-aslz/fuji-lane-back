Feature: Accounts Management

  Scenario: Creating an Owner's Account
    Given the following users:
      | Email              |
      | diego@selzlein.com |
    And the following countries:
      | Name  |
      | Japan |
    And I am authenticated with "diego@selzlein.com"
    When I create the following account:
      | UserName | Diego Selzlein    |
      | Name     | Diego Apartments  |
      | Phone    | +55 44 99999-9999 |
      | Country  | Japan             |
    Then the system should respond with "CREATED" and the following body:
      | name  | Diego Apartments  |
      | phone | +55 44 99999-9999 |
    And I should have the following accounts:
      | Name             | Phone             | Country |
      | Diego Apartments | +55 44 99999-9999 | Japan   |
    And I should have the following users:
      | Account          | Email              |
      | Diego Apartments | diego@selzlein.com |
