CREATE TABLE images(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  name varchar not null,
  type varchar not null,
  size integer not null,
  url varchar not null,
  property_id integer references properties,
  unit_id integer references units,
  uploaded boolean default false
);

ALTER TABLE units ADD floor_plan_image_id int references images;

CREATE INDEX images_property_id ON images(property_id);
CREATE INDEX images_unit_id ON images(unit_id);
