Feature: Email Sign Up

  Background:
    Given it is currently "01 Jun 18 08:00"

  Scenario: Signing Up using Email
    When the following user signs up with his email:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    Then the system should respond with "CREATED" and the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |
    And I should have the following users:
      | Email              | LastSignedIn         |
      | diego@selzlein.com | 2018-06-01T08:00:00Z |

  Scenario: Signing Up with invalid information
    When the following user signs up with his email:
      | email    | diego   |
      | password | 1234567 |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid email: diego                |
      | Invalid password: minimum size is 8 |
    And I should have no users

  Scenario: Signing Up with existing email
    Given the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |
    When the following user signs up with his email:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid email: already in use |
    And I should have the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |
