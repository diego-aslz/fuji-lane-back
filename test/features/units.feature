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
      | ID | Unit         | Name          | Uploaded |
      | 3  | Standard Apt | blueprint.jpg | true     |
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

  Scenario: Updating a unit with a floor plan image that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account          | State | Name            |
      | 1  | Diego Apartments | Draft | ACME Downtown   |
      | 2  | Other            | Draft | ACME Apartments |
    And the following units:
      | Property        | Name               | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown   | Standard Apt       | 1        | 52     | 3            | 15    |
      | ACME Apartments | Standard Other Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Unit               | Name          | Uploaded |
      | 3  | Standard Other Apt | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name             | Std Apartment |
      | Bedrooms         | 2             |
      | SizeM2           | 50            |
      | MaxOccupancy     | 2             |
      | Count            | 20            |
      | FloorPlanImageID | 3             |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | floor plan image does not exist |
    And I should have the following units:
      | Property        | Name               | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown   | Standard Apt       | 1        | 52     | 3            | 15    |
      | ACME Apartments | Standard Other Apt | 1        | 52     | 3            | 15    |

  Scenario: Updating unit amenities
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
      | Standard Apt | desk      |
    And the following amenities:
      | Unit         | Type   | Name   |
      | Standard Apt | custom | Towels |
      | Standard Apt | custom | Soap   |
    When I update unit "Standard Apt" with the following amenities:
      | Type             | Name             |
      | desk             | Desk             |
      | air_conditioning | Air Conditioning |
      | custom           | Soap             |
      | custom           | Windows          |
    Then the system should respond with "OK"
    And I should have the following amenities:
      | Unit         | Type             | Name    |
      | Standard Apt | desk             |         |
      | Standard Apt | custom           | Soap    |
      | Standard Apt | air_conditioning |         |
      | Standard Apt | custom           | Windows |

  Scenario: Getting unit details
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    And the following units:
      | ID | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents | OneNightPriceCents | OneWeekPriceCents | ThreeMonthsPriceCents | SixMonthsPriceCents | TwelveMonthsPriceCents |
      | 2  | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    | 12000          | 11000              | 40000             | 350000                | 650000              | 1200000                |
    And the following images:
      | ID | Unit         | Name          | Uploaded | URL                                | Type       | Size |
      | 3  | Standard Apt | blueprint.jpg | true     | https://s3.amazonaws.com/front.jpg | image/jpeg | 5000 |
    And unit "Standard Apt" has:
      | FloorPlanImageID | 3 |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
    When I get details for unit "Standard Apt"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "id": 2,
        "propertyID": 1,
        "name": "Standard Apt",
        "bedrooms": 1,
        "sizeM2": 52,
        "maxOccupancy": 3,
        "count": 15,
        "basePriceCents": 12000,
        "oneNightPriceCents": 11000,
        "oneWeekPriceCents": 40000,
        "threeMonthsPriceCents": 350000,
        "sixMonthsPriceCents": 650000,
        "twelveMonthsPriceCents": 1200000,
        "floorPlanImage": {
          "id": 3,
          "name": "blueprint.jpg",
          "type": "image/jpeg",
          "size": 5000,
          "url": "https://s3.amazonaws.com/front.jpg",
          "uploaded": true
        },
        "images": [],
        "amenities": [
          {
            "type": "bathrobes",
            "name": null
          }
        ]
      }
      """

  Scenario: Getting unit details for a unit the user does not have access to
    Given the following accounts:
      | Name            |
      | John Apartments |
    And the following properties:
      | Account         | State | Name          |
      | John Apartments | Draft | ACME Downtown |
    And the following units:
      | ID | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | 2  | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    When I get details for unit "Standard Apt"
    Then the system should respond with "NOT FOUND"
