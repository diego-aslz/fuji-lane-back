Feature: Amenities

  Scenario: Listing amenity types for properties
    When I list amenity types for "properties"
    Then I should receive an "OK" response with "60" amenity types like:
      | Code    | Name    |
      | daycare | Daycare |
      | gym     | Gym     |
      | pool    | Pool    |

  Scenario: Listing amenity types for units
    When I list amenity types for "units"
    Then I should receive an "OK" response with "60" amenity types like:
      | Code             | Name             |
      | air-conditioning | Air Conditioning |
      | bathrobes        | Bathrobes        |
      | desk             | Desk             |
      | phone            | Phone            |
