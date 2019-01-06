Feature: Bookings

  Background:
    Given the following users:
      | Email                |
      | diego@selzlein.com   |
      | djeison@selzlein.com |
    And the following accounts:
      | Name             |
      | Diego Apartments |
    And the following properties:
      | Account          | Name          |
      | Diego Apartments | ACME Downtown |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents | PublishedAt          |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    | 11000          | 2018-06-09T15:00:00Z |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Listing my Bookings
    Given the following bookings:
      | ID | User                 | Unit         | CheckIn    | CheckOut   | Nights |
      | 1  | diego@selzlein.com   | Standard Apt | 2018-06-09 | 2018-06-11 | 2      |
      | 2  | diego@selzlein.com   | Standard Apt | 2018-06-15 | 2018-06-16 | 1      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-07-09 | 2018-07-11 | 2      |
    When I list my bookings
    Then the system should respond with "OK" and the following JSON:
      """
      [{
        "id": 2,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-15",
        "checkOut": "2018-06-16",
        "nights": 1
      }, {
        "id": 1,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-09",
        "checkOut": "2018-06-11",
        "nights": 2
      }]
      """

  Scenario: Paginating my Bookings
    Given the following bookings:
      | ID | User               | Unit         | CheckIn    | CheckOut   | Nights |
      | 1  | diego@selzlein.com | Standard Apt | 2018-06-09 | 2018-06-11 | 2      |
      | 2  | diego@selzlein.com | Standard Apt | 2018-06-15 | 2018-06-16 | 1      |
    When I list my bookings for page "2"
    Then the system should respond with "OK" and the following JSON:
      """
      []
      """

  Scenario: Booking a Unit
    Given it is currently "01 Jun 18 08:00"
    When I submit the following booking:
      | Unit           | Standard Apt |
      | CheckIn        | 2018-06-09   |
      | CheckOut       | 2018-06-11   |
      | AdditionalInfo | Nothing      |
    Then the system should respond with "CREATED"
    And I should have the following bookings:
      | User               | Unit         | CheckIn    | CheckOut   | AdditionalInfo | Nights | NightPriceCents | ServiceFeeCents | TotalCents |
      | diego@selzlein.com | Standard Apt | 2018-06-09 | 2018-06-11 | Nothing        | 2      | 11000           | 0               | 22000      |

  Scenario: Booking a Unit with invalid information
    Given it is currently "01 Jun 18 08:00"
    When I submit the following booking:
      | AdditionalInfo | Nothing |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | unit is required                             |
      | check in date is required                    |
      | check in date should be in the future        |
      | check out date is required                   |
      | check out date should be after check in date |

  Scenario: Trying to Book a Unit that's not published
    Given it is currently "01 Jun 18 08:00"
    And the following units:
      | Property      | Name       | Bedrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents |
      | ACME Downtown | Double Apt | 1        | 32     | 3            | 15    | 11000          |
    When I submit the following booking:
      | Unit           | Double Apt |
      | CheckIn        | 2018-06-09 |
      | CheckOut       | 2018-06-11 |
      | AdditionalInfo | Nothing    |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | unit is invalid |
