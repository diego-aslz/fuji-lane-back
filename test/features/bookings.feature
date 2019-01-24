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
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Standard Apt | 1         | 11000 |
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
    When I create the following booking:
      | Unit     | Standard Apt |
      | CheckIn  | 2018-06-09   |
      | CheckOut | 2018-06-11   |
      | Message  | Nothing      |
    Then the system should respond with "CREATED"
    And I should have the following bookings:
      | User               | Unit         | CheckIn    | CheckOut   | Message | Nights | NightPriceCents | ServiceFeeCents | TotalCents |
      | diego@selzlein.com | Standard Apt | 2018-06-09 | 2018-06-11 | Nothing | 2      | 11000           | 0               | 22000      |

  Scenario: Booking a Unit with invalid information
    Given it is currently "01 Jun 18 08:00"
    When I create the following booking:
      | Message | Nothing |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | unit is required                   |
      | check in is required               |
      | check in should be in the future   |
      | check out is required              |
      | check out should be after check in |

  Scenario: Trying to Book a Unit that's not published
    Given it is currently "01 Jun 18 08:00"
    And the following units:
      | Property      | Name       | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Double Apt | 1        | 32     | 3            | 15    |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 11000 |
    When I create the following booking:
      | Unit     | Double Apt |
      | CheckIn  | 2018-06-09 |
      | CheckOut | 2018-06-11 |
      | Message  | Nothing    |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | unit is invalid |

  Scenario Outline: Calculating Booking Prices
    Given the following units:
      | Property      | Name            | Bedrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | ACME Downtown | Specific Prices | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
      | ACME Downtown | Single Price    | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
    And the following prices:
      | Unit            | MinNights | Cents   |
      | Specific Prices | 1         | 13000   |
      | Specific Prices | 2         | 20000   |
      | Specific Prices | 7         | 60000   |
      | Specific Prices | 30        | 240000  |
      | Specific Prices | 90        | 750000  |
      | Specific Prices | 180       | 1200000 |
      | Specific Prices | 365       | 2200000 |
      | Single Price    | 1         | 10000   |
    And it is currently "01 Jun 18 08:00"
    When I create the following booking:
      | Unit     | <Unit>     |
      | CheckIn  | <CheckIn>  |
      | CheckOut | <CheckOut> |
    Then the system should respond with "CREATED"
    And I should have the following bookings:
      | User               | Unit   | CheckIn   | CheckOut   | Nights   | NightPriceCents | TotalCents |
      | diego@selzlein.com | <Unit> | <CheckIn> | <CheckOut> | <Nights> | <PerNight>      | <Total>    |

    Examples:
      | Unit            | CheckIn    | CheckOut   | Nights | PerNight | Total   |
      | Specific Prices | 2018-07-01 | 2018-07-02 | 1      | 13000    | 13000   |
      | Specific Prices | 2018-07-01 | 2018-07-07 | 6      | 10000    | 60000   |
      | Specific Prices | 2018-07-01 | 2018-07-08 | 7      | 8571     | 60000   |
      | Specific Prices | 2018-07-01 | 2018-07-30 | 29     | 8571     | 248571  |
      | Specific Prices | 2018-07-01 | 2018-07-31 | 30     | 8000     | 240000  |
      | Specific Prices | 2018-07-01 | 2018-09-28 | 89     | 8000     | 712000  |
      | Specific Prices | 2018-07-01 | 2018-09-29 | 90     | 8333     | 750000  |
      | Specific Prices | 2018-07-01 | 2018-12-27 | 179    | 8333     | 1491667 |
      | Specific Prices | 2018-07-01 | 2018-12-28 | 180    | 6667     | 1200000 |
      | Specific Prices | 2018-07-01 | 2019-06-30 | 364    | 6667     | 2426667 |
      | Specific Prices | 2018-07-01 | 2019-07-01 | 365    | 6027     | 2200000 |
      | Single Price    | 2018-07-01 | 2018-07-02 | 1      | 10000    | 10000   |
      | Single Price    | 2018-07-01 | 2018-07-08 | 7      | 10000    | 70000   |
      | Single Price    | 2018-07-01 | 2018-12-28 | 180    | 10000    | 1800000 |
      | Single Price    | 2018-07-01 | 2019-07-01 | 365    | 10000    | 3650000 |
