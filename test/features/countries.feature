Feature: Countries

  Scenario: Listing supported countries
    Given the following countries:
      | ID | Name      |
      | 1  | Hong Kong |
    When I list countries
    Then the system should respond with "OK" and the following countries:
      | ID | Name      |
      | 1  | Hong Kong |
