Feature: Countries

  Scenario: Listing supported countries
    Given the following countries:
      | ID | Name      |
      | 1  | Hong Kong |
    When I list countries
    Then I should receive an "OK" response with the following countries:
      | ID | Name      |
      | 1  | Hong Kong |
