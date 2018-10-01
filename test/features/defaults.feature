Feature: Application Defaults

  Scenario: Loading database defaults
    When database defaults are loaded
    Then we should have the following countries:
      | ID | Name      |
      | 1  | China     |
      | 2  | Hong Kong |
      | 3  | Japan     |
      | 4  | Singapore |
      | 5  | Vietnam   |
