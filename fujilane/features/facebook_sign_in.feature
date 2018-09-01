Feature: Facebook Signing in

  Background:
    Given it is currently "01 Jun 18 08:00"

  Scenario: Signing in for the first time
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email              | Name           | FacebookID |
      | diego@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK" and the following body:
      | email       | diego@selzlein.com                                                                                                                                                                                                    |
      | issued_at   | 2018-06-01T08:00:00Z                                                                                                                                                                                                  |
      | renew_after | 2018-06-05T08:00:00Z                                                                                                                                                                                                  |
      | expires_at  | 2018-06-08T08:00:00Z                                                                                                                                                                                                  |
      | token       | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImRpZWdvQHNlbHpsZWluLmNvbSIsIkV4cGlyZXNBdCI6MTUyODQ0NDgwMCwiSXNzdWVkQXQiOjE1Mjc4NDAwMDAsIlJlbmV3QWZ0ZXIiOjE1MjgxODU2MDB9.k1dEBzwNMxYFsaBjMzkJFHctUk6Y-txk_GfrR6NX1Vk |

  Scenario: Signing in for a second time updating attributes
    Given the following users:
      | Email              | Name                 | FacebookID |
      | diego@selzlein.com | Diego Aguir Selzlein | 123        |
    And Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-222   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-222          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email              | Name           | FacebookID |
      | diego@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK" and the following body:
      | email       | diego@selzlein.com   |
      | issued_at   | 2018-06-01T08:00:00Z |
      | renew_after | 2018-06-05T08:00:00Z |
      | expires_at  | 2018-06-08T08:00:00Z |

  Scenario: Signing in when user already exists, but didn't come from Facebook
    Given the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    And Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-222   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-222          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email              | Name           | FacebookID |
      | diego@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK" and the following body:
      | email       | diego@selzlein.com   |
      | issued_at   | 2018-06-01T08:00:00Z |
      | renew_after | 2018-06-05T08:00:00Z |
      | expires_at  | 2018-06-08T08:00:00Z |

  Scenario: Signing in when user already exists, but Facebook email changed
    Given the following users:
      | Email                    | Name                 | FacebookID |
      | diego+other@selzlein.com | Diego Aguir Selzlein | 123        |
    And Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-222   | 111   | true    | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-222          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have the following users:
      | Email                    | Name           | FacebookID |
      | diego+other@selzlein.com | Diego Selzlein | 123        |
    And the system should respond with "OK" and the following body:
      | email       | diego+other@selzlein.com |
      | issued_at   | 2018-06-01T08:00:00Z     |
      | renew_after | 2018-06-05T08:00:00Z     |
      | expires_at  | 2018-06-08T08:00:00Z     |

  Scenario: Signing in with unrecognized token
    When the following user signs in via Facebook:
      | accessToken | unrecognized-token |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have no users
    And the system should respond with "UNAUTHORIZED" and no body

  Scenario: Signing in with invalid token
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | false   | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then we should have no users
    And the system should respond with "UNAUTHORIZED" and no body
