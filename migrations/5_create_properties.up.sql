CREATE TABLE properties(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  account_id int not null references accounts,
  state smallint not null default 1,
  name varchar,
  address_1 varchar,
  address_2 varchar,
  address_3 varchar,
  city_id int references cities,
  postal_code varchar,
  country_id int references countries
);

CREATE INDEX properties_account_id_idx ON properties(account_id);
CREATE INDEX properties_state_idx ON properties(state);
CREATE INDEX properties_deleted_at ON properties(deleted_at);
