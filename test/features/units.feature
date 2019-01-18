Feature: Units Management

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following properties:
      | ID | Account          | Name          |
      | 1  | Diego Apartments | ACME Downtown |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Adding an unit to my property
    When I add the following unit:
      | Property     | ACME Downtown               |
      | Name         | Standard Apt                |
      | Overview     | <strong>Big rooms!</strong> |
      | Bedrooms     | 1                           |
      | SizeM2       | 52                          |
      | MaxOccupancy | 3                           |
      | Count        | 15                          |
    Then the system should respond with "CREATED"
    And I should have the following units:
      | Property      | Name         | Slug         | Bedrooms | SizeM2 | MaxOccupancy | Count | Overview                    |
      | ACME Downtown | Standard Apt | standard-apt | 1        | 52     | 3            | 15    | <strong>Big rooms!</strong> |

  Scenario: Adding an unit to my property which would duplicate slugs
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    When I add the following unit:
      | Property     | ACME Downtown |
      | Name         | Standard  Apt |
      | Bedrooms     | 1             |
      | SizeM2       | 52            |
      | MaxOccupancy | 3             |
      | Count        | 15            |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Name is already in use |

  Scenario: Adding an invalid unit
    When I add the following unit:
      | MaxOccupancy | 3                                            |
      | Overview     | <strong>Big rooms!<script></script></strong> |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | property is required                  |
      | name is required                      |
      | bedrooms is required                  |
      | size is required                      |
      | number of unit type is required       |
      | overview: script tags are not allowed |
    And I should have no units

  Scenario: Adding an unit to a property that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | Name        |
      | 2  | Other   | ACME Uptown |
    When I add the following unit:
      | Property     | ACME Uptown  |
      | Name         | Standard Apt |
      | Bedrooms     | 1            |
      | SizeM2       | 52           |
      | MaxOccupancy | 3            |
      | Count        | 15           |
    Then the system should respond with "NOT FOUND"
    And I should have no units

  Scenario: Updating an unit
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Unit         | Name          | Uploaded |
      | 3  | Standard Apt | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name                   | Std Apartment                 |
      | Bedrooms               | 2                             |
      | Bathrooms              | 3                             |
      | SizeM2                 | 50                            |
      | MaxOccupancy           | 2                             |
      | Count                  | 20                            |
      | BasePriceCents         | 12000                         |
      | OneNightPriceCents     | 11000                         |
      | OneWeekPriceCents      | 40000                         |
      | ThreeMonthsPriceCents  | 350000                        |
      | SixMonthsPriceCents    | 650000                        |
      | TwelveMonthsPriceCents | 1200000                       |
      | FloorPlanImageID       | 3                             |
      | Overview               | <strong>Big windows!</strong> |
    Then the system should respond with "OK"
    And I should have the following units:
      | Property      | Name          | Slug          | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents | OneNightPriceCents | OneWeekPriceCents | ThreeMonthsPriceCents | SixMonthsPriceCents | TwelveMonthsPriceCents | FloorPlanImageID | Overview                      |
      | ACME Downtown | Std Apartment | std-apartment | 2        | 3         | 50     | 2            | 20    | 12000          | 11000              | 40000             | 350000                | 650000              | 1200000                | 3                | <strong>Big windows!</strong> |

  Scenario: Updating an unit with a name which would duplicate slugs
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 10    |
      | ACME Downtown | Double Apt   | 1        | 52     | 3            | 15    |
    When I update unit "Double Apt" with the following attributes:
      | Name | Standard  Apt |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Name is already in use |

  Scenario: Updating an unit with invalid Overview
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    When I update unit "Standard Apt" with the following attributes:
      | Overview | <strong>Big windows!</strong><script></script> |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | overview: script tags are not allowed |

  Scenario: Updating an unit that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | Name        |
      | 2  | Other   | ACME Uptown |
    And the following units:
      | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |
    When I update unit "Standard Apt" with the following attributes:
      | Name         | Std Apartment |
      | Bedrooms     | 2             |
      | SizeM2       | 50            |
      | MaxOccupancy | 2             |
      | Count        | 20            |
    Then the system should respond with "NOT FOUND"
    And I should have the following units:
      | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |

  Scenario: Updating an unit with a floor plan image that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | Name            |
      | 2  | Other   | ACME Apartments |
    And the following units:
      | Property        | Name               | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown   | Standard Apt       | 1        | 52     | 3            | 15    |
      | ACME Apartments | Standard Other Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Unit               | Name          | Uploaded |
      | 3  | Standard Other Apt | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name             | Std Apartment |
      | Bedrooms         | 2             |
      | SizeM2           | 50            |
      | MaxOccupancy     | 2             |
      | Count            | 20            |
      | FloorPlanImageID | 3             |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | floor plan image does not exist |
    And I should have the following units:
      | Property        | Name               | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown   | Standard Apt       | 1        | 52     | 3            | 15    |
      | ACME Apartments | Standard Other Apt | 1        | 52     | 3            | 15    |

  Scenario: Updating unit amenities
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
      | Standard Apt | desk      |
    And the following amenities:
      | Unit         | Type   | Name   |
      | Standard Apt | custom | Towels |
      | Standard Apt | custom | Soap   |
    When I update unit "Standard Apt" with the following amenities:
      | Type             | Name             |
      | desk             | Desk             |
      | air_conditioning | Air Conditioning |
      | custom           | Soap             |
      | custom           | Windows          |
    Then the system should respond with "OK"
    And I should have the following amenities:
      | Unit         | Type             | Name    |
      | Standard Apt | desk             |         |
      | Standard Apt | custom           | Soap    |
      | Standard Apt | air_conditioning |         |
      | Standard Apt | custom           | Windows |

  Scenario: Getting unit details
    Given the following units:
      | ID | Property      | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents | OneNightPriceCents | OneWeekPriceCents | ThreeMonthsPriceCents | SixMonthsPriceCents | TwelveMonthsPriceCents | Overview                    |
      | 2  | ACME Downtown | Standard Apt | 1        | 2         | 52     | 3            | 15    | 12000          | 11000              | 40000             | 350000                | 650000              | 1200000                | <strong>Good view!</strong> |
    And the following images:
      | ID | Unit         | Name          | Uploaded | URL                                 | Type       | Size | Position |
      | 3  | Standard Apt | blueprint.jpg | true     | https://s3.amazonaws.com/blue.jpg   | image/jpeg | 5000 | 0        |
      | 4  | Standard Apt | front.jpg     | true     | https://s3.amazonaws.com/front.jpg  | image/jpeg | 5000 | 2        |
      | 5  | Standard Apt | back.jpg      | true     | https://s3.amazonaws.com/back.jpg   | image/jpeg | 5000 | 1        |
      | 6  | Standard Apt | failed.jpg    | false    | https://s3.amazonaws.com/failed.jpg | image/jpeg | 5000 | 3        |
    And unit "Standard Apt" has:
      | FloorPlanImageID | 3 |
    And the following amenities:
      | ID | Unit         | Type      |
      | 1  | Standard Apt | bathrobes |
    When I get details for unit "Standard Apt"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "id": 2,
        "publishedAt": null,
        "propertyID": 1,
        "name": "Standard Apt",
        "slug": "standard-apt",
        "bedrooms": 1,
        "bathrooms": 2,
        "sizeM2": 52,
        "maxOccupancy": 3,
        "count": 15,
        "basePriceCents": 12000,
        "oneNightPriceCents": 11000,
        "oneWeekPriceCents": 40000,
        "threeMonthsPriceCents": 350000,
        "sixMonthsPriceCents": 650000,
        "twelveMonthsPriceCents": 1200000,
        "overview": "<strong>Good view!</strong>",
        "floorPlanImage": {
          "id": 3,
          "name": "blueprint.jpg",
          "type": "image/jpeg",
          "size": 5000,
          "url": "https://s3.amazonaws.com/blue.jpg",
          "uploaded": true,
          "position": 0
        },
        "images": [
          {
            "id": 5,
            "name": "back.jpg",
            "type": "image/jpeg",
            "size": 5000,
            "url": "https://s3.amazonaws.com/back.jpg",
            "uploaded": true,
            "position": 1
          }, {
            "id": 4,
            "name": "front.jpg",
            "type": "image/jpeg",
            "size": 5000,
            "url": "https://s3.amazonaws.com/front.jpg",
            "uploaded": true,
            "position": 2
          }
        ],
        "amenities": [
          {
            "id": 1,
            "type": "bathrobes",
            "name": "Bathrobes"
          }
        ]
      }
      """

  Scenario: Getting unit details for an unit the user does not have access to
    Given the following accounts:
      | Name            |
      | John Apartments |
    And the following properties:
      | ID | Account         | Name        |
      | 2  | John Apartments | ACME Uptown |
    And the following units:
      | ID | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | 2  | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |
    When I get details for unit "Standard Apt"
    Then the system should respond with "NOT FOUND"

  Scenario: Publishing my unit
    Given the following units:
      | ID | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | BasePriceCents |
      | 2  | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    | 10000          |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
    And it is currently "05 Jun 18 08:00"
    When I publish unit "2"
    Then the system should respond with "OK"
    And I should have the following units:
      | ID | Name         | PublishedAt          |
      | 2  | Standard Apt | 2018-06-05T08:00:00Z |

  Scenario: Publishing an unit with missing information
    Given the following units:
      | ID | Property      | Name | Bedrooms | SizeM2 | Count |
      | 2  | ACME Downtown |      | 0        | 0      | 0     |
    And the following images:
      | ID | PropertyID | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | 1          | false    | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
    When I publish unit "2"
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Name is required                 |
      | Bedrooms is required             |
      | Size is required                 |
      | Number of Unit Type is required  |
      | At least one amenity is required |
      | At least one image is required   |
    And I should have the following units:
      | ID | Name | PublishedAt |
      | 2  |      |             |
