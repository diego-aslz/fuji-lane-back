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
      | ID  | CountryID | Name        |
      | 101 | 1         | Beijing     |
      | 102 | 1         | Chengdu     |
      | 103 | 1         | Chongqing   |
      | 104 | 1         | Dongguan    |
      | 105 | 1         | Guangzhou   |
      | 106 | 1         | Shanghai    |
      | 107 | 1         | Shenyang    |
      | 108 | 1         | Shenzhen    |
      | 109 | 1         | Tianjin     |
      | 110 | 1         | Wuhan       |
      | 201 | 2         | Hong Kong   |
      | 301 | 3         | Fukuoka     |
      | 302 | 3         | Kawasaki    |
      | 303 | 3         | Kobe        |
      | 304 | 3         | Kyoto       |
      | 305 | 3         | Nagoya      |
      | 306 | 3         | Osaka       |
      | 307 | 3         | Saitama     |
      | 308 | 3         | Sapporo     |
      | 309 | 3         | Tokyo       |
      | 310 | 3         | Yokohama    |
      | 401 | 4         | Singapore   |
      | 501 | 5         | Ho Chi Minh |
      | 502 | 5         | Hanoi       |
