Feature: Amenities

  Scenario: Listing amenity types for properties
    When I list amenity types for "properties"
    Then I should receive an "OK" response with the following amenity types:
      | Code          | Name          |
      | daycare       | Daycare       |
      | gym           | Gym           |
      | meeting_rooms | Meeting Rooms |
      | pool          | Pool          |
      | restaurant    | Restaurant    |

  Scenario: Listing amenity types for units
    When I list amenity types for "units"
    Then I should receive an "OK" response with the following amenity types:
      | Code              | Name               |
      | air_conditioning  | Air Conditioning   |
      | bathrobes         | Bathrobes          |
      | blackout_curtains | Blackout Curtains  |
      | housekeeping      | Daily Housekeeping |
      | desk              | Desk               |
      | dvd               | DVD Player         |
      | minibar           | Minibar            |
      | phone             | Phone              |
      | toilet            | Toilet             |
