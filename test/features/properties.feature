Feature: Properties Management

  Scenario: Adding a new property
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    When I add a new property
    Then the system should respond with "CREATED"
    And we should have the following properties:
      | Account          | State |
      | Diego Apartments | draft |

  Scenario: Adding a new property without having an Account
    Given the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And I am authenticated with "diego@selzlein.com"
    When I add a new property
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | You need a company account |
    And we should have no properties

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
    Then the system should respond with "OK" and the following pre-signed URL:
      | bucket     | fujilane-test                            |
      | key        | public/properties/20/images/building.jpg |
      | expiration | 3600                                     |

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

  Scenario: Getting property details
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And the following properties:
      | Account          | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country |
      | Diego Apartments | Draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
    And I am authenticated with "diego@selzlein.com"
    When I get details for property "ACME Downtown"
    Then the system should respond with "OK" and the following body:
      | state      | draft         |
      | name       | ACME Downtown |
      | address1   | Add. One      |
      | address2   | Add. Two      |
      | address3   | Add. Three    |
      | cityID     | 3             |
      | postalCode | 223344        |
      | countryID  | 2             |

  Scenario: Getting property details for a property the user does not have access to
    Given the following accounts:
      | Name             |
      | Diego Apartments |
      | John Apartments  |
    And the following users:
      | Account          | Email              | Name                 |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And the following properties:
      | Account         | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country |
      | John Apartments | Draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
    And I am authenticated with "diego@selzlein.com"
    When I get details for property "ACME Downtown"
    Then the system should respond with "NOT FOUND"
