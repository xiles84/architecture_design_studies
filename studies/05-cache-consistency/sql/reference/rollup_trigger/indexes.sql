-- rollup-trigger indexes: byte-identical to the flattened baseline, so the only
-- difference between the two cells is the stored aggregate and who maintains it.

CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);
CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);
