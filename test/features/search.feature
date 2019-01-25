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
    And the following properties:
      | ID | Account          | Name                  | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 1  | Diego Apartments | Awesome Property      | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 2           |
      | 2  | Diego Apartments | Nice Property         | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 2      | 1         | 3           |
      | 5  | Diego Apartments | Other City's Property | 2018-06-01T08:00:00Z | <p>Uptown Overview</p>   | 100      | 200       | Add 1                   | 3      | 1         | 3           |
    And the following properties:
      | ID | Account          | Name           |
      | 3  | Diego Apartments | Draft Property |

  Scenario: Searching for units in a city
    Given the following units:
      | ID | Property              | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property      | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 10 | Awesome Property      | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Other City's Property | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Double Apt   | 1         | 12000 |
      | Standard Apt | 1         | 10000 |
      | Triple Apt   | 1         | 20000 |
    And the following images:
      | ID | Property         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | Awesome Property | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
      | 2  | Awesome Property | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 1        |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 10 | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
      | 11 | Standard Apt | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 | 2        |
    And the following amenities:
      | Property         | Type |
      | Awesome Property | gym  |
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
        "name": "Awesome Property",
        "overview": "<p>Property Overview</p>",
        "postalCode": null,
        "slug": "awesome-property",
        "units": [{
          "amenities": [{
            "name": "Bathrobes",
            "type": "bathrobes"
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
          "perNightPriceCents": 10000,
          "sizeM2": 52,
          "slug": "standard-apt"
        }, {
          "amenities": [],
          "bathrooms": 2,
          "bedrooms": 2,
          "id": 11,
          "images": [],
          "maxOccupancy": 6,
          "name": "Double Apt",
          "perNightPriceCents": 12000,
          "sizeM2": 62,
          "slug": "double-apt"
        }]
      }]
      """

  Scenario: Ignoring unpublished units or properties
    Given the following units:
      | ID | Property         | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 13 | Draft Property   | Penthouse  | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 12000 |
    And the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count |
      | 10 | Awesome Property | Standard Apt | 1        | 1         | 52     | 3            | 15    |
      | 12 | Nice Property    | Triple Apt   | 3        | 4         | 80     | 6            | 5     |
    When I search for units with the following filters:
      | cityID | 2 |
    Then the system should respond with "OK" and the following search results:
      | PropertyName     | Name       | PerNightPriceCents |
      | Awesome Property | Double Apt | 12000              |

  Scenario: Paginating listings
    Given the following units:
      | ID | Property         | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    When I search for units with the following filters:
      | cityID | 2 |
      | page   | 2 |
    Then the system should respond with "OK" and the following JSON:
      """
      []
      """

  Scenario: Searching for units with at least 2 bedrooms
    Given the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 10 | Awesome Property | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Awesome Property | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
      | 13 | Nice Property    | Basic Apt    | 1        | 1         | 20     | 1            | 5     | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 12000 |
      | Triple Apt | 1         | 20000 |
    When I search for units with the following filters:
      | cityID   | 2 |
      | bedrooms | 2 |
    Then the system should respond with "OK" and the following search results:
      | PropertyName     | Name       | PerNightPriceCents |
      | Awesome Property | Double Apt | 12000              |
      | Awesome Property | Triple Apt | 20000              |

  Scenario: Searching for units with at least 2 bathrooms
    Given the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 10 | Awesome Property | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Awesome Property | Triple Apt   | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
      | 13 | Nice Property    | Basic Apt    | 1        | 1         | 20     | 1            | 5     | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 12000 |
      | Triple Apt | 1         | 20000 |
    When I search for units with the following filters:
      | cityID    | 2 |
      | bathrooms | 2 |
    Then the system should respond with "OK" and the following search results:
      | PropertyName     | Name       | PerNightPriceCents |
      | Awesome Property | Double Apt | 12000              |
      | Awesome Property | Triple Apt | 20000              |

  Scenario: Obtaining the right price for a specific period of time and validating minimum stay
    Given the following properties:
      | ID | Account          | Name            | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 11 | Diego Apartments | Min 1 Week Stay | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 7           |
    And the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 10 | Awesome Property | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Min 1 Week Stay  | Penthouse    | 3        | 4         | 80     | 6            | 5     | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Double Apt   | 1         | 12000 |
      | Double Apt   | 2         | 20000 |
      | Standard Apt | 1         | 11000 |
      | Penthouse    | 1         | 50000 |
    When I search for units with the following filters:
      | cityID   | 2          |
      | checkIn  | 2019-01-01 |
      | checkOut | 2019-01-03 |
    Then the system should respond with "OK" and the following search results:
      | PropertyName     | Name         | PerNightPriceCents |
      | Awesome Property | Double Apt   | 10000              |
      | Awesome Property | Standard Apt | 11000              |
