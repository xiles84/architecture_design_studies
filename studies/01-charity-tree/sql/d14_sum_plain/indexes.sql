-- Foreign-key support: makes parent->child navigation and cascading checks cheap.
CREATE INDEX person_charity_idx ON person (charity_id);

-- The workhorse. Covers "donations of a person, newest first", "first/last donation
-- of a person" (via forward and backward scan) and "count donations of a person".
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);

-- Global recency: "the last donation" without scanning the table.
CREATE INDEX donation_time_idx ON donation (donated_at DESC);

CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);
CREATE INDEX donation_charity_amount_idx ON donation (charity_id);
