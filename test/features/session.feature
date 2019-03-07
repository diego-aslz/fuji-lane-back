Feature: Session

  Scenario: Accessing a protected resource without an authentication token
    When I add a new property
    Then I should receive an "UNAUTHORIZED" response
    And I should have no properties

  Scenario: Accessing a protected resource with an authentication token for an invalid user
    Given it is currently "01 Jun 18 08:00"
    And the following session:
      | IssuedAt   | 2018-06-01T08:00:00Z |
      | RenewAfter | 2018-06-05T08:00:00Z |
      | ExpiresAt  | 2018-06-08T08:00:00Z |
    When I add a new property
    Then I should receive an "UNAUTHORIZED" response with the following errors:
      | You need to sign in |
    And I should have no properties

  Scenario: Accessing a protected resource with an expired session
    Given it is currently "01 Jun 18 08:00"
    And the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And the following session:
      | IssuedAt   | 2018-05-21T08:00:00Z |
      | RenewAfter | 2018-05-25T08:00:00Z |
      | ExpiresAt  | 2018-05-28T08:00:00Z |
    When I add a new property
    Then I should receive an "UNAUTHORIZED" response with the following errors:
      | Your session expired |
    And I should have no properties

  Scenario: Getting session details for authenticated user
    Given the following countries:
      | ID | Name  |
      | 1  | China |
    And the following accounts:
      | Name             | Phone   | Country |
      | Diego Apartments | 555 555 | China   |
    And the following users:
      | Account          | Email              | Name                 | Password | FacebookID |
      | Diego Apartments | diego@selzlein.com | Diego Aguir Selzlein | 12345678 | 12345      |
    And it is currently "01 Jun 18 08:00"
    When I sign in with:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    And I should receive an "OK" response with the following JSON:
      """
      {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImRpZWdvQHNlbHpsZWluLmNvbSIsIkV4cGlyZXNBdCI6MTUyODQ0NDgwMCwiSXNzdWVkQXQiOjE1Mjc4NDAwMDAsIlJlbmV3QWZ0ZXIiOjE1MjgxODU2MDB9.k1dEBzwNMxYFsaBjMzkJFHctUk6Y-txk_GfrR6NX1Vk",
        "issuedAt": "2018-06-01T08:00:00Z",
        "expiresAt": "2018-06-08T08:00:00Z",
        "renewAfter": "2018-06-05T08:00:00Z",
        "user": {
          "picture": "https://graph.facebook.com/12345/picture?width=64&height=64",
          "name": "Diego Aguir Selzlein",
          "email": "diego@selzlein.com",
          "unreadBookingsCount": 0
        },
        "account": {
          "name": "Diego Apartments",
          "phone": "555 555",
          "countryID": 1
        }
      }
      """

  Scenario: Renewing session token
    Given it is currently "05 Jun 18 08:00"
    And the following users:
      | Email              | Name                 | Password | FacebookID |
      | diego@selzlein.com | Diego Aguir Selzlein | 12345678 | 12345      |
    And the following session:
      | Email      | diego@selzlein.com   |
      | IssuedAt   | 2018-06-01T08:00:00Z |
      | RenewAfter | 2018-06-05T08:00:00Z |
      | ExpiresAt  | 2018-06-08T08:00:00Z |
    When I renew the session token
    Then I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-05T08:00:00Z |
      | renewAfter | 2018-06-09T08:00:00Z |
      | expiresAt  | 2018-06-12T08:00:00Z |

  Scenario: Renewing session token before time
    Given it is currently "05 Jun 18 07:00"
    And the following users:
      | Email              | Name                 | Password | FacebookID |
      | diego@selzlein.com | Diego Aguir Selzlein | 12345678 | 12345      |
    And the following session:
      | Email      | diego@selzlein.com   |
      | IssuedAt   | 2018-06-01T08:00:00Z |
      | RenewAfter | 2018-06-05T08:00:00Z |
      | ExpiresAt  | 2018-06-08T08:00:00Z |
    When I renew the session token
    Then I should receive a "NOT MODIFIED" response with no body
