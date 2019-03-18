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
      | ID | Account          | Name               | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 1  | Diego Apartments | ACME Downtown      | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         | 2           |
      | 2  | Diego Apartments | ACME Uptown        | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 3      | 2         | 3           |
      | 3  | Diego Apartments | No Published Units | 2018-06-01T08:00:00Z | <p>Overview</p>          | 100      | 200       | Add 1                   | 3      | 2         | 4           |
      | 4  | Diego Apartments | Different City     | 2018-06-01T08:00:00Z | <p>Overview</p>          | 100      | 200       | Add 1                   | 2      | 1         | 5           |
      | 10 | Diego Apartments | Not Published      |                      | <p>Overview</p>          | 100      | 200       |                         | 3      | 2         | 1           |
    And the following units:
      | ID | Property           | Name                  | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | ACME Downtown      | Double Apt            | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 10 | ACME Downtown      | Standard Apt          | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | ACME Uptown        | Triple Apt            | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
      | 13 | Different City     | Penthouse             | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
      | 22 | ACME Downtown      | Unpublished Penthouse | 1        | 1         | 1      | 1            | 1     |                      |
      | 23 | No Published Units | Unpublished           | 1        | 1         | 1      | 1            | 1     |                      |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Double Apt   | 1         | 12000 |
      | Standard Apt | 1         | 10000 |
      | Triple Apt   | 1         | 20000 |
      | Penthouse    | 1         | 20000 |
    And the following images:
      | ID | Property      | Unit         | Uploaded | Name          | URL                                    | Type       | Size    | Position |
      | 1  | ACME Downtown |              | true     | front.jpg     | https://s3.amazonaws.com/front.jpg     | image/jpeg | 1000000 | 1        |
      | 2  | ACME Downtown |              | false    | back.jpg      | https://s3.amazonaws.com/back.jpg      | image/jpeg | 1000000 | 2        |
      | 3  | ACME Uptown   |              | true     | property.jpg  | https://s3.amazonaws.com/property.jpg  | image/jpeg | 1000000 | 1        |
      | 4  | ACME Uptown   |              | true     | reception.jpg | https://s3.amazonaws.com/reception.jpg | image/jpeg | 1000000 | 2        |
      | 10 |               | Standard Apt | true     | front.jpg     | https://s3.amazonaws.com/front.jpg     | image/jpeg | 1000000 | 1        |
      | 11 |               | Standard Apt | false    | back.jpg      | https://s3.amazonaws.com/back.jpg      | image/jpeg | 1000000 | 2        |
    And the following amenities:
      | Property      | Unit         | Type      |
      | ACME Downtown |              | gym       |
      |               | Standard Apt | bathrobes |
    When I get listing details for "ACME Downtown"
    Then I should receive an "OK" response with the following JSON:
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
        "minimumStay": 2,
        "name": "ACME Downtown",
        "overview": "<p>Property Overview</p>",
        "postalCode": null,
        "similarListings": [{
          "address1": "Add 1",
          "address2": null,
          "address3": null,
          "prices": [{
            "minNights": 1,
            "cents": 20000
          }],
          "bathrooms": 4,
          "bedrooms": 3,
          "id": 2,
          "images": [{
            "name": "property.jpg",
            "url": "https://s3.amazonaws.com/property.jpg"
          }, {
            "name": "reception.jpg",
            "url": "https://s3.amazonaws.com/reception.jpg"
          }],
          "name": "ACME Uptown",
          "overview": "<p>Uptown Overview</p>",
          "sizeM2": 80,
          "slug": "acme-uptown"
        }],
        "slug": "acme-downtown",
        "units": [{
          "amenities": [{
            "name": "Bathrobes",
            "type": "bathrobes"
          }],
          "prices": [{
            "minNights": 1,
            "cents": 10000
          }],
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
          "sizeM2": 52,
          "slug": "standard-apt"
        }, {
          "amenities": [],
          "prices": [{
            "minNights": 1,
            "cents": 12000
          }],
          "bathrooms": 2,
          "bedrooms": 2,
          "id": 11,
          "images": [],
          "maxOccupancy": 6,
          "name": "Double Apt",
          "overview": null,
          "sizeM2": 62,
          "slug": "double-apt"
        }]
      }
      """

  Scenario: Getting Listing Details for a not published Property
    Given the following properties:
      | ID | Account          | Name          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    When I get listing details for "ACME Downtown"
    Then I should receive a "NOT FOUND" response

  Scenario: Getting Listing Details for a published Property with no published Units
    Given the following properties:
      | ID | Account          | Name          | Overview                 | PublishedAt          | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 2018-06-01T08:00:00Z | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    |
    When I get listing details for "ACME Downtown"
    Then I should receive a "NOT FOUND" response

  Scenario: Getting Listing Details for a not published Property as its owner
    Given the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    And the following properties:
      | ID | Account          | Name          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    When I get listing details for "ACME Downtown"
    Then I should receive an "OK" response

  Scenario: Getting Listing Details for a published Property with no published Units as its owner
    Given the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    And the following properties:
      | ID | Account          | Name          | Overview                 | PublishedAt          | Latitude | Longitude | Address1                | CityID | CountryID |
      | 1  | Diego Apartments | ACME Downtown | <p>Property Overview</p> | 2018-06-01T08:00:00Z | 100      | 200       | 88 Tai Tam Reservoir Rd | 3      | 2         |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    |
    When I get listing details for "ACME Downtown"
    Then I should receive an "OK" response
