Feature: Application Defaults

  Scenario: Loading defaults
    When defaults are loaded
    Then we should have the following countries:
      | ID | Name      |
      | 1  | China     |
      | 2  | Hong Kong |
      | 3  | Japan     |
      | 4  | Singapore |
      | 5  | Vietnam   |
    And we should have the following cities:
      | ID | CountryID | Name  |
      | 1  | 3         | Osaka |
