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

  Scenario: Updating a unit
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Property      | Name          | Uploaded |
      | 3  | ACME Downtown | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name                   | Std Apartment |
      | Bedrooms               | 2             |
      | SizeM2                 | 50            |
      | MaxOccupancy           | 2             |
      | Count                  | 20            |
      | BasePriceCents         | 12000         |
      | OneNightPriceCents     | 11000         |
      | OneWeekPriceCents      | 40000         |
      | ThreeMonthsPriceCents  | 350000        |
      | SixMonthsPriceCents    | 650000        |
      | TwelveMonthsPriceCents | 1200000       |
      | FloorPlanImageID       | 3             |
    Then the system should respond with "OK"
    And I should have the following units:
      | Property      | Name          | Bedrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents | OneNightPriceCents | OneWeekPriceCents | ThreeMonthsPriceCents | SixMonthsPriceCents | TwelveMonthsPriceCents | FloorPlanImageID |
      | ACME Downtown | Std Apartment | 2        | 50     | 2            | 20    | 12000          | 11000              | 40000             | 350000                | 650000              | 1200000                | 3                |

  Scenario: Updating a unit with invalid attributes
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Property      | Name          | Uploaded |
      | 3  | ACME Downtown | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name     |  |
      | Bedrooms |  |
      | SizeM2   |  |
      | Count    |  |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | name is required                |
      | bedrooms is required            |
      | size is required                |
      | number of unit type is required |
    And I should have the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |

  Scenario: Updating a unit that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | State | Name          |
      | 1  | Other   | Draft | ACME Downtown |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Property      | Name          | Uploaded |
      | 3  | ACME Downtown | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name         | Std Apartment |
      | Bedrooms     | 2             |
      | SizeM2       | 50            |
      | MaxOccupancy | 2             |
      | Count        | 20            |
    Then the system should respond with "NOT FOUND"
    And I should have the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
