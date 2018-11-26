CREATE TABLE properties(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  account_id int not null references accounts,
  state smallint not null default 1,
  name varchar unique,
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
CREATE INDEX properties_state ON properties(state);
CREATE INDEX properties_deleted_at ON properties(deleted_at);
