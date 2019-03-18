Feature: Properties Sitemap

  Background:
    And the following countries:
      | ID | Name  |
      | 2  | Japan |
    And the following cities:
      | ID | Country | Name  |
      | 3  | Japan   | Osaka |

  Scenario: Exporting properties sitemap
    Given the following accounts:
      | Name             |
      | Diego Apartments |
    And the following properties:
      | ID | Account          | Name          | PublishedAt          | UpdatedAt            | Country | City  |
      | 1  | Diego Apartments | ACME Downtown | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z | Japan   | Osaka |
      | 2  | Diego Apartments | Unpublished   |                      | 2018-06-05T08:00:00Z | Japan   | Osaka |
    And the following units:
      | ID | Property      | Name         | PublishedAt          | UpdatedAt            |
      | 11 | ACME Downtown | Standard Apt | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z |
      | 12 | Unpublished   | Double Apt   | 2018-06-05T08:00:00Z | 2018-06-05T08:00:00Z |
    When I get properties sitemap
    Then I should receive an "OK" response with the following XML:
      """
      <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
        <url>
          <loc>http://test.fujilane.com/listings/acme-downtown-1</loc>
          <lastmod>2018-06-05T08:00:00Z</lastmod>
          <changefreq>daily</changefreq>
        </url>
      </urlset>
      """
