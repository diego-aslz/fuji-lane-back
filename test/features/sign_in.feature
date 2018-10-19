Feature: Signing In

  Scenario: Signing in with email and password
    Given the following users:
      | Email              | Name                 | LastSignedIn         | Password |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z | 12345678 |
    And it is currently "01 Jun 18 08:00"
    When I sign in with:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    Then the system should respond with "OK" and the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |
    And I should have the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-06-01T08:00:00Z |

  Scenario: Signing in with invalid password
    Given the following users:
      | Email              | Name                 | LastSignedIn         | Password |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z | 12345678 |
    When I sign in with:
      | email    | diego@selzlein.com |
      | password | 123456789          |
    Then the system should respond with "UNAUTHORIZED" and the following errors:
      | Invalid email or password |
    And I should have the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |

  Scenario: Signing in with invalid email
    Given the following users:
      | Email              | Name                 | LastSignedIn         | Password |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z | 12345678 |
    When I sign in with:
      | email    | dieg@selzlein.com |
      | password | 12345678          |
    Then the system should respond with "UNAUTHORIZED" and the following errors:
      | Invalid email or password |
    And I should have the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |
