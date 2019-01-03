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
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    |

  Scenario: Listing my Bookings
    Given the following bookings:
      | ID | User                 | Unit         | CheckInAt            | CheckOutAt           | Nights |
      | 1  | diego@selzlein.com   | Standard Apt | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | 2  | diego@selzlein.com   | Standard Apt | 2018-06-15T15:00:00Z | 2018-06-16T11:00:00Z | 1      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-07-09T15:00:00Z | 2018-07-11T11:00:00Z | 2      |
    And I am authenticated with "diego@selzlein.com"
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
    And I am authenticated with "diego@selzlein.com"
    When I list my bookings for page "2"
    Then the system should respond with "OK" and the following JSON:
      """
      []
      """
