Feature: Images Management

  Background:
    Given the following configuration:
      | MAX_IMAGE_SIZE_MB | 20 |
    And the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |

  Scenario Outline: Adding an image
    Given the following properties:
      | ID   | Account          | Name            |
      | <ID> | Diego Apartments | ACME Skyscraper |
    And the following units:
      | ID   | Property        | Name         |
      | <ID> | ACME Skyscraper | Standard Apt |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     | build/ing.jpg |
      | Size     | 15000000      |
      | Type     | image/jpeg    |
      | <Target> | <Name>        |
      | Position | 3             |
    Then I should receive an "OK" response with the following image:
      | Name | building.jpg                                                                                                                                                                                                                                                                                         |
      | URL  | https://fujilane-test.s3.amazonaws.com/public/<Collection>/<ID>/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=CREDENTIAL&X-Amz-Date=DATE&X-Amz-Expires=3600&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost%3Bx-amz-acl&X-Amz-Signature=SIGNATURE |
    And I should have the following images:
      | <Target> | Name         | URL                                                                                                   | Uploaded | Type       | Size     | Position | StorageKey                                                     |
      | <Name>   | building.jpg | https://fujilane-test.s3.amazonaws.com/public/<Collection>/<ID>/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    | image/jpeg | 15000000 | 3        | public/<Collection>/<ID>/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |

    Examples:
      | Target   | Name            | Collection | ID |
      | Property | ACME Skyscraper | properties | 20 |
      | Unit     | Standard Apt    | units      | 25 |

  Scenario: Validating file size
    Given the following properties:
      | ID | Account          | Name            |
      | 20 | Diego Apartments | ACME Skyscraper |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     |                 |
      | Size     | 25000000        |
      | Type     | text/csv        |
      | Property | ACME Skyscraper |
      | Position | 1               |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | name is required                           |
      | Invalid size: maximum is 20971520          |
      | Invalid type: needs to be JPEG, PNG or GIF |

  Scenario: Obtaining a signed URL to upload a property image without having an account
    Given the following accounts:
      | Name                |
      | Somebody Apartments |
    And the following users:
      | Email               | Name  |
      | noaccount@gmail.com | Diego |
    And the following properties:
      | Account             | Name            |
      | Somebody Apartments | ACME Skyscraper |
    And I am authenticated with "noaccount@gmail.com"
    When I request an URL to upload the following image:
      | Name     | building.jpg    |
      | Size     | 15000000        |
      | Type     | image/png       |
      | Property | ACME Skyscraper |
      | Position | 1               |
    Then I should receive a "PRECONDITION REQUIRED" response with the following errors:
      | You need a company account to perform this action |

  Scenario Outline: Adding an image to a target I don't have access to
    Given the following accounts:
      | Name                |
      | Somebody Apartments |
    And the following properties:
      | Account             | Name            |
      | Somebody Apartments | ACME Skyscraper |
    And the following units:
      | Property        | Name         |
      | ACME Skyscraper | Standard Apt |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     | building.jpg |
      | Size     | 15000000     |
      | Type     | image/png    |
      | <Target> | <Name>       |
      | Position | 1            |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Could not find <Target> |

    Examples:
      | Target   | Name            |
      | Property | ACME Skyscraper |
      | Unit     | Standard Apt    |

  Scenario: Marking a unit image as uploaded
    Given the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following units:
      | Property        | Name         |
      | ACME Skyscraper | Standard Apt |
    And the following images:
      | Unit         | Name         | URL                                                                                               | Uploaded |
      | Standard Apt | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    |
    And I am authenticated with "diego@selzlein.com"
    When I mark image "building.jpg" as uploaded
    Then I should receive an "OK" response
    And I should have the following images:
      | Name         | Uploaded |
      | building.jpg | true     |

  Scenario: Marking a unit image as uploaded
    Given the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following images:
      | Property        | Name         | URL                                                                                               | Uploaded |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    |
    And I am authenticated with "diego@selzlein.com"
    When I mark image "building.jpg" as uploaded
    Then I should receive an "OK" response
    And I should have the following images:
      | Name         | Uploaded |
      | building.jpg | true     |

  Scenario: Removing a property image
    Given the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following images:
      | Property        | Name         | URL                                                                                               |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
    And I am authenticated with "diego@selzlein.com"
    When I remove the image "building.jpg"
    Then I should receive an "OK" response
    And I should have no images

  Scenario: Removing an unit image
    Given the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following units:
      | Property        | Name         |
      | ACME Skyscraper | Standard Apt |
    And the following images:
      | Unit         | Name         | URL                                                                                               |
      | Standard Apt | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
    And I am authenticated with "diego@selzlein.com"
    When I remove the image "building.jpg"
    Then I should receive an "OK" response
    And I should have no images

  Scenario: Removing a floor plan image
    Given the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following units:
      | Property        | Name         |
      | ACME Skyscraper | Standard Apt |
      | ACME Skyscraper | Double Apt   |
    And the following images:
      | ID | Unit         | Name          | URL                                                                                               |
      | 3  | Standard Apt | building3.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
      | 4  | Double Apt   | building4.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb |
    And unit "Standard Apt" has:
      | FloorPlanImageID | 3 |
    And unit "Double Apt" has:
      | FloorPlanImageID | 4 |
    And I am authenticated with "diego@selzlein.com"
    When I remove the image "building3.jpg"
    Then I should receive an "OK" response
    And I should have the following images:
      | Name          |
      | building4.jpg |
    And I should have the following units:
      | Property        | Name         | FloorPlanImageID |
      | ACME Skyscraper | Standard Apt |                  |
      | ACME Skyscraper | Double Apt   | 4                |

  Scenario: Removing an image that does not belong to me
    Given the following accounts:
      | Name      |
      | Other Acc |
    And the following properties:
      | Account   | Name            |
      | Other Acc | ACME Skyscraper |
    And the following images:
      | Property        | Name         | URL                                                                                               |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa |
    And I am authenticated with "diego@selzlein.com"
    When I remove the image "building.jpg"
    Then I should receive a "NOT FOUND" response
    And I should have the following images:
      | Name         |
      | building.jpg |

  Scenario: Sorting my images
    Given the following properties:
      | ID | Account          | Name            |
      | 5  | Diego Apartments | ACME Skyscraper |
    And the following images:
      | ID | Position | Property        | Name      | URL                                                                    |
      | 5  | 1        | ACME Skyscraper | front.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaa |
      | 10 | 2        | ACME Skyscraper | back.jpg  | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/bbb |
    And I am authenticated with "diego@selzlein.com"
    When I sort images to be like:
      | back.jpg  |
      | front.jpg |
    Then I should have the following images:
      | ID | Position | Name      |
      | 5  | 2        | front.jpg |
      | 10 | 1        | back.jpg  |

  Scenario: Sorting images that do not belong to me
    Given the following accounts:
      | Name      |
      | Other Acc |
    And the following properties:
      | Account   | Name            |
      | Other Acc | ACME Skyscraper |
    And the following images:
      | ID | Position | Property        | Name      | URL                                                                    |
      | 5  | 1        | ACME Skyscraper | front.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaa |
      | 10 | 2        | ACME Skyscraper | back.jpg  | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/bbb |
    And I am authenticated with "diego@selzlein.com"
    When I sort images to be like:
      | back.jpg  |
      | front.jpg |
    Then I should have the following images:
      | ID | Position | Name      |
      | 5  | 1        | front.jpg |
      | 10 | 2        | back.jpg  |
