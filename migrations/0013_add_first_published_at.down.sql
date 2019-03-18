ALTER TABLE properties ADD ever_published boolean;
ALTER TABLE units ADD ever_published boolean;

UPDATE properties SET ever_published = t WHERE first_published_at IS NOT NULL;
UPDATE units SET ever_published = t WHERE first_published_at IS NOT NULL;

ALTER TABLE properties DROP first_published_at;
ALTER TABLE units DROP first_published_at;
