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
      | ID | User                 | Unit         | CheckInAt            | CheckOutAt           | Nights |
      | 1  | diego@selzlein.com   | Standard Apt | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | 2  | diego@selzlein.com   | Standard Apt | 2018-06-15T15:00:00Z | 2018-06-16T11:00:00Z | 1      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-07-09T15:00:00Z | 2018-07-11T11:00:00Z | 2      |
    When I list my bookings
    Then the system should respond with "OK" and the following JSON:
      """
      [{
        "id": 2,
        "unitName": "Standard Apt",
        "checkInAt": "2018-06-15T15:00:00Z",
        "checkOutAt": "2018-06-16T11:00:00Z",
        "nights": 1
      }, {
        "id": 1,
        "unitName": "Standard Apt",
        "checkInAt": "2018-06-09T15:00:00Z",
        "checkOutAt": "2018-06-11T11:00:00Z",
        "nights": 2
      }]
      """

  Scenario: Paginating my Bookings
    Given the following bookings:
      | ID | User               | Unit         | CheckInAt            | CheckOutAt           | Nights |
      | 1  | diego@selzlein.com | Standard Apt | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | 2  | diego@selzlein.com | Standard Apt | 2018-06-15T15:00:00Z | 2018-06-16T11:00:00Z | 1      |
    When I list my bookings for page "2"
    Then the system should respond with "OK" and the following JSON:
      """
      []
      """

  Scenario: Booking a Unit
    Given it is currently "01 Jun 18 08:00"
    When I submit the following booking:
      | Unit           | Standard Apt         |
      | CheckInAt      | 2018-06-09T15:00:00Z |
      | CheckOutAt     | 2018-06-11T11:00:00Z |
      | AdditionalInfo | Nothing              |
    Then the system should respond with "CREATED"
    And I should have the following bookings:
      | User               | Unit         | CheckInAt            | CheckOutAt           | AdditionalInfo | Nights | NightPriceCents | ServiceFeeCents | TotalCents |
      | diego@selzlein.com | Standard Apt | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | Nothing        | 2      | 11000           | 0               | 22000      |

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
