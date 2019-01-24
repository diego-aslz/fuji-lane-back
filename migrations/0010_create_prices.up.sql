CREATE TABLE prices(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  unit_id bigint not null references units,
  min_nights smallint not null,
  cents int not null
);

CREATE UNIQUE INDEX prices_unit_id_min_nights ON prices(unit_id, min_nights);
