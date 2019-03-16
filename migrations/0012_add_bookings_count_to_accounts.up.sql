ALTER TABLE accounts ADD bookings_count INT NOT NULL DEFAULT 0;

UPDATE accounts SET bookings_count = (
  SELECT COUNT(*) FROM bookings
  INNER JOIN units ON units.id = bookings.unit_id
  INNER JOIN properties ON properties.id = units.property_id
  WHERE properties.account_id = accounts.id
);
