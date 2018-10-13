Feature: Cities

  Scenario: Listing supported cities
    Given the following countries:
      | ID | Name  |
      | 1  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 10 | Japan   | Osaka |
    When I list cities
    Then the system should respond with "OK" and the following cities:
      | ID | CountryID | Name  |
      | 10 | 1         | Osaka |
