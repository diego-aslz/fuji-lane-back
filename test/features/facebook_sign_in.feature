Feature: Facebook Signing In

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
    Then I should have the following users:
      | Email              | Name           | FacebookID | LastSignedIn         |
      | diego@selzlein.com | Diego Selzlein | 123        | 2018-06-01T08:00:00Z |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |

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
    Then I should have the following users:
      | Email              | Name           | FacebookID | LastSignedIn         |
      | diego@selzlein.com | Diego Selzlein | 123        | 2018-06-01T08:00:00Z |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |

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
    Then I should have the following users:
      | Email              | Name           | FacebookID | LastSignedIn         |
      | diego@selzlein.com | Diego Selzlein | 123        | 2018-06-01T08:00:00Z |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |

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
    Then I should have the following users:
      | Email                    | Name           | FacebookID | LastSignedIn         |
      | diego+other@selzlein.com | Diego Selzlein | 123        | 2018-06-01T08:00:00Z |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z     |
      | renewAfter | 2018-06-05T08:00:00Z     |
      | expiresAt  | 2018-06-08T08:00:00Z     |

  Scenario: Signing in with unrecognized token
    When the following user signs in via Facebook:
      | accessToken | unrecognized-token |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then I should have no users
    And I should receive an "UNAUTHORIZED" response with the following errors:
      | You could not be authenticated |

  Scenario: Signing in with invalid token
    Given Facebook recognizes the following tokens:
      | accessToken | AppID | IsValid | UserID |
      | token-123   | 111   | false   | 123    |
    When the following user signs in via Facebook:
      | accessToken | token-123          |
      | id          | 123                |
      | name        | Diego Selzlein     |
      | email       | diego@selzlein.com |
    Then I should have no users
    And I should receive an "UNAUTHORIZED" response with the following errors:
      | You could not be authenticated |
