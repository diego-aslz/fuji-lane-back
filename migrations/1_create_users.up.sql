CREATE TABLE users(
  id bigserial primary key not null,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  email varchar not null,
  name varchar not null,
  facebook_id varchar,
  last_signed_in timestamp with time zone
);

CREATE UNIQUE INDEX users_email_unique ON users(email);
CREATE INDEX users_deleted_at ON users(deleted_at);
