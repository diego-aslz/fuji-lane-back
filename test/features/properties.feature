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

  Scenario: Creating a new property
    When I create the following property:
      | Name       | ACME Downtown                |
      | Overview   | <strong>Good place!</strong> |
      | Address1   | Add. One                     |
      | Address2   | Add. Two                     |
      | Address3   | Add. Three                   |
      | CityID     | 3                            |
      | PostalCode | 223344                       |
      | Latitude   | 34.69374                     |
      | Longitude  | 135.50218                    |
    Then I should receive a "CREATED" response
    And I should have the following properties:
      | Account          | Name          | Slug          | Address1 | Address2 | Address3   | City  | PostalCode | Country | Overview                     | Latitude | Longitude |
      | Diego Apartments | ACME Downtown | acme-downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | <strong>Good place!</strong> | 34.69374 | 135.50218 |

  Scenario: Creating a new property with invalid data
    When I create the following property:
      | Overview | <strong>Good place!</strong> |
    Then I should receive a "UNPROCESSABLE ENTITY" response with the following errors:
      | name is required    |
      | address is required |
      | city is required    |

  Scenario: Adding a new property without having an Account
    Given the following users:
      | Email                | Name             |
      | djeison@selzlein.com | Djeison Selzlein |
    And I am authenticated with "djeison@selzlein.com"
    When I create the following property:
      | Name | ACME Downtown |
    Then I should receive a "PRECONDITION REQUIRED" response with the following errors:
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
      | ID | Property      | Name         | Bedrooms | SizeM2 | SizeFT2 | MaxOccupancy | Count |
      | 11 | ACME Downtown | Standard Apt | 1        | 52     | 560     | 3            | 15    |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 5  | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 6  | Standard Apt | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 2        |
    And the following amenities:
      | ID | Unit         | Type |
      | 2  | Standard Apt | desk |
    When I get details for property "ACME Downtown"
    Then I should receive an "OK" response with the following JSON:
      """
      {
        "id": 1,
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": null,
        "firstPublishedAt": null,
        "name": "ACME Downtown",
        "slug": "acme-downtown",
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
          "firstPublishedAt": null,
          "propertyID": 1,
          "name": "Standard Apt",
          "slug": "standard-apt",
          "bedrooms": 1,
          "bathrooms": 0,
          "sizeM2": 52,
          "sizeFT2": 560,
          "maxOccupancy": 3,
          "count": 15,
          "prices": [],
          "floorPlanImage": null,
          "overview": null,
          "amenities": [{
            "id": 2,
            "type": "desk",
            "name": "Desk"
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
    Then I should receive a "NOT FOUND" response

  Scenario: Updating my property without changing amenities
    Given the following properties:
      | ID | Account          | Country | City  |
      | 1  | Diego Apartments | Japan   | Osaka |
    And the following amenities:
      | PropertyID | Type |
      | 1          | gym  |
      | 1          | pool |
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
    Then I should receive an "OK" response
    And I should have the following properties:
      | Account          | Name          | Slug          | Address1 | Address2 | Address3   | City  | PostalCode | Country | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview                     | Latitude | Longitude |
      | Diego Apartments | ACME Downtown | acme-downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   | 3           | 150     | daily    | IGU            | Central Park  | Pharmacy        | <strong>Good place!</strong> | 34.69374 | 135.50218 |
    And I should have the following amenities:
      | Property      | Type |
      | ACME Downtown | gym  |
      | ACME Downtown | pool |

  Scenario: Updating my property with an used name
    Given the following accounts:
      | Name             |
      | Other Apartments |
    And the following properties:
      | ID | Account          | Name | Country | City  |
      | 1  | Diego Apartments |      | Japan   | Osaka |
      | 2  | Other Apartments | ACME | Japan   | Osaka |
    When I update the property "1" with the following details:
      | Name | ACME |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is already in use |

  Scenario: Updating my property with a name that will create a duplicated slug
    Given the following accounts:
      | Name             |
      | Other Apartments |
    And the following properties:
      | ID | Account          | Name      | Country | City  |
      | 1  | Diego Apartments |           | Japan   | Osaka |
      | 2  | Other Apartments | ACME Down | Japan   | Osaka |
    When I update the property "1" with the following details:
      | Name | ACME  Down |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is already in use |

  Scenario: Updating my property with invalid Overview
    Given the following properties:
      | ID | Account          | Country | City  |
      | 1  | Diego Apartments | Japan   | Osaka |
    When I update the property "1" with the following details:
      | Overview | <strong>Big windows!</strong><script></script> |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | overview: script tags are not allowed |

  Scenario: Updating a property that does not belong to me
    Given the following accounts:
      | Name             |
      | Other Apartments |
    And the following properties:
      | ID | Account          | Country | City  |
      | 1  | Other Apartments | Japan   | Osaka |
    When I update the property "1" with the following details:
      | Name | ACME Downtown |
    Then I should receive a "NOT FOUND" response
    And I should have the following properties:
      | Account          | Name |
      | Other Apartments |      |

  Scenario: Updating property amenities
    Given the following properties:
      | ID | Account          | Name          | Country | City  |
      | 1  | Diego Apartments | ACME Downtown | Japan   | Osaka |
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
    Then I should receive an "OK" response
    And I should have the following amenities:
      | Property      | Type       | Name          |
      | ACME Downtown | pool       |               |
      | ACME Downtown | custom     | Casino        |
      | ACME Downtown | restaurant |               |
      | ACME Downtown | custom     | All Inclusive |

  Scenario: Updating property with invalid or duplicated amenities
    Given the following properties:
      | ID | Account          | Name          | Country | City  |
      | 1  | Diego Apartments | ACME Downtown | Japan   | Osaka |
    When I update the property "1" with the following amenities:
      | Type    | Name    |
      | invalid | Invalid |
      | custom  | Casino  |
      | custom  | Casino  |
      | custom  |         |
    Then I should receive an "OK" response
    And I should have the following amenities:
      | Property      | Type   | Name   |
      | ACME Downtown | custom | Casino |

  Scenario Outline: Publishing my property
    Given the following properties:
      | ID | Account          | Name          | Address1 | Address2 | Address3   | Country | City  | PostalCode | MinimumStay | Deposit | Cleaning | NearestAirport | NearestSubway | NearbyLocations | Overview   | Latitude | Longitude | FirstPublishedAt         |
      | 1  | Diego Apartments | ACME Downtown | Add. One | Add. Two | Add. Three | Japan   | Osaka | 223344     | 3           | 150     | daily    | IGU            | Central Park  | Pharmacy        | Nice place | 34.69374 | 135.50218 | <FirstPublishedAtBefore> |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
    And the following amenities:
      | Property      | Type   | Name      |
      | ACME Downtown | custom | Breakfast |
    And it is currently "05 Jun 18 08:00"
    When I publish property "1"
    Then I should receive an "OK" response
    And I should have the following properties:
      | ID | Account          | Name          | PublishedAt          | FirstPublishedAt        |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-05T08:00:00Z | <FirstPublishedAtAfter> |

    Examples:
      | FirstPublishedAtBefore | FirstPublishedAtAfter |
      |                        | 2018-06-05T08:00:00Z  |
      | 2015-06-05T08:00:00Z   | 2015-06-05T08:00:00Z  |

  Scenario: Publishing a property with missing information
    Given the following properties:
      | ID | Account          | Country | City  |
      | 1  | Diego Apartments | Japan   | Osaka |
    And the following images:
      | ID | PropertyID | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | 1          | false    | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
    When I publish property "1"
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is required                 |
      | Address is incomplete            |
      | At least one amenity is required |
      | At least one image is required   |
    And I should have the following properties:
      | ID | Account          | PublishedAt | FirstPublishedAt |
      | 1  | Diego Apartments |             |                  |

  Scenario: Unpublishing my property
    Given the following properties:
      | ID | Account          | Name          | PublishedAt          | Country | City  |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-05T08:00:00Z | Japan   | Osaka |
    When I unpublish property "1"
    Then I should receive an "OK" response
    And I should have the following properties:
      | ID | Account          | Name          | PublishedAt |
      | 1  | Diego Apartments | ACME Downtown |             |

  Scenario: Listing my properties
    Given the following accounts:
      | Name              |
      | Antoni Apartments |
    And the following properties:
      | ID | Account           | Name          | Address1                | Address2 | Country | City  | PostalCode | PublishedAt          | FirstPublishedAt     | UpdatedAt            |
      | 1  | Diego Apartments  | ACME Downtown | 88 Tai Tam Reservoir Rd | Tai Tam  | Japan   | Osaka | 111        | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z |
      | 2  | Diego Apartments  | ACME Uptown   | 90 Tai Tam Reservoir Rd | Tai Tam  | Japan   | Osaka | 222        | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z |
      | 3  | Antoni Apartments | ACME          | Add. One                | Add. Two | Japan   | Osaka | 333        | 2018-06-05T08:00:00Z |                      | 2018-06-05T08:00:00Z |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
      | 2  | ACME Downtown | true     | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 1        |
    And the following units:
      | ID | Property      | Name           | FirstPublishedAt     |
      | 2  | ACME Downtown | Standard Apt   |                      |
      | 3  | ACME Downtown | Double-bed Apt | 2018-06-05T08:00:00Z |
    And the following prices:
      | Unit           | MinNights | Cents |
      | Standard Apt   | 1         | 11000 |
      | Standard Apt   | 2         | 20000 |
      | Double-bed Apt | 1         | 12000 |
      | Double-bed Apt | 2         | 22000 |
    When I list my properties
    Then I should receive an "OK" response with the following JSON:
      """
      [{
        "id": 1,
        "name": "ACME Downtown",
        "slug": "acme-downtown",
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": "2018-06-05T08:00:00Z",
        "firstPublishedAt": "2018-06-05T08:00:00Z",
        "address1": "88 Tai Tam Reservoir Rd",
        "address2": "Tai Tam",
        "address3": null,
        "postalCode": "111",
        "cityID": 3,
        "countryID": 2,
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
          "firstPublishedAt": null,
          "propertyID": 1,
          "name": "Standard Apt",
          "slug": "standard-apt",
          "bedrooms": 0,
          "bathrooms": 0,
          "sizeM2": 0,
          "sizeFT2": 0,
          "maxOccupancy": null,
          "count": 0,
          "prices": [{
            "minNights": 1,
            "cents": 11000
          }, {
            "minNights": 2,
            "cents": 20000
          }],
          "floorPlanImage": null,
          "amenities": null,
          "images": null,
          "overview": null
        }, {
          "id": 3,
          "publishedAt": null,
          "firstPublishedAt": "2018-06-05T08:00:00Z",
          "propertyID": 1,
          "name": "Double-bed Apt",
          "slug": "double-bed-apt",
          "bedrooms": 0,
          "bathrooms": 0,
          "sizeM2": 0,
          "sizeFT2": 0,
          "maxOccupancy": null,
          "count": 0,
          "prices": [{
            "minNights": 1,
            "cents": 12000
          }, {
            "minNights": 2,
            "cents": 22000
          }],
          "floorPlanImage": null,
          "amenities": null,
          "images": null,
          "overview": null
        }]
      }, {
        "id": 2,
        "updatedAt": "2018-06-05T08:00:00Z",
        "publishedAt": "2018-06-05T08:00:00Z",
        "firstPublishedAt": "2018-06-05T08:00:00Z",
        "name": "ACME Uptown",
        "slug": "acme-uptown",
        "address1": "90 Tai Tam Reservoir Rd",
        "address2": "Tai Tam",
        "address3": null,
        "postalCode": "222",
        "cityID": 3,
        "countryID": 2,
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
