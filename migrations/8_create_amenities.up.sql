CREATE TABLE amenities(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  type varchar not null,
  name varchar check (type <> 'custom' or name is not null),
  property_id integer references properties,
  unit_id integer references units check (property_id is not null or unit_id is not null)
);

CREATE INDEX amenities_property_id ON amenities(property_id);
CREATE INDEX amenities_unit_id ON amenities(unit_id);
