Feature: Images Management

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
    When I request an URL to upload an image called "build/ing.jpg" for property "ACME Skyscraper"
    Then the system should respond with "OK" and the following image:
      | Name | building.jpg                                                                                                                                                                                                                                                                                                                           |
      | URL  | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=11111111111111111111%2F20181008%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=TEST_DATE&X-Amz-Expires=3600&X-Amz-SignedHeaders=host%3Bx-amz-acl&X-Amz-Signature=TEST_SIGNATURE |
    And I should have the following images:
      | Property        | Name         | URL                                                                                               | Uploaded |
      | ACME Skyscraper | building.jpg | https://fujilane-test.s3.amazonaws.com/public/properties/20/images/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | false    |

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
    When I request an URL to upload an image called "building.jpg" for property "ACME Skyscraper"
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
    When I request an URL to upload an image called "building.jpg" for property "ACME Skyscraper"
    Then the system should respond with "NOT FOUND"
