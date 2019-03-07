Feature: Users

  Scenario: Getting my own user details
    Given the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I get my user details
    Then I should receive an "OK" response with the following body:
      | name                | Diego Aguir Selzlein |
      | email               | diego@selzlein.com   |
      | unreadBookingsCount | 15                   |

  Scenario: Marking my bookings as read
    Given the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 15                  |
    And I am authenticated with "diego@selzlein.com"
    When I mark my bookings as read
    Then I should have the following users:
      | Email              | Name                 | UnreadBookingsCount |
      | diego@selzlein.com | Diego Aguir Selzlein | 0                   |
    And I should receive an "OK" response with the following body:
      | name                | Diego Aguir Selzlein |
      | email               | diego@selzlein.com   |
      | unreadBookingsCount | 0                    |
