Feature: Units Management

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And the following properties:
      | ID | Account          | Name          | Country | City  |
      | 1  | Diego Apartments | ACME Downtown | Japan   | Osaka |
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
    Then I should receive a "CREATED" response
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
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is already in use |

  Scenario: Adding an invalid unit
    When I add the following unit:
      | MaxOccupancy | 3                                            |
      | Overview     | <strong>Big rooms!<script></script></strong> |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
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
      | ID | Account | Name        | Country | City  |
      | 2  | Other   | ACME Uptown | Japan   | Osaka |
    When I add the following unit:
      | Property     | ACME Uptown  |
      | Name         | Standard Apt |
      | Bedrooms     | 1            |
      | SizeM2       | 52           |
      | MaxOccupancy | 3            |
      | Count        | 15           |
    Then I should receive a "NOT FOUND" response
    And I should have no units

  Scenario: Updating an unit
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    And the following images:
      | ID | Unit         | Name          | Uploaded |
      | 3  | Standard Apt | blueprint.jpg | true     |
    When I update unit "Standard Apt" with the following attributes:
      | Name             | Std Apartment                                                                   |
      | Bedrooms         | 2                                                                               |
      | Bathrooms        | 3                                                                               |
      | SizeM2           | 50                                                                              |
      | MaxOccupancy     | 2                                                                               |
      | Count            | 20                                                                              |
      | Prices           | 1: 11000, 2: 12000, 7: 40000, 30: 160000, 90: 350000, 180: 650000, 365: 1200000 |
      | FloorPlanImageID | 3                                                                               |
      | Overview         | <strong>Big windows!</strong>                                                   |
    Then I should receive an "OK" response
    And I should have the following units:
      | Property      | Name          | Slug          | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | FloorPlanImageID | Overview                      |
      | ACME Downtown | Std Apartment | std-apartment | 2        | 3         | 50     | 2            | 20    | 3                | <strong>Big windows!</strong> |
    And I should have the following prices:
      | Unit          | MinNights | Cents   |
      | Std Apartment | 1         | 11000   |
      | Std Apartment | 2         | 12000   |
      | Std Apartment | 7         | 40000   |
      | Std Apartment | 30        | 160000  |
      | Std Apartment | 90        | 350000  |
      | Std Apartment | 180       | 650000  |
      | Std Apartment | 365       | 1200000 |

  Scenario: Updating unit prices
    Given the following units:
      | Property      | Name         |
      | ACME Downtown | Standard Apt |
    And the following prices:
      | Unit         | MinNights | Cents  |
      | Standard Apt | 1         | 11000  |
      | Standard Apt | 2         | 12000  |
      | Standard Apt | 30        | 160000 |
      | Standard Apt | 90        | 300000 |
      | Standard Apt | 180       | 600000 |
    When I update unit "Standard Apt" with the following attributes:
      | Prices | 1: 9000, 7: 40000, 30: 155000, 180: 0 |
    Then I should receive an "OK" response
    And I should have the following prices:
      | Unit         | MinNights | Cents  |
      | Standard Apt | 1         | 9000   |
      | Standard Apt | 30        | 155000 |
      | Standard Apt | 7         | 40000  |

  Scenario: Updating an unit with a name which would duplicate slugs
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 10    |
      | ACME Downtown | Double Apt   | 1        | 52     | 3            | 15    |
    When I update unit "Double Apt" with the following attributes:
      | Name | Standard  Apt |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is already in use |

  Scenario: Updating an unit with invalid Overview
    Given the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    |
    When I update unit "Standard Apt" with the following attributes:
      | Overview | <strong>Big windows!</strong><script></script> |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | overview: script tags are not allowed |

  Scenario: Updating an unit that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | Name        | Country | City  |
      | 2  | Other   | ACME Uptown | Japan   | Osaka |
    And the following units:
      | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |
    When I update unit "Standard Apt" with the following attributes:
      | Name         | Std Apartment |
      | Bedrooms     | 2             |
      | SizeM2       | 50            |
      | MaxOccupancy | 2             |
      | Count        | 20            |
    Then I should receive a "NOT FOUND" response
    And I should have the following units:
      | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |

  Scenario: Updating an unit with a floor plan image that does not belong to me
    Given the following accounts:
      | Name  |
      | Other |
    And the following properties:
      | ID | Account | Name            | Country | City  |
      | 2  | Other   | ACME Apartments | Japan   | Osaka |
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
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
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
      | air-conditioning | Air Conditioning |
      | custom           | Soap             |
      | custom           | Windows          |
    Then I should receive an "OK" response
    And I should have the following amenities:
      | Unit         | Type             | Name    |
      | Standard Apt | desk             |         |
      | Standard Apt | custom           | Soap    |
      | Standard Apt | air-conditioning |         |
      | Standard Apt | custom           | Windows |

  Scenario: Getting unit details
    Given the following units:
      | ID | Property      | Name         | Bedrooms | Bathrooms | SizeM2 | MaxOccupancy | Count | Overview                    |
      | 2  | ACME Downtown | Standard Apt | 1        | 2         | 52     | 3            | 15    | <strong>Good view!</strong> |
    And the following prices:
      | Unit         | MinNights | Cents   |
      | Standard Apt | 1         | 11000   |
      | Standard Apt | 2         | 12000   |
      | Standard Apt | 7         | 40000   |
      | Standard Apt | 90        | 350000  |
      | Standard Apt | 180       | 650000  |
      | Standard Apt | 365       | 1200000 |
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
    Then I should receive an "OK" response with the following JSON:
      """
      {
        "id": 2,
        "publishedAt": null,
        "everPublished": false,
        "propertyID": 1,
        "name": "Standard Apt",
        "slug": "standard-apt",
        "bedrooms": 1,
        "bathrooms": 2,
        "sizeM2": 52,
        "maxOccupancy": 3,
        "count": 15,
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
        ],
        "prices": [{
          "minNights": 1,
          "cents": 11000
        }, {
          "minNights": 2,
          "cents": 12000
        }, {
          "minNights": 7,
          "cents": 40000
        }, {
          "minNights": 90,
          "cents": 350000
        }, {
          "minNights": 180,
          "cents": 650000
        }, {
          "minNights": 365,
          "cents": 1200000
        }]
      }
      """

  Scenario: Getting unit details for an unit the user does not have access to
    Given the following accounts:
      | Name            |
      | John Apartments |
    And the following properties:
      | ID | Account         | Name        | Country | City  |
      | 2  | John Apartments | ACME Uptown | Japan   | Osaka |
    And the following units:
      | ID | Property    | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | 2  | ACME Uptown | Standard Apt | 1        | 52     | 3            | 15    |
    When I get details for unit "Standard Apt"
    Then I should receive a "NOT FOUND" response

  Scenario: Publishing my unit
    Given the following units:
      | ID | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | EverPublished |
      | 2  | ACME Downtown | Standard Apt | 1        | 52     | 3            | 15    | false         |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Standard Apt | 1         | 11000 |
    And the following images:
      | ID | Unit         | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | Standard Apt | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 1        |
    And the following amenities:
      | Unit         | Type      |
      | Standard Apt | bathrobes |
    And it is currently "05 Jun 18 08:00"
    When I publish unit "2"
    Then I should receive an "OK" response
    And I should have the following units:
      | ID | Name         | PublishedAt          | EverPublished |
      | 2  | Standard Apt | 2018-06-05T08:00:00Z | true          |

  Scenario: Publishing an unit with missing information
    Given the following units:
      | ID | Property      | Name | Bedrooms | SizeM2 | Count |
      | 2  | ACME Downtown |      | 0        | 0      | 0     |
    And the following images:
      | ID | PropertyID | Uploaded | Name      | URL                                | Type       | Size    | Position |
      | 1  | 1          | false    | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 | 2        |
    And the following prices:
      | UnitID | MinNights | Cents |
      | 2      | 7         | 11000 |
    When I publish unit "2"
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Name is required                 |
      | Bedrooms is required             |
      | Size is required                 |
      | Number of Unit Type is required  |
      | At least one amenity is required |
      | At least one image is required   |
      | Please provide a base unit price |
    And I should have the following units:
      | ID | Name | PublishedAt | EverPublished |
      | 2  |      |             | false         |

  Scenario: Unpublishing my unit
    Given the following units:
      | ID | Property      | Name         | PublishedAt          |
      | 2  | ACME Downtown | Standard Apt | 2018-06-05T08:00:00Z |
    When I unpublish unit "2"
    Then I should receive an "OK" response
    And I should have the following units:
      | ID | Property      | Name         | PublishedAt |
      | 2  | ACME Downtown | Standard Apt |             |
