Feature: Email Sign Up

  Background:
    Given it is currently "01 Jun 18 08:00"

  Scenario: Signing Up using Email
    When the following user signs up with his email:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    Then the system should respond with "CREATED" and the following body:
      | email       | diego@selzlein.com                                                                                                                                                                                                    |
      | issued_at   | 2018-06-01T08:00:00Z                                                                                                                                                                                                  |
      | renew_after | 2018-06-05T08:00:00Z                                                                                                                                                                                                  |
      | expires_at  | 2018-06-08T08:00:00Z                                                                                                                                                                                                  |
      | token       | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImRpZWdvQHNlbHpsZWluLmNvbSIsIkV4cGlyZXNBdCI6MTUyODQ0NDgwMCwiSXNzdWVkQXQiOjE1Mjc4NDAwMDAsIlJlbmV3QWZ0ZXIiOjE1MjgxODU2MDB9.k1dEBzwNMxYFsaBjMzkJFHctUk6Y-txk_GfrR6NX1Vk |
    And we should have the following users:
      | Email              | LastSignedIn         |
      | diego@selzlein.com | 2018-06-01T08:00:00Z |

  Scenario: Signing Up with invalid information
    When the following user signs up with his email:
      | email    | diego   |
      | password | 1234567 |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid email: diego                |
      | Invalid password: minimum size is 8 |
    And we should have no users

  Scenario: Signing Up with existing email
    Given the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |
    When the following user signs up with his email:
      | email    | diego@selzlein.com |
      | password | 12345678           |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid email: already in use |
    And we should have the following users:
      | Email              | Name                 | LastSignedIn         |
      | diego@selzlein.com | Diego Aguir Selzlein | 2018-05-01T08:00:00Z |
