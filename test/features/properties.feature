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
    Then the system should respond with "PRECONDITION REQUIRED" and the following errors:
      | You need a company account to perform this action |
    And we should have no properties

  Scenario: Getting property details
    Given the following accounts:
      | Name             |
      | Diego Apartments |
      | Other            |
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
      | ID | Account          | State | Name          | Address1 | Address2 | Address3   | City  | PostalCode | Country |
      | 1  | Diego Apartments | Draft | ACME Downtown | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
      | 2  | Other            | Draft | Other Prop    | Add. One | Add. Two | Add. Three | Osaka | 223344     | Japan   |
    And the following images:
      | ID | Property      | Uploaded | Name      | URL                                | Type       | Size    |
      | 1  | ACME Downtown | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 |
      | 2  | ACME Downtown | false    | back.jpg  | https://s3.amazonaws.com/back.jpg  | image/jpeg | 1000000 |
      | 3  | Other Prop    | true     | front.jpg | https://s3.amazonaws.com/front.jpg | image/jpeg | 1000000 |
    And I am authenticated with "diego@selzlein.com"
    When I get details for property "ACME Downtown"
    Then the system should respond with "OK" and the following JSON:
      """
      {
        "id": 1,
        "state": "draft",
        "name": "ACME Downtown",
        "address1": "Add. One",
        "address2": "Add. Two",
        "address3": "Add. Three",
        "postalCode": "223344",
        "cityID": 3,
        "postalCode": "223344",
        "countryID":  2,
        "images": [
          {
            "id": 1,
            "name": "front.jpg",
            "type": "image/jpeg",
            "size": 1000000,
            "url": "https://s3.amazonaws.com/front.jpg",
            "uploaded":true,
            "propertyID":1
          }
        ]
      }
      """

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
