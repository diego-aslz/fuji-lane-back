CREATE TABLE bookings(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  user_id integer not null references users,
  unit_id integer not null references units,
  check_in date not null,
  check_out date not null,
  additional_info text,
  night_price_cents integer not null,
  nights integer not null check (nights > 0),
  service_fee_cents integer not null,
  total_cents integer not null
);

CREATE INDEX bookings_user_id ON bookings(user_id);
CREATE INDEX bookings_unit_id ON bookings(unit_id);
