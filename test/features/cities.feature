Feature: Cities

  Scenario: Listing supported cities
    Given the following countries:
      | ID | Name  |
      | 1  | Japan |
    And the following cities:
      | ID | Country | Name  | Latitude | Longitude |
      | 10 | Japan   | Osaka | 90       | 100       |
    When I list cities
    Then I should receive an "OK" response with the following cities:
      | ID | CountryID | Name  | Latitude | Longitude |
      | 10 | 1         | Osaka | 90       | 100       |
