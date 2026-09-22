-- embedded-locked indexes: byte-identical to the flattened baseline, so the only
-- difference between the cells is where the recent slice lives.
--
-- The parent's embedded array is read with the parent row, so it needs no index.
-- The child table keeps every index it had, because it is still the source of
-- truth for every other question and for the audits.

CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);
CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);
