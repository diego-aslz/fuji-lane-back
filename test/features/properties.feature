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
      | Account          |
      | Diego Apartments |

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
      | ID | Account          | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview                     | UpdatedAt            |
      | 1  | Diego Apartments | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 3           | 300     | 50       | IGU            | Ines          | Pharmacy        | <strong>Good place!</strong> | 2018-06-05T08:00:00Z |
      | 2  | Other            | Other Prop    | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 4           | 350     | 50       | IGU            | Ines          | Restaurant      | <strong>Nice place!</strong> | 2018-06-05T08:00:00Z |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 2  | ACME Downtown | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 2        |
      | 3  | Other Prop    | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 4  | ACME Downtown | true     | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 1        |
    And the following amenities:
      | ID | Property      | Type |
      | 1  | ACME Downtown | gym  |
    And the following units:
      | ID | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | 11 | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 5  | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 6  | Standard Apt | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 2        |
    And the following amenities:
      | ID | Unit         | Type   |
      | 2  | Standard Apt | toilet |
    When I get details for property "ACME Downtown"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "id": 1,
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": null,
        "everPublished": false,
        "name": "ACME Downtown",
        "address1": "Add. One",
        "address2": "Add. Two",
        "address3": "Add. Three",
        "cityID": 3,
        "postalCode": "223344",
        "countryID": 2,
        "latitude": 0,
        "longitude": 0,
        "minimumStay": 3,
        "deposit": "300",
        "cleaning": "50",
        "nearestAirport": "IGU",
        "nearestSubway": "Ines",
        "nearbyLocations": "Pharmacy",
        "overview": "<strong>Good place!</strong>",
        "images": [
          {
            "id": 4,
            "name": "back.jpg",
            "type": "image/jpeg",
            "size": 1000000,
            "url": "https://s3.amazonaws.com/back.jpg",
            "uploaded": true,
            "position": 1
          }, {
            "id": 1,
            "name": "front.jpg",
            "type": "image/jpeg",
            "size": 1000000,
            "url": "https://s3.amazonaws.com/front.jpg",
            "uploaded": true,
            "position": 2
          }
        ],
        "amenities": [{
          "id": 1,
          "type": "gym",
          "name": "Gym"
        }],
        "units": [{
          "id": 11,
          "publishedAt": null,
          "propertyID": 1,
          "name": "Standard Apt",
          "bedrooms": 1,
          "bathrooms": 0,
          "sizeM2": 52,
          "maxOccupancy": 3,
          "count": 15,
          "basePriceCents": null,
          "oneNightPriceCents": null,
          "oneWeekPriceCents": null,
          "threeMonthsPriceCents": null,
          "sixMonthsPriceCents": null,
          "twelveMonthsPriceCents": null,
          "floorPlanImage": null,
          "overview": null,
          "amenities": [{
            "id": 2,
            "type": "toilet",
            "name": "Toilet"
          }],
          "images": [{
            "id": 5,
            "name": "front.jpg",
            "type": "image/jpeg",
            "size": 1000000,
            "url": "https://s3.amazonaws.com/front.jpg",
            "uploaded": true,
            "position": 2
          }]
        }]
      }
      """

  Scenario: Getting property details for a property the user does not have access to
    Given the following accounts:
      | Name            |
      | John Apartments |
    And the following properties:
      | Account         | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country |
      | John Apartments | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
    When I get details for property "ACME Downtown"
    Then the system should respond with "NOT FOUND"

  Scenario: Updating my property
    Given the following properties:
      | ID | Account          |
      | 1  | Diego Apartments |
    When I update the property "1" with the following details:
      | Name            | ACME Downtown                |
      | Address1        | Add. One                     |
      | Address2        | Add. Two                     |
      | Address3        | Add. Three                   |
      | CityID          | 3                            |
      | PostalCode      | 223344                       |
      | MinimumStay     | 3                            |
      | Deposit         | 150                          |
      | Cleaning        | daily                        |
      | NearestAirport  | IGU                          |
      | NearestSubway   | Central Park                 |
      | NearbyLocations | Pharmacy                     |
      | Overview        | <strong>Good place!</strong> |
      | Latitude        | 34.69374                     |
      | Longitude       | 135.50218                    |
    Then the system should respond with "OK"
    And I should have the following properties:
      | Account          | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview                     | Latitude | Longitude |
      | Diego Apartments | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 3           | 150     | daily    | IGU            | Central Park  | Pharmacy        | <strong>Good place!</strong> | 34.69374 | 135.50218 |

  Scenario: Updating my property with invalid Overview
    Given the following properties:
      | ID | Account          |
      | 1  | Diego Apartments |
    When I update the property "1" with the following details:
      | Overview | <strong>Big windows!</strong><script></script> |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | overview: script tags are not allowed |

  Scenario: Updating a property that does not belong to me
    Given the following accounts:
      | Name             |
      | Other Apartments |
    And the following properties:
      | ID | Account          |
      | 1  | Other Apartments |
    When I update the property "1" with the following details:
      | Name | ACME Downtown |
    Then the system should respond with "NOT FOUND"
    And I should have the following properties:
      | Account          | Name |
      | Other Apartments |      |

  Scenario: Updating property amenities
    Given the following properties:
      | ID | Account          | Name          |
      | 1  | Diego Apartments | ACME Downtown |
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
      | ID | Account          | Name          |
      | 1  | Diego Apartments | ACME Downtown |
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

  Scenario: Publishing my property
    Given the following properties:
      | ID | Account          | Name          | Address1 | Address2 | Address3   | CityID | PostalCode | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview   | Latitude | Longitude | EverPublished |
      | 1  | Diego Apartments | ACME Downtown | Add. One | Add. Two | Add. Three | 3      | 223344     | 3           | 150     | daily    | IGU            | Central Park  | Pharmacy        | Nice place | 34.69374 | 135.50218 | false         |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
    And the following amenities:
      | Property      | Type   | Name      |
      | ACME Downtown | custom | Breakfast |
    And it is currently "05 Jun 18 08:00"
    When I publish property "1"
    Then the system should respond with "OK"
    And I should have the following properties:
      | ID | Account          | Name          | PublishedAt          | EverPublished |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-05T08:00:00Z | true          |

  Scenario: Publishing a property with missing information
    Given the following properties:
      | ID | Account          | EverPublished |
      | 1  | Diego Apartments | false         |
    And the following images:
      | ID | PropertyID | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | 1          | false    | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
    When I publish property "1"
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Name is required                 |
      | Address is incomplete            |
      | At least one amenity is required |
      | At least one image is required   |
    And I should have the following properties:
      | ID | Account          | PublishedAt | EverPublished |
      | 1  | Diego Apartments |             | false         |

  Scenario: Listing my properties
    Given the following accounts:
      | Name              |
      | Antoni Apartments |
    And the following properties:
      | ID | Account           | Name          | Address1                | Address2 | CityID | PostalCode | PublishedAt          | UpdatedAt            | EverPublished |
      | 1  | Diego Apartments  | ACME Downtown | 88 Tai Tam Reservoir Rd | Tai Tam  | 3      | 111        | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | false         |
      | 2  | Diego Apartments  | ACME Uptown   | 90 Tai Tam Reservoir Rd | Tai Tam  | 3      | 222        | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | true          |
      | 3  | Antoni Apartments | ACME          | Add. One                | Add. Two | 3      | 333        | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | false         |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 2  | ACME Downtown | true     | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 1        |
    And the following units:
      | ID | Property      | Name           | BasePriceCents | OneNightPriceCents | OneWeekPriceCents | ThreeMonthsPriceCents | SixMonthsPriceCents | TwelveMonthsPriceCents |
      | 2  | ACME Downtown | Standard Apt   | 10000          | 11000              | 40000             | 350000                | 650000              | 1200000                |
      | 3  | ACME Downtown | Double-bed Apt | 11000          | 12000              | 42000             | 370000                | 670000              | 1220000                |
    When I list my properties
    Then the system should respond with "OK" and the following JSON:
      """
      [{
        "id": 1,
        "name": "ACME Downtown",
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": "2018-06-05T08:00:00Z",
        "everPublished": false,
        "address1": "88 Tai Tam Reservoir Rd",
        "address2": "Tai Tam",
        "address3": null,
        "postalCode": "111",
        "cityID": 3,
        "countryID": null,
        "latitude": 0,
        "longitude": 0,
        "images": [{
          "id": 2,
          "name": "back.jpg",
          "type": "image/jpeg",
          "size": 1000000,
          "url": "https://s3.amazonaws.com/back.jpg",
          "uploaded": true,
          "position": 1
        }, {
          "id": 1,
          "name": "front.jpg",
          "type": "image/jpeg",
          "size": 1000000,
          "url": "https://s3.amazonaws.com/front.jpg",
          "uploaded": true,
          "position": 2
        }],
        "minimumStay": null,
        "deposit": null,
        "cleaning": null,
        "nearestAirport": null,
        "nearestSubway": null,
        "nearbyLocations": null,
        "overview": null,
        "amenities": null,
        "units": [{
          "id": 2,
          "publishedAt": null,
          "propertyID": 1,
          "name": "Standard Apt",
          "bedrooms": 0,
          "bathrooms": 0,
          "sizeM2": 0,
          "maxOccupancy": null,
          "count": 0,
          "basePriceCents": 10000,
          "oneNightPriceCents": 11000,
          "oneWeekPriceCents": 40000,
          "threeMonthsPriceCents": 350000,
          "sixMonthsPriceCents": 650000,
          "twelveMonthsPriceCents": 1200000,
          "floorPlanImage": null,
          "amenities": null,
          "images": null,
          "overview": null
        }, {
          "id": 3,
          "publishedAt": null,
          "propertyID": 1,
          "name": "Double-bed Apt",
          "bedrooms": 0,
          "bathrooms": 0,
          "sizeM2": 0,
          "maxOccupancy": null,
          "count": 0,
          "basePriceCents": 11000,
          "oneNightPriceCents": 12000,
          "oneWeekPriceCents": 42000,
          "threeMonthsPriceCents": 370000,
          "sixMonthsPriceCents": 670000,
          "twelveMonthsPriceCents": 1220000,
          "floorPlanImage": null,
          "amenities": null,
          "images": null,
          "overview": null
        }]
      }, {
        "id": 2,
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": "2018-06-05T08:00:00Z",
        "everPublished": true,
        "name": "ACME Uptown",
        "address1": "90 Tai Tam Reservoir Rd",
        "address2": "Tai Tam",
        "address3": null,
        "postalCode": "222",
        "cityID": 3,
        "countryID": null,
        "latitude": 0,
        "longitude": 0,
        "images": [],
        "minimumStay": null,
        "deposit": null,
        "cleaning": null,
        "nearestAirport": null,
        "nearestSubway": null,
        "nearbyLocations": null,
        "overview": null,
        "amenities": null,
        "units": []
      }]
      """
