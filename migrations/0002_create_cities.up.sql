CREATE TABLE cities(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  name varchar not null unique,
  slug varchar not null unique,
  country_id int not null references countries,
  latitude decimal not null,
  longitude decimal not null
);
