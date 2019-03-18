Feature: Dashboard

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
      | Alex Apartments  |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And the following properties:
      | Account          | Name          | Country | City  |
      | Diego Apartments | ACME Downtown | Japan   | Osaka |
      | Alex Apartments  | ACME Uptown   | Japan   | Osaka |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    |
      | ACME Uptown   | Double Apt   | 2        | 52     | 6            | 10    |
    And the following users:
      | Account          | Email                | Name                 |
      | Diego Apartments | diego@selzlein.com   | Diego Aguir Selzlein |
      |                  | antoni@gmail.com     | Antoni               |
      |                  | djeison@selzlein.com | Djeison              |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Getting dashboard details
    Given the following bookings:
      | User             | Unit         | CreatedAt            | CheckIn    | CheckOut   | Nights |
      | antoni@gmail.com | Standard Apt | 2018-05-31T23:59:00Z | 2018-06-09 | 2018-06-11 | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-01T00:01:00Z | 2018-06-09 | 2018-06-11 | 2      |
      | antoni@gmail.com | Double Apt   | 2018-06-01T00:01:00Z | 2018-06-19 | 2018-06-21 | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-07T23:59:00Z | 2018-06-09 | 2018-06-11 | 2      |
      | antoni@gmail.com | Standard Apt | 2018-06-08T00:01:00Z | 2018-06-09 | 2018-06-11 | 2      |
    When I get dashboard details for:
      | since | 2018-06-01 |
      | until | 2018-06-07 |
    Then I should receive an "OK" response with the following JSON:
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
          "bookingsCount": 0,
          "visitsCount": 0
        }, {
          "date": "2018-06-07",
          "bookingsCount": 1,
          "visitsCount": 0
        }]
      }
      """

  Scenario: Listing my Properties' Bookings
    Given the following bookings:
      | ID | User                 | Unit         | CheckIn    | CheckOut   | Nights | PerNightCents | TotalCents |
      | 1  | antoni@gmail.com     | Standard Apt | 2018-06-09 | 2018-06-11 | 2      | 10000         | 20000      |
      | 2  | antoni@gmail.com     | Standard Apt | 2018-06-15 | 2018-06-16 | 1      | 10000         | 10000      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-05-09 | 2018-05-11 | 2      | 10000         | 20000      |
      | 4  | djeison@selzlein.com | Double Apt   | 2018-07-19 | 2018-07-21 | 2      | 10000         | 20000      |
    When I list my properties' bookings
    Then I should receive an "OK" response with the following JSON:
      """
      [{
        "id": 3,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-05-09",
        "checkOut": "2018-05-11",
        "nights": 2,
        "perNightCents": 10000,
        "totalCents": 20000,
        "user": {
          "name": "Djeison",
          "email": "djeison@selzlein.com"
        }
      }, {
        "id": 2,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-15",
        "checkOut": "2018-06-16",
        "nights": 1,
        "perNightCents": 10000,
        "totalCents": 10000,
        "user": {
          "name": "Antoni",
          "email": "antoni@gmail.com"
        }
      }, {
        "id": 1,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-09",
        "checkOut": "2018-06-11",
        "nights": 2,
        "perNightCents": 10000,
        "totalCents": 20000,
        "user": {
          "name": "Antoni",
          "email": "antoni@gmail.com"
        }
      }]
      """

  Scenario: Paginating my properties' Bookings
    Given the following bookings:
      | ID | User                 | Unit         | CheckIn    | CheckOut   | Nights |
      | 1  | antoni@gmail.com     | Standard Apt | 2018-06-09 | 2018-06-11 | 2      |
      | 2  | antoni@gmail.com     | Standard Apt | 2018-06-15 | 2018-06-16 | 1      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-05-09 | 2018-05-11 | 2      |
      | 4  | djeison@selzlein.com | Double Apt   | 2018-07-19 | 2018-07-21 | 2      |
    When I list my properties' bookings for page "2"
    Then I should receive an "OK" response with the following JSON:
      """
      []
      """
