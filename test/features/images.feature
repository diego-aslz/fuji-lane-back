Feature: Images Management

  Background:
    Given the following configuration:
      | MAX_IMAGE_SIZE_MB | 20 |

  Scenario: Obtaining a signed URL to upload a property image
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following properties:
      | ID | Account          | Name            |
      | 20 | Diego Apartments | ACME Skyscraper |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     | build/ing.jpg   |
      | Size     | 15000000        |
      | Type     | image/jpeg      |
      | Property | ACME Skyscraper |
    Then the system should respond with "OK" and the following image:
      | Name | building.jpg                                                                                                                                                                                                                                                                                     |
      | URL  | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=CREDENTIAL&X-Amz-Date=DATE&X-Amz-Expires=3600&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost%3Bx-amz-acl&X-Amz-Signature=SIGNATURE |
    And I should have the following images:
      | Property        | Name         | URL                                                                                               | Uploaded | Type       | Size     |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    | image/jpeg | 15000000 |

  Scenario: Validating file size
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following properties:
      | ID | Account          | Name            |
      | 20 | Diego Apartments | ACME Skyscraper |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     |                 |
      | Size     | 25000000        |
      | Type     | text/csv        |
      | Property | ACME Skyscraper |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid name: cannot be blank      |
      | Invalid size: maximum is 20971520  |
      | Invalid type: needs to be an image |

  Scenario: Obtaining a signed URL to upload a property image without having an account
    Given the following accounts:
      | Name                |
      | Somebody Apartments |
    And the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And the following properties:
      | Account             | Name            |
      | Somebody Apartments | ACME Skyscraper |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     | building.jpg    |
      | Size     | 15000000        |
      | Type     | image/png       |
      | Property | ACME Skyscraper |
    Then the system should respond with "PRECONDITION REQUIRED" and the following errors:
      | You do not have an owner account |

  Scenario: Obtaining a signed URL to upload a property image for a property I don't have access to
    Given the following accounts:
      | Name                |
      | Somebody Apartments |
      | Diego Apartments    |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following properties:
      | Account             | Name            |
      | Somebody Apartments | ACME Skyscraper |
    And I am authenticated with "diego@selzlein.com"
    When I request an URL to upload the following image:
      | Name     | building.jpg    |
      | Size     | 15000000        |
      | Type     | image/png       |
      | Property | ACME Skyscraper |
    Then the system should respond with "NOT FOUND"

  Scenario: Marking an image as uploaded
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following properties:
      | Account          | Name            |
      | Diego Apartments | ACME Skyscraper |
    And the following images:
      | Property        | Name         | URL                                                                                               | Uploaded |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    |
    And I am authenticated with "diego@selzlein.com"
    When I mark image "building.jpg" as uploaded
    And I should have the following images:
      | Name         | Uploaded |
      | building.jpg | true     |
