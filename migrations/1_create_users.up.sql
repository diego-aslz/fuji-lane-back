CREATE TABLE users(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  deleted_at timestamp without time zone,
  email varchar not null,
  name varchar not null,
  facebook_id varchar,
  encrypted_password varchar,
  last_signed_in timestamp without time zone
);

CREATE UNIQUE INDEX users_email_unique ON users(email);
CREATE INDEX users_deleted_at ON users(deleted_at);
