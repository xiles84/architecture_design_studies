CREATE INDEX person_charity_idx ON person (charity_id);

-- A GIN index makes containment ("which person has donation id X?") possible
-- without a sequential scan. It cannot help with ordering or aggregation, which
-- is exactly the limitation the embedded design runs into.
CREATE INDEX person_donations_gin ON person USING GIN (donations jsonb_path_ops);
