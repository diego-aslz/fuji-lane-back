ALTER TABLE units ADD size_ft2 INTEGER;

UPDATE units SET size_ft2 = size_m2 * 10.764;

ALTER TABLE units ALTER size_ft2 SET NOT NULL;
