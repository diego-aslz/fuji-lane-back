ALTER TABLE properties ADD first_published_at timestamp without time zone;
ALTER TABLE units ADD first_published_at timestamp without time zone;

UPDATE properties SET first_published_at = COALESCE(published_at, created_at) WHERE ever_published;
UPDATE units SET first_published_at = COALESCE(published_at, created_at) WHERE ever_published;

ALTER TABLE properties DROP ever_published;
ALTER TABLE units DROP ever_published;
