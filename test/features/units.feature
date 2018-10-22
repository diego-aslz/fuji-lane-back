Feature: Units Management

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Adding a unit to my property
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    When I add the following unit:
      | Property     | ACME Downtown |
      | Name         | Standard Apt  |
      | Bedrooms     | 1             |
      | SizeM2       | 52            |
      | MaxOccupancy | 3             |
      | Count        | 15            |
    Then the system should respond with "CREATED"
    And I should have the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |

  Scenario: Adding an invalid unit to my property
    When I add the following unit:
      | MaxOccupancy | 3 |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | property is required            |
      | name is required                |
      | bedrooms is required            |
      | size is required                |
      | number of unit type is required |
    And I should have no units

  Scenario: Adding a unit to a property that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | State | Name          |
      | 1  | Other   | Draft | ACME Downtown |
    When I add the following unit:
      | Property     | ACME Downtown |
      | Name         | Standard Apt  |
      | Bedrooms     | 1             |
      | SizeM2       | 52            |
      | MaxOccupancy | 3             |
      | Count        | 15            |
    Then the system should respond with "NOT FOUND"
    And I should have no units
