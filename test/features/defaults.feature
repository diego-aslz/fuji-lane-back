Feature: Application Defaults

  Scenario: Loading defaults
    When defaults are loaded
    Then I should have the following countries:
      | ID | Name      |
      | 1  | China     |
      | 2  | Hong Kong |
      | 3  | Japan     |
      | 4  | Singapore |
      | 5  | Vietnam   |
    And I should have the following cities:
      | ID  | CountryID | Name        | Slug        |
      | 101 | 1         | Beijing     | beijing     |
      | 102 | 1         | Chengdu     | chengdu     |
      | 103 | 1         | Chongqing   | chongqing   |
      | 104 | 1         | Dongguan    | dongguan    |
      | 105 | 1         | Guangzhou   | guangzhou   |
      | 106 | 1         | Shanghai    | shanghai    |
      | 107 | 1         | Shenyang    | shenyang    |
      | 108 | 1         | Shenzhen    | shenzhen    |
      | 109 | 1         | Tianjin     | tianjin     |
      | 110 | 1         | Wuhan       | wuhan       |
      | 201 | 2         | Hong Kong   | hong-kong   |
      | 301 | 3         | Fukuoka     | fukuoka     |
      | 302 | 3         | Kawasaki    | kawasaki    |
      | 303 | 3         | Kobe        | kobe        |
      | 304 | 3         | Kyoto       | kyoto       |
      | 305 | 3         | Nagoya      | nagoya      |
      | 306 | 3         | Osaka       | osaka       |
      | 307 | 3         | Saitama     | saitama     |
      | 308 | 3         | Sapporo     | sapporo     |
      | 309 | 3         | Tokyo       | tokyo       |
      | 310 | 3         | Yokohama    | yokohama    |
      | 401 | 4         | Singapore   | singapore   |
      | 501 | 5         | Ho Chi Minh | ho-chi-minh |
      | 502 | 5         | Hanoi       | hanoi       |
