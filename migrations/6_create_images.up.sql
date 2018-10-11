CREATE TABLE images(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  name varchar not null,
  type varchar not null,
  size integer not null,
  url varchar not null,
  property_id integer not null references properties,
  uploaded boolean default false
);

CREATE INDEX images_property_id_idx ON images(property_id);
