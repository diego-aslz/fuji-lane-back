CREATE TABLE units(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  property_id integer not null references properties,
  published_at timestamp without time zone,
  name varchar not null,
  bedrooms integer not null,
  bathrooms integer not null,
  size_m2 integer not null,
  max_occupancy integer,
  count integer not null,
  base_price_cents integer,
  one_night_price_cents integer,
  one_week_price_cents integer,
  three_months_price_cents integer,
  six_months_price_cents integer,
  twelve_months_price_cents integer,
  overview text
);

CREATE INDEX units_property_id ON units(property_id);
CREATE INDEX units_published_at ON units(published_at);
CREATE INDEX units_deleted_at ON units(deleted_at);
