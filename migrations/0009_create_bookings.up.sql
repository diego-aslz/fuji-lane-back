CREATE TABLE bookings(
  id bigserial primary key not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  user_id bigint not null references users,
  unit_id bigint not null references units,
  check_in date not null,
  check_out date not null,
  message text,
  per_night_cents int not null,
  nights int not null check (nights > 0),
  service_fee_cents int not null,
  total_cents int not null
);

CREATE INDEX bookings_user_id ON bookings(user_id);
CREATE INDEX bookings_unit_id ON bookings(unit_id);
