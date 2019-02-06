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
      | ID | Account          | Name           | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 3  | Diego Apartments | Draft Property | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 1           |

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
    Then I should receive an "OK" response with the following JSON:
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
          "perNightCents": 10000,
          "totalCents": 10000,
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
          "perNightCents": 12000,
          "totalCents": 12000,
          "sizeM2": 62,
          "slug": "double-apt"
        }]
      }]
      """
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 10000 |
      | Max-Per-Night-Cents    | 12000 |
      | Avg-Per-Night-Cents    | 11000 |

  Scenario: Ignoring unpublished units or properties
    Given the following units:
      | ID | Property         | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 13 | Draft Property   | Penthouse  | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    And the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count |
      | 10 | Awesome Property | Standard Apt | 1        | 1         | 52     | 3            | 15    |
      | 12 | Nice Property    | Triple Apt   | 3        | 4         | 80     | 6            | 5     |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Standard Apt | 1         | 10000 |
      | Double Apt   | 1         | 12000 |
      | Triple Apt   | 1         | 15000 |
      | Penthouse    | 1         | 80000 |
    When I search for units with the following filters:
      | cityID | 2 |
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name       | PerNightCents | TotalCents |
      | Awesome Property | Double Apt | 12000         | 12000      |
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 12000 |
      | Max-Per-Night-Cents    | 12000 |
      | Avg-Per-Night-Cents    | 12000 |

  Scenario: Paginating listings
    Given the following units:
      | ID | Property         | Name       | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 11 | Awesome Property | Double Apt | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 12000 |
    When I search for units with the following filters:
      | cityID | 2 |
      | page   | 2 |
    Then I should receive an "OK" response with the following JSON:
      """
      []
      """
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 12000 |
      | Max-Per-Night-Cents    | 12000 |
      | Avg-Per-Night-Cents    | 12000 |

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
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name       | PerNightCents | TotalCents |
      | Awesome Property | Double Apt | 12000         | 12000      |
      | Awesome Property | Triple Apt | 20000         | 20000      |
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 12000 |
      | Max-Per-Night-Cents    | 20000 |
      | Avg-Per-Night-Cents    | 16000 |

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
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name       | PerNightCents | TotalCents |
      | Awesome Property | Double Apt | 12000         | 12000      |
      | Awesome Property | Triple Apt | 20000         | 20000      |
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 12000 |
      | Max-Per-Night-Cents    | 20000 |
      | Avg-Per-Night-Cents    | 16000 |

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
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name         | PerNightCents | TotalCents |
      | Awesome Property | Double Apt   | 10000         | 20000      |
      | Awesome Property | Standard Apt | 11000         | 22000      |
    And I should receive the following headers:
      | Total-Properties-Count | 1     |
      | Min-Per-Night-Cents    | 10000 |
      | Max-Per-Night-Cents    | 11000 |
      | Avg-Per-Night-Cents    | 10500 |

  Scenario: Filtering by price range
    Given the following properties:
      | ID | Account          | Name          | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 11 | Diego Apartments | Too Expensive | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 1           |
    And the following units:
      | ID | Property         | Name          | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 10 | Awesome Property | Double Apt    | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 11 | Awesome Property | Standard Apt  | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Awesome Property | Kitchen       | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 13 | Too Expensive    | Too Expensive | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit          | MinNights | Cents |
      | Kitchen       | 1         | 5000  |
      | Standard Apt  | 1         | 11000 |
      | Standard Apt  | 7         | 70000 |
      | Double Apt    | 1         | 12000 |
      | Double Apt    | 7         | 70000 |
      | Too Expensive | 1         | 50000 |
    When I search for units with the following filters:
      | cityID        | 2     |
      | minPriceCents | 9000  |
      | maxPriceCents | 11000 |
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name         | PerNightCents | TotalCents |
      | Awesome Property | Standard Apt | 11000         | 11000      |
    And I should receive the following headers:
      | Total-Properties-Count | 2     |
      | Min-Per-Night-Cents    | 5000  |
      | Max-Per-Night-Cents    | 50000 |
      | Avg-Per-Night-Cents    | 19500 |

  Scenario: Filtering by price range with dates, considering longer periods
    Given the following properties:
      | ID | Account          | Name          | PublishedAt          | Overview                 | Latitude | Longitude | Address1                | CityID | CountryID | MinimumStay |
      | 11 | Diego Apartments | Too Expensive | 2018-06-01T08:00:00Z | <p>Property Overview</p> | 100      | 200       | 88 Tai Tam Reservoir Rd | 2      | 1         | 1           |
    And the following units:
      | ID | Property         | Name          | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 10 | Awesome Property | Double Apt    | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 11 | Awesome Property | Standard Apt  | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 12 | Awesome Property | Kitchen       | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
      | 13 | Too Expensive    | Too Expensive | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit          | MinNights | Cents |
      | Kitchen       | 1         | 5000  |
      | Standard Apt  | 1         | 11000 |
      | Standard Apt  | 7         | 70000 |
      | Double Apt    | 1         | 12000 |
      | Double Apt    | 7         | 69000 |
      | Too Expensive | 1         | 50000 |
    When I search for units with the following filters:
      | cityID        | 2          |
      | checkIn       | 2019-01-01 |
      | checkOut      | 2019-01-08 |
      | minPriceCents | 9000       |
      | maxPriceCents | 11000      |
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name         | PerNightCents | TotalCents |
      | Awesome Property | Double Apt   | 9857          | 69000      |
      | Awesome Property | Standard Apt | 10000         | 70000      |
    And I should receive the following headers:
      | Total-Properties-Count | 2     |
      | Min-Per-Night-Cents    | 5000  |
      | Max-Per-Night-Cents    | 50000 |
      | Avg-Per-Night-Cents    | 18714 |

  Scenario: No listings match
    When I search for units with the following filters:
      | cityID | 2 |
    Then I should receive an "OK" response with an empty list
    And I should receive the following headers:
      | Total-Properties-Count | 0 |
      | Min-Per-Night-Cents    | 0 |
      | Max-Per-Night-Cents    | 0 |
      | Avg-Per-Night-Cents    | 0 |

  Scenario: Not duplicating units when multiple prices match
    Given the following units:
      | ID | Property         | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | 10 | Awesome Property | Double Apt   | 2        | 2         | 62     | 6            | 10    | 2018-06-01T08:00:00Z |
      | 11 | Nice Property    | Standard Apt | 1        | 1         | 52     | 3            | 15    | 2018-06-01T08:00:00Z |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Double Apt   | 1         | 12000 |
      | Double Apt   | 2         | 22000 |
      | Standard Apt | 1         | 11500 |
      | Standard Apt | 2         | 21000 |
    When I search for units with the following filters:
      | cityID   | 2          |
      | checkIn  | 2019-01-01 |
      | checkOut | 2019-01-08 |
    Then I should receive an "OK" response with the following search results:
      | PropertyName     | Name         |
      | Awesome Property | Double Apt   |
      | Nice Property    | Standard Apt |
