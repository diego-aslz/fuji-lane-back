CREATE TABLE amenities(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  type varchar not null,
  name varchar check (type <> 'custom' or name is not null),
  property_id bigint references properties,
  unit_id bigint references units check (property_id is not null or unit_id is not null)
);

CREATE UNIQUE INDEX amenities_name_uniq_by_property ON amenities(property_id, name);
CREATE UNIQUE INDEX amenities_name_uniq_by_unit ON amenities(unit_id, name);
