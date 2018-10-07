CREATE TABLE cities(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  name varchar not null unique,
  country_id integer not null references countries
);
