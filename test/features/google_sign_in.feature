Feature: Facebook Signing In

  Background:
    Given it is currently "01 Jun 18 08:00"

  Scenario: Signing in for the first time
    When the following user signs in via Google:
      | email    | diego@selzlein.com                                                                                 |
      | name     | Diego Selzlein                                                                                     |
      | googleID | 123                                                                                                |
      | picture  | https://lh6.googleusercontent.com/-HJZOHb_H5pA/AAAAAAAAAAI/AAAAAAAAAFI/sSP5vw-siwc/s96-c/photo.jpg |
    Then I should have the following users:
      | Email              | Name           | GoogleID | LastSignedIn         | PictureURL                                                                                         |
      | diego@selzlein.com | Diego Selzlein | 123      | 2018-06-01T08:00:00Z | https://lh6.googleusercontent.com/-HJZOHb_H5pA/AAAAAAAAAAI/AAAAAAAAAFI/sSP5vw-siwc/s96-c/photo.jpg |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |

  Scenario: Signing in for a second time updating attributes
    Given the following users:
      | Email              | Name                 | GoogleID |
      | diego@selzlein.com | Diego Aguir Selzlein | 123      |
    When the following user signs in via Google:
      | email    | diego@selzlein.com                                                                                 |
      | name     | Diego Selzlein                                                                                     |
      | googleID | 123                                                                                                |
      | picture  | https://lh6.googleusercontent.com/-HJZOHb_H5pA/AAAAAAAAAAI/AAAAAAAAAFI/sSP5vw-siwc/s96-c/photo.jpg |
    Then I should have the following users:
      | Email              | Name           | GoogleID | LastSignedIn         | PictureURL                                                                                         |
      | diego@selzlein.com | Diego Selzlein | 123      | 2018-06-01T08:00:00Z | https://lh6.googleusercontent.com/-HJZOHb_H5pA/AAAAAAAAAAI/AAAAAAAAAFI/sSP5vw-siwc/s96-c/photo.jpg |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |

  Scenario: Signing in when user already exists, but didn't come from Google
    Given the following users:
      | Email              | Name                 |
      | diego@selzlein.com | Diego Aguir Selzlein |
    When the following user signs in via Google:
      | email    | diego@selzlein.com                                                                                 |
      | name     | Diego Selzlein                                                                                     |
      | googleID | 123                                                                                                |
      | picture  | https://lh6.googleusercontent.com/-HJZOHb_H5pA/AAAAAAAAAAI/AAAAAAAAAFI/sSP5vw-siwc/s96-c/photo.jpg |
    Then I should have the following users:
      | Email              | Name           | GoogleID | LastSignedIn         |
      | diego@selzlein.com | Diego Selzlein | 123      | 2018-06-01T08:00:00Z |
    And I should receive an "OK" response with the following body:
      | issuedAt   | 2018-06-01T08:00:00Z |
      | renewAfter | 2018-06-05T08:00:00Z |
      | expiresAt  | 2018-06-08T08:00:00Z |
