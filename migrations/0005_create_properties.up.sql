CREATE TABLE properties(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  account_id bigint not null references accounts,
  published_at timestamp without time zone,
  ever_published boolean not null default false,
  name varchar unique,
  slug varchar unique,
  address1 varchar,
  address2 varchar,
  address3 varchar,
  city_id int references cities,
  postal_code varchar,
  country_id int references countries,
  latitude decimal,
  longitude decimal,
  minimum_stay smallint,
  deposit varchar,
  cleaning varchar,
  nearest_airport varchar,
  nearest_subway varchar,
  nearby_locations varchar,
  overview varchar
);

CREATE INDEX properties_account_id ON properties(account_id);
CREATE INDEX properties_city_id ON properties(city_id);
CREATE INDEX properties_published_at ON properties(published_at);
CREATE INDEX properties_deleted_at ON properties(deleted_at);
