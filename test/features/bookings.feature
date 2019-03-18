Feature: Bookings

  Background:
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following users:
      | Email                | Account          | UnreadBookingsCount |
      | diego@selzlein.com   |                  | 0                   |
      | djeison@selzlein.com | Diego Apartments | 5                   |
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |
    And the following properties:
      | ID | Account          | Name          | Country | City  |
      | 19 | Diego Apartments | ACME Downtown | Japan   | Osaka |
    And the following units:
      | Property      | Name         | Bedrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | ACME Downtown | Standard Apt | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
    And the following prices:
      | Unit         | MinNights | Cents |
      | Standard Apt | 1         | 11000 |
    And I am authenticated with "diego@selzlein.com"

  Scenario: Listing my Bookings
    Given the following bookings:
      | ID | User                 | Unit         | CheckIn    | CheckOut   | Nights |
      | 1  | diego@selzlein.com   | Standard Apt | 2018-06-09 | 2018-06-11 | 2      |
      | 2  | diego@selzlein.com   | Standard Apt | 2018-06-15 | 2018-06-16 | 1      |
      | 3  | djeison@selzlein.com | Standard Apt | 2018-07-09 | 2018-07-11 | 2      |
    When I list my bookings
    Then I should receive an "OK" response with the following JSON:
      """
      [{
        "id": 2,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-15",
        "checkOut": "2018-06-16",
        "nights": 1
      }, {
        "id": 1,
        "propertyName": "ACME Downtown",
        "unitName": "Standard Apt",
        "checkIn": "2018-06-09",
        "checkOut": "2018-06-11",
        "nights": 2
      }]
      """

  Scenario: Paginating my Bookings
    Given the following bookings:
      | ID | User               | Unit         | CheckIn    | CheckOut   | Nights |
      | 1  | diego@selzlein.com | Standard Apt | 2018-06-09 | 2018-06-11 | 2      |
      | 2  | diego@selzlein.com | Standard Apt | 2018-06-15 | 2018-06-16 | 1      |
    When I list my bookings for page "2"
    Then I should receive an "OK" response with the following JSON:
      """
      []
      """

  Scenario: Booking a Unit
    Given it is currently "01 Jun 18 08:00"
    When I create the following booking:
      | Unit     | Standard Apt |
      | CheckIn  | 2018-06-09   |
      | CheckOut | 2018-06-11   |
      | Message  | Nothing      |
    Then I should receive a "CREATED" response
    And I should have the following bookings:
      | User               | Unit         | CheckIn    | CheckOut   | Message | Nights | PerNightCents | ServiceFeeCents | TotalCents |
      | diego@selzlein.com | Standard Apt | 2018-06-09 | 2018-06-11 | Nothing | 2      | 11000         | 0               | 22000      |
    And "djeison@selzlein.com" should have received the following email:
      """
      <<TEXT>>
      Dear property owner,

      You have received a new booking request:

      * User: diego@selzlein.com
      * Unit: ACME Downtown > Standard Apt
      * Check In: Sat, 09 Jun 2018
      * Check Out: Mon, 11 Jun 2018
      * Nights: 2
      * Price: $110.00/night
      * Total: $220.00

      Respond to this email to get in touch with them.
      <<HTML>>
      <div style="font-family: Proxima-Nova,Arial,sans-serif; line-height: 1.8; color: #464855; width: 500px; padding: 40px 20px;">
        <img style="margin: 0 auto; display: block; margin-bottom: 50px;"
          src="https://s3.amazonaws.com/fujilane-production/public/assets/leaf%402x.png" />

        <p>Dear property owner,</p>

        <p>You have received a new booking request:</p>

        <ul>
          <li>User: diego@selzlein.com</li>
          <li>Unit: ACME Downtown &gt; Standard Apt</li>
          <li>Check In: Sat, 09 Jun 2018</li>
          <li>Check Out: Mon, 11 Jun 2018</li>
          <li>Nights: 2</li>
          <li>Price: $110.00/night</li>
          <li>Total: $220.00</li>
        </ul>

        <p>Respond to this email to get in touch with them.</p>
      </div>
      """
    And I should have the following users:
      | Email                | Account          | UnreadBookingsCount |
      | diego@selzlein.com   |                  | 0                   |
      | djeison@selzlein.com | Diego Apartments | 6                   |
    And I should have the following accounts:
      | Name             | BookingsCount |
      | Diego Apartments | 1             |

  Scenario: Booking a Unit with invalid information
    Given it is currently "01 Jun 18 08:00"
    When I create the following booking:
      | Message | Nothing |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | unit is required                   |
      | check in is required               |
      | check in should be in the future   |
      | check out is required              |
      | check out should be after check in |
    And no emails should have been sent

  Scenario: Trying to Book a Unit that's not published
    Given it is currently "01 Jun 18 08:00"
    And the following units:
      | Property      | Name       | Bedrooms | SizeM2 | MaxOccupancy | Count |
      | ACME Downtown | Double Apt | 1        | 32     | 3            | 15    |
    And the following prices:
      | Unit       | MinNights | Cents |
      | Double Apt | 1         | 11000 |
    When I create the following booking:
      | Unit     | Double Apt |
      | CheckIn  | 2018-06-09 |
      | CheckOut | 2018-06-11 |
      | Message  | Nothing    |
    Then I should receive an "UNPROCESSABLE ENTITY" response with the following errors:
      | unit is invalid |
    And no emails should have been sent

  Scenario Outline: Calculating Booking Prices
    Given the following units:
      | Property      | Name            | Bedrooms | SizeM2 | MaxOccupancy | Count | PublishedAt          |
      | ACME Downtown | Specific Prices | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
      | ACME Downtown | Single Price    | 1        | 32     | 3            | 15    | 2018-06-09T15:00:00Z |
    And the following prices:
      | Unit            | MinNights | Cents   |
      | Specific Prices | 1         | 13000   |
      | Specific Prices | 2         | 20000   |
      | Specific Prices | 7         | 60000   |
      | Specific Prices | 30        | 240000  |
      | Specific Prices | 90        | 750000  |
      | Specific Prices | 180       | 1200000 |
      | Specific Prices | 365       | 2200000 |
      | Single Price    | 1         | 10000   |
    And it is currently "01 Jun 18 08:00"
    When I create the following booking:
      | Unit     | <Unit>     |
      | CheckIn  | <CheckIn>  |
      | CheckOut | <CheckOut> |
    Then I should receive a "CREATED" response
    And I should have the following bookings:
      | User               | Unit   | CheckIn   | CheckOut   | Nights   | PerNightCents | TotalCents |
      | diego@selzlein.com | <Unit> | <CheckIn> | <CheckOut> | <Nights> | <PerNight>    | <Total>    |
    And "djeison@selzlein.com" should have received the following email:
      """
      <<TEXT>>
      Dear property owner,

      You have received a new booking request:

      * User: diego@selzlein.com
      * Unit: ACME Downtown > <Unit>
      * Check In: <EmailCheckIn>
      * Check Out: <EmailCheckOut>
      * Nights: <Nights>
      * Price: $<EmailPrice>/night
      * Total: $<EmailTotal>

      Respond to this email to get in touch with them.
      <<HTML>>
      <div style="font-family: Proxima-Nova,Arial,sans-serif; line-height: 1.8; color: #464855; width: 500px; padding: 40px 20px;">
        <img style="margin: 0 auto; display: block; margin-bottom: 50px;"
          src="https://s3.amazonaws.com/fujilane-production/public/assets/leaf%402x.png" />

        <p>Dear property owner,</p>

        <p>You have received a new booking request:</p>

        <ul>
          <li>User: diego@selzlein.com</li>
          <li>Unit: ACME Downtown &gt; <Unit></li>
          <li>Check In: <EmailCheckIn></li>
          <li>Check Out: <EmailCheckOut></li>
          <li>Nights: <Nights></li>
          <li>Price: $<EmailPrice>/night</li>
          <li>Total: $<EmailTotal></li>
        </ul>

        <p>Respond to this email to get in touch with them.</p>
      </div>
      """

    Examples:
      | Unit            | CheckIn    | CheckOut   | Nights | PerNight | Total   | EmailPrice | EmailTotal | EmailCheckIn     | EmailCheckOut    |
      | Specific Prices | 2018-07-01 | 2018-07-02 | 1      | 13000    | 13000   | 130.00     | 130.00     | Sun, 01 Jul 2018 | Mon, 02 Jul 2018 |
      | Specific Prices | 2018-07-01 | 2018-07-07 | 6      | 10000    | 60000   | 100.00     | 600.00     | Sun, 01 Jul 2018 | Sat, 07 Jul 2018 |
      | Specific Prices | 2018-07-01 | 2018-07-08 | 7      | 8571     | 60000   | 85.71      | 600.00     | Sun, 01 Jul 2018 | Sun, 08 Jul 2018 |
      | Specific Prices | 2018-07-01 | 2018-07-30 | 29     | 8571     | 248571  | 85.71      | 2,485.71   | Sun, 01 Jul 2018 | Mon, 30 Jul 2018 |
      | Specific Prices | 2018-07-01 | 2018-07-31 | 30     | 8000     | 240000  | 80.00      | 2,400.00   | Sun, 01 Jul 2018 | Tue, 31 Jul 2018 |
      | Specific Prices | 2018-07-01 | 2018-09-28 | 89     | 8000     | 712000  | 80.00      | 7,120.00   | Sun, 01 Jul 2018 | Fri, 28 Sep 2018 |
      | Specific Prices | 2018-07-01 | 2018-09-29 | 90     | 8333     | 750000  | 83.33      | 7,500.00   | Sun, 01 Jul 2018 | Sat, 29 Sep 2018 |
      | Specific Prices | 2018-07-01 | 2018-12-27 | 179    | 8333     | 1491667 | 83.33      | 14,916.67  | Sun, 01 Jul 2018 | Thu, 27 Dec 2018 |
      | Specific Prices | 2018-07-01 | 2018-12-28 | 180    | 6667     | 1200000 | 66.67      | 12,000.00  | Sun, 01 Jul 2018 | Fri, 28 Dec 2018 |
      | Specific Prices | 2018-07-01 | 2019-06-30 | 364    | 6667     | 2426667 | 66.67      | 24,266.67  | Sun, 01 Jul 2018 | Sun, 30 Jun 2019 |
      | Specific Prices | 2018-07-01 | 2019-07-01 | 365    | 6027     | 2200000 | 60.27      | 22,000.00  | Sun, 01 Jul 2018 | Mon, 01 Jul 2019 |
      | Single Price    | 2018-07-01 | 2018-07-02 | 1      | 10000    | 10000   | 100.00     | 100.00     | Sun, 01 Jul 2018 | Mon, 02 Jul 2018 |
      | Single Price    | 2018-07-01 | 2018-07-08 | 7      | 10000    | 70000   | 100.00     | 700.00     | Sun, 01 Jul 2018 | Sun, 08 Jul 2018 |
      | Single Price    | 2018-07-01 | 2018-12-28 | 180    | 10000    | 1800000 | 100.00     | 18,000.00  | Sun, 01 Jul 2018 | Fri, 28 Dec 2018 |
      | Single Price    | 2018-07-01 | 2019-07-01 | 365    | 10000    | 3650000 | 100.00     | 36,500.00  | Sun, 01 Jul 2018 | Mon, 01 Jul 2019 |
