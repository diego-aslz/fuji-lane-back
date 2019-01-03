Feature: User Profile

  Scenario: Updating my Profile
    Given the following users:
      | Email              | Password |
      | diego@selzlein.com | 12345678 |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Name     | Diego Selzlein        |
      | Email    | diego+fl@selzlein.com |
      | Password | 12345678              |
    Then the system should respond with "OK"
    And I should have the following users:
      | Email                 | Name           |
      | diego+fl@selzlein.com | Diego Selzlein |

  Scenario: Updating my Profile with an invalid password
    Given the following users:
      | Email              | Password |
      | diego@selzlein.com | 12345678 |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Name     | Diego Selzlein        |
      | Email    | diego+fl@selzlein.com |
      | Password | 123456789             |
    Then the system should respond with "UNAUTHORIZED" and the following errors:
      | Password does not match |
    And I should have the following users:
      | Email              | Name |
      | diego@selzlein.com |      |

  Scenario: Updating my Profile with invalid attributes
    Given the following users:
      | Email              | Password |
      | diego@selzlein.com | 12345678 |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Email    | diego+fl |
      | Password | 12345678 |
    Then the system should respond with "UNPROCESSABLE ENTITY" and the following errors:
      | Invalid email: diego+fl |
    And I should have the following users:
      | Email              | Name |
      | diego@selzlein.com |      |
