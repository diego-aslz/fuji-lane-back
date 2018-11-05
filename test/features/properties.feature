Feature: Properties Management

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Adding a new property
    When I add a new property
    Then the system should respond with "CREATED"
    And I should have the following properties:
      | Account          | State |
      | Diego Apartments | draft |

  Scenario: Adding a new property without having an Account
    Given the following users:
      | Email                | Name             |
      | djeison@selzlein.com | Djeison Selzlein |
    And I am authenticated with "djeison@selzlein.com"
    When I add a new property
    Then the system should respond with "PRECONDITION REQUIRED" and the following errors:
      | You need a company account to perform this action |
    And I should have no properties

  Scenario: Getting property details
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account          | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview   |
      | 1  | Diego Apartments | Draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 3           | 300     | 50       | IGU            | Ines          | Pharmacy        | Good place |
      | 2  | Other            | Draft | Other Prop    | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 4           | 350     | 50       | IGU            | Ines          | Restaurant      | Nice place |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 |
      | 2  | ACME Downtown | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 |
      | 3  | Other Prop    | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 |
    And the following amenities:
      | Property      | Type |
      | ACME Downtown | gym  |
    When I get details for property "ACME Downtown"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "id": 1,
        "state": "draft",
        "name": "ACME Downtown",
        "address1": "Add. One",
        "address2": "Add. Two",
        "address3": "Add. Three",
        "postalCode": "223344",
        "cityID": 3,
        "postalCode": "223344",
        "countryID":  2,
        "minimumStay": "3",
        "deposit": "300",
        "cleaning": "50",
        "nearestAirport": "IGU",
        "nearestSubway": "Ines",
        "nearbyLocations": "Pharmacy",
        "overview": "Good place",
        "images": [
          {
            "id": 1,
            "name": "front.jpg",
            "type": "image/jpeg",
            "size": 1000000,
            "url": "https://s3.amazonaws.com/front.jpg",
            "uploaded":true
          }
        ],
        "amenities": [
          {
            "type": "gym",
            "name": null
          }
        ]
      }
      """

  Scenario: Getting property details for a property the user does not have access to
    Given the following accounts:
      | Name            |
      | John Apartments |
    And the following properties:
      | Account         | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country |
      | John Apartments | Draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
    When I get details for property "ACME Downtown"
    Then the system should respond with "NOT FOUND"

  Scenario: Updating my property
    Given the following properties:
      | ID | Account          | State |
      | 1  | Diego Apartments | Draft |
    When I update the property "1" with the following details:
      | Name            | ACME Downtown |
      | Address1        | Add. One      |
      | Address2        | Add. Two      |
      | Address3        | Add. Three    |
      | CityID          | 3             |
      | PostalCode      | 223344        |
      | MinimumStay     | 3 days        |
      | Deposit         | 150           |
      | Cleaning        | daily         |
      | NearestAirport  | IGU           |
      | NearestSubway   | Central Park  |
      | NearbyLocations | Pharmacy      |
      | Overview        | Nice place    |
    Then the system should respond with "OK"
    And I should have the following properties:
      | Account          | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview   |
      | Diego Apartments | draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 3 days      | 150     | daily    | IGU            | Central Park  | Pharmacy        | Nice place |

  Scenario: Updating a property that does not belong to me
    Given the following accounts:
      | Name             |
      | Other Apartments |
    And the following properties:
      | ID | Account          | State |
      | 1  | Other Apartments | Draft |
    When I update the property "1" with the following details:
      | Name | ACME Downtown |
    Then the system should respond with "NOT FOUND"
    And I should have the following properties:
      | Account          | State | Name |
      | Other Apartments | draft |      |

  Scenario: Updating property amenities
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    And the following amenities:
      | Property      | Type |
      | ACME Downtown | gym  |
      | ACME Downtown | pool |
    And the following amenities:
      | Property      | Type   | Name      |
      | ACME Downtown | custom | Breakfast |
      | ACME Downtown | custom | Casino    |
    When I update the property "1" with the following amenities:
      | Type       | Name          |
      | pool       | Pool          |
      | restaurant | Restaurant    |
      | custom     | Casino        |
      | custom     | All Inclusive |
    Then the system should respond with "OK"
    And I should have the following amenities:
      | Property      | Type       | Name          |
      | ACME Downtown | pool       |               |
      | ACME Downtown | custom     | Casino        |
      | ACME Downtown | restaurant |               |
      | ACME Downtown | custom     | All Inclusive |

  Scenario: Updating property with invalid or duplicated amenities
    Given the following properties:
      | ID | Account          | State | Name          |
      | 1  | Diego Apartments | Draft | ACME Downtown |
    When I update the property "1" with the following amenities:
      | Type    | Name    |
      | invalid | Invalid |
      | custom  | Casino  |
      | custom  | Casino  |
      | custom  |         |
    Then the system should respond with "OK"
    And I should have the following amenities:
      | Property      | Type   | Name   |
      | ACME Downtown | custom | Casino |
