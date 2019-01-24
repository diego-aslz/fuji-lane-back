CREATE TABLE units(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  property_id bigint not null references properties,
  published_at timestamp without time zone,
  ever_published boolean not null default false,
  name varchar not null,
  slug varchar not null,
  bedrooms int not null,
  bathrooms int not null,
  size_m2 int not null,
  max_occupancy int,
  count int not null,
  base_price_cents int,
  one_night_price_cents int,
  one_week_price_cents int,
  three_months_price_cents int,
  six_months_price_cents int,
  twelve_months_price_cents int,
  overview text
);

CREATE INDEX units_property_id ON units(property_id);
CREATE INDEX units_published_at ON units(published_at);
CREATE INDEX units_deleted_at ON units(deleted_at);

CREATE UNIQUE INDEX units_property_id_slug ON units(property_id, slug);
