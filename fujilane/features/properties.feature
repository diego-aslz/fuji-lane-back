Feature: Properties Management

  Scenario: Adding a new property
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    When I add a new property
    Then the system should respond with "CREATED"
    And we should have the following properties:
      | Account          | State |
      | Diego Apartments | Draft |

  Scenario: Adding a new property without having an Account
    Given the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    When I add a new property
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | You need a company account |
    And we should have no properties
