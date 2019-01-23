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
      | ID  | CountryID | Name        | Slug        | Latitude   | Longitude   |
      | 101 | 1         | Beijing     | beijing     | 39.9390731 | 116.11728   |
      | 102 | 1         | Chengdu     | chengdu     | 30.6587488 | 103.935463  |
      | 103 | 1         | Chongqing   | chongqing   | 29.5551377 | 106.4084703 |
      | 104 | 1         | Dongguan    | dongguan    | 22.9764535 | 113.654243  |
      | 105 | 1         | Guangzhou   | guangzhou   | 23.1259819 | 112.9476602 |
      | 106 | 1         | Shanghai    | shanghai    | 31.2246325 | 121.1965709 |
      | 107 | 1         | Shenyang    | shenyang    | 41.8055019 | 123.2964156 |
      | 108 | 1         | Shenzhen    | shenzhen    | 22.5554167 | 113.913795  |
      | 109 | 1         | Tianjin     | tianjin     | 39.1252291 | 117.015353  |
      | 110 | 1         | Wuhan       | wuhan       | 30.5683366 | 114.1603012 |
      | 201 | 2         | Hong Kong   | hong-kong   | 22.284736  | 114.1414606 |
      | 301 | 3         | Fukuoka     | fukuoka     | 33.625038  | 130.0258401 |
      | 302 | 3         | Kawasaki    | kawasaki    | 35.5562073 | 139.5723855 |
      | 303 | 3         | Kobe        | kobe        | 34.6943656 | 135.1556806 |
      | 304 | 3         | Kyoto       | kyoto       | 35.0061653 | 135.7259306 |
      | 305 | 3         | Nagoya      | nagoya      | 35.1680838 | 136.8940904 |
      | 306 | 3         | Osaka       | osaka       | 34.69374   | 135.50218   |
      | 307 | 3         | Saitama     | saitama     | 35.915717  | 139.5787164 |
      | 308 | 3         | Sapporo     | sapporo     | 43.0595074 | 141.3354807 |
      | 309 | 3         | Tokyo       | tokyo       | 35.6735408 | 139.570305  |
      | 310 | 3         | Yokohama    | yokohama    | 35.4620149 | 139.5842306 |
      | 401 | 4         | Singapore   | singapore   | 1.3439166  | 103.7540049 |
      | 501 | 5         | Ho Chi Minh | ho-chi-minh | 10.7659164 | 106.4034602 |
      | 502 | 5         | Hanoi       | hanoi       | 20.9740874 | 105.3724915 |
