CREATE TABLE countries(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  name varchar not null unique
);
