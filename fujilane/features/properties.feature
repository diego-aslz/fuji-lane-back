Feature: Properties Management

  Scenario: Adding a new property
    Given the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    When I add a new property
    Then the system should respond with "CREATED"
    And I should have the following properties:
      | User                 | State   |
      | Diego Aguir Selzlein | Pending |
