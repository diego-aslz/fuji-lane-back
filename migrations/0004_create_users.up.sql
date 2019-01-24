CREATE TABLE users(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  account_id bigint references accounts,
  email varchar not null unique,
  name varchar,
  facebook_id varchar,
  encrypted_password varchar,
  last_signed_in timestamp without time zone
);

CREATE INDEX users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX users_email ON users(email);
