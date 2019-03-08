Feature: User Profile

  Scenario: Getting my own profile details
    Given the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I get my profile details
    Then I should receive an "OK" response with the following body:
      | name                | Diego Aguir Selzlein |
      | email               | diego@selzlein.com   |
      | unreadBookingsCount | 15                   |

  Scenario: Updating my Profile
    Given the following users:
      | Email              | Password | UnreadBookingsCount |
      | diego@selzlein.com | 12345678 | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Name     | Diego Selzlein        |
      | Email    | diego+fl@selzlein.com |
      | Password | 12345678              |
    Then I should receive an "OK" response
    And I should have the following users:
      | Email                 | Name           | UnreadBookingsCount |
      | diego+fl@selzlein.com | Diego Selzlein | 15                  |

  Scenario: Updating my Profile with an invalid password
    Given the following users:
      | Email              | Password | UnreadBookingsCount |
      | diego@selzlein.com | 12345678 | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Name     | Diego Selzlein        |
      | Email    | diego+fl@selzlein.com |
      | Password | 123456789             |
    Then I should receive an "UNAUTHORIZED" response with the following errors:
      | Password is incorrect |
    And I should have the following users:
      | Email              | Name | UnreadBookingsCount |
      | diego@selzlein.com |      | 15                  |

  Scenario: Updating my Profile with invalid attributes
    Given the following users:
      | Email              | Password | UnreadBookingsCount |
      | diego@selzlein.com | 12345678 | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | Email    | diego+fl |
      | Password | 12345678 |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | Invalid email: diego+fl |
    And I should have the following users:
      | Email              | Name | UnreadBookingsCount |
      | diego@selzlein.com |      | 15                  |

  Scenario: Marking my bookings as read
    Given the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I update my user details with:
      | ResetUnreadBookingsCount | true |
    Then I should have the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 0                   |
    And I should receive an "OK" response
