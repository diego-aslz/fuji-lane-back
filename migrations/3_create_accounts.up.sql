CREATE TABLE accounts(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  status smallint not null,
  name varchar not null unique,
  phone varchar,
  country_id int references countries
);

CREATE INDEX accounts_deleted_at ON accounts(deleted_at);
