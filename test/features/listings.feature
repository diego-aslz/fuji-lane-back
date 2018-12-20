Feature: Listings

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following countries:
      | ID | Name      |
      | 1  | Japan     |
      | 2  | Hong Kong |
    And the following cities:
      | ID | Country   | Name      |
      | 2  | Japan     | Osaka     |
      | 3  | Hong Kong | Hong Kong |

  Scenario: Getting Listing Details
    Given the following properties:
      | ID | Account          | Name               | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown      | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
      | 2  | Diego Apartments | ACME Uptown        | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 3      | 2         |
      | 3  | Diego Apartments | No Published Units | 2018-06-01T08:00:00Z | <p>Overview</p>          | 100      | 200       | Add 1                   | 3      | 2         |
      | 4  | Diego Apartments | Different City     | 2018-06-01T08:00:00Z | <p>Overview</p>          | 100      | 200       | Add 1                   | 2      | 1         |
    And the following properties:
      | ID | Account          | Name          | Overview        | Latitude | Longitude | CityID | CountryID |
      | 10 | Diego Apartments | Not Published | <p>Overview</p> | 100      | 200       | 3      | 2         |
    And the following units:
      | ID | Property       | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          | BasePriceCents |
      | 11 | ACME Downtown  | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z | 12000          |
      | 10 | ACME Downtown  | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z | 10000          |
      | 12 | ACME Uptown    | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z | 20000          |
      | 13 | Different City | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z | 20000          |
    And the following units:
      | ID | Property           | Name                  |
      | 22 | ACME Downtown      | Unpublished Penthouse |
      | 23 | No Published Units | Unpublished           |
    And the following images:
      | ID | Property      | Uploaded | Name         | URL                                   | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg    | https://s3.amazonaws.com/front.jpg    | image/jpeg | 1000000 | 1        |
      | 2  | ACME Downtown | false    | back.jpg     | https://s3.amazonaws.com/back.jpg     | image/jpeg | 1000000 | 2        |
      | 3  | ACME Uptown   | true     | property.jpg | https://s3.amazonaws.com/property.jpg | image/jpeg | 1000000 | 1        |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 10 | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
      | 11 | Standard Apt | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 2        |
    And the following amenities:
      | Property      | Type |
      | ACME Downtown | gym  |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
    When I get listing details for "ACME Downtown"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "address1": "88 Tai Tam Reservoir Rd",
        "address2": null,
        "address3": null,
        "amenities": [{
          "name": "Gym",
          "type": "gym"
        }],
        "cityID": 3,
        "countryID": 2,
        "id": 1,
        "images": [{
          "name": "front.jpg",
          "url": "https://s3.amazonaws.com/front.jpg"
        }],
        "latitude": 100,
        "longitude": 200,
        "name": "ACME Downtown",
        "overview": "<p>Property Overview</p>",
        "postalCode": null,
        "similarListings": [{
          "address1": "Add 1",
          "address2": null,
          "address3": null,
          "basePriceCents": 20000,
          "bathrooms": 4,
          "bedrooms": 3,
          "id": 2,
          "images": [{
            "name": "property.jpg",
            "url": "https://s3.amazonaws.com/property.jpg"
          }],
          "name": "ACME Uptown",
          "overview": "<p>Uptown Overview</p>",
          "sizeM2": 80
        }],
        "units": [{
          "amenities": [{
            "name": "Bathrobes",
            "type": "bathrobes"
          }],
          "basePriceCents": 10000,
          "bathrooms": 1,
          "bedrooms": 1,
          "id": 10,
          "images": [{
            "name": "front.jpg",
            "url": "https://s3.amazonaws.com/front.jpg"
          }],
          "maxOccupancy": 3,
          "name": "Standard Apt",
          "overview": null,
          "sizeM2": 52
        }, {
          "amenities": [],
          "basePriceCents": 12000,
          "bathrooms": 2,
          "bedrooms": 2,
          "id": 11,
          "images": [],
          "maxOccupancy": 6,
          "name": "Double Apt",
          "overview": null,
          "sizeM2": 62
        }]
      }
      """

  Scenario: Getting Listing Details for a not published Property
    Given the following properties:
      | ID | Account          | Name          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          | BasePriceCents |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z | 12000          |
    When I get listing details for "ACME Downtown"
    Then the system should respond with "NOT FOUND"

  Scenario: Getting Listing Details for a published Property with no published Units
    Given the following properties:
      | ID | Account          | Name          | Overview                 | PublishedAt          | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 2018-06-01T08:00:00Z | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 12000          |
    When I get listing details for "ACME Downtown"
    Then the system should respond with "NOT FOUND"

  Scenario: Getting Listing Details for a not published Property as its owner
    Given the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    And the following properties:
      | ID | Account          | Name          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          | BasePriceCents |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z | 12000          |
    When I get listing details for "ACME Downtown"
    Then the system should respond with "OK"

  Scenario: Getting Listing Details for a published Property with no published Units as its owner
    Given the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    And the following properties:
      | ID | Account          | Name          | Overview                 | PublishedAt          | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 2018-06-01T08:00:00Z | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 12000          |
    When I get listing details for "ACME Downtown"
    Then the system should respond with "OK"
