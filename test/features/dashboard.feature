Feature: Dashboard

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
      | Alex Apartments  |
    And the following properties:
      | Account          | Name          |
      | Diego Apartments | ACME Downtown |
      | Alex Apartments  | ACME Uptown   |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    |
      | ACME Uptown   | Double Apt   | 2        | 52     | 6            | 10    |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
      |                  | antoni@gmail.com   | Antoni               |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Getting dashboard details
    Given the following bookings:
      | User             | Unit         | CreatedAt            | CheckInAt            | CheckOutAt           | Nights |
      | antoni@gmail.com | Standard Apt | 2018-05-31T23:59:00Z | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-01T00:01:00Z | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | antoni@gmail.com | Double Apt   | 2018-06-01T00:01:00Z | 2018-06-19T15:00:00Z | 2018-06-21T11:00:00Z | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-06T23:59:00Z | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-07T00:01:00Z | 2018-06-09T15:00:00Z | 2018-06-11T11:00:00Z | 2      |
    When I get dashboard details for:
      | since | 2018-06-01T00:00:00Z |
      | until | 2018-06-07T00:00:00Z |
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "totals": {
          "searches": 0,
          "visits": 0,
          "requests": 2,
          "favorites": 0
        },
        "daily": [{
          "date": "2018-06-01",
          "bookingsCount": 1,
          "visitsCount": 0
        }, {
          "date": "2018-06-02",
          "bookingsCount": 0,
          "visitsCount": 0
        }, {
          "date": "2018-06-03",
          "bookingsCount": 0,
          "visitsCount": 0
        }, {
          "date": "2018-06-04",
          "bookingsCount": 0,
          "visitsCount": 0
        }, {
          "date": "2018-06-05",
          "bookingsCount": 0,
          "visitsCount": 0
        }, {
          "date": "2018-06-06",
          "bookingsCount": 1,
          "visitsCount": 0
        }]
      }
      """
