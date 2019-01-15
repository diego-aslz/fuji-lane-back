Feature: Searching for Units

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following countries:
      | ID | Name  |
      | 1  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 2  | Japan   | Osaka |
      | 3  | Japan   | Tokio |

  Scenario: Searching for units in a city
    Given the following properties:
      | ID | Account          | Name          | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 2           |
      | 2  | Diego Apartments | ACME Uptown   | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 3      | 1         | 3           |
    And the following units:
      | ID | Property      | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          | BasePriceCents |
      | 11 | ACME Downtown | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z | 12000          |
      | 10 | ACME Downtown | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z | 10000          |
      | 12 | ACME Uptown   | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z | 20000          |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
      | 2  | ACME Downtown | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 1        |
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
    When I search for units with the following filters:
      | cityID | 2 |
    Then the system should respond with "OK" and the following JSON:
      """
      [{
        "address1": "88 Tai Tam Reservoir Rd",
        "address2": null,
        "address3": null,
        "amenities": [{
          "name": "Gym",
          "type": "gym"
        }],
        "cityID": 2,
        "countryID": 1,
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
        "units": [{
          "amenities": [{
            "name": "Bathrobes",
            "type": "bathrobes"
          }],
          "perNightPriceCents": 10000,
          "bathrooms": 1,
          "bedrooms": 1,
          "id": 10,
          "images": [{
            "name": "front.jpg",
            "url": "https://s3.amazonaws.com/front.jpg"
          }],
          "maxOccupancy": 3,
          "name": "Standard Apt",
          "sizeM2": 52
        }, {
          "amenities": [],
          "perNightPriceCents": 12000,
          "bathrooms": 2,
          "bedrooms": 2,
          "id": 11,
          "images": [],
          "maxOccupancy": 6,
          "name": "Double Apt",
          "sizeM2": 62
        }]
      }]
      """

  Scenario: Ignored unpublished units
    Given the following properties:
      | ID | Account          | Name          | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 2           |
      | 2  | Diego Apartments | ACME Uptown   | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 2      | 1         | 3           |
    And the following units:
      | ID | Property      | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          | BasePriceCents |
      | 11 | ACME Downtown | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z | 12000          |
    And the following units:
      | ID | Property      | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents |
      | 10 | ACME Downtown | Standard Apt | 1        | 1         | 52     | 3            | 15    | 10000          |
      | 12 | ACME Uptown   | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 20000          |
    When I search for units with the following filters:
      | cityID | 2 |
    Then the system should respond with "OK" and the following JSON:
      """
      [{
        "address1": "88 Tai Tam Reservoir Rd",
        "address2": null,
        "address3": null,
        "amenities": [],
        "cityID": 2,
        "countryID": 1,
        "id": 1,
        "images": [],
        "latitude": 100,
        "longitude": 200,
        "name": "ACME Downtown",
        "overview": "<p>Property Overview</p>",
        "postalCode": null,
        "units": [{
          "amenities": [],
          "perNightPriceCents": 12000,
          "bathrooms": 2,
          "bedrooms": 2,
          "id": 11,
          "images": [],
          "maxOccupancy": 6,
          "name": "Double Apt",
          "sizeM2": 62
        }]
      }]
      """
