-- y2_noncolocated indexes: byte-identical to the flattened baseline's, and to
-- y1_colocated's after accounting for the fact that y1's key choice moves the
-- identity lookup into a secondary index. The pair holds the charity-scoped and
-- time indexes constant so the measured difference is the placement of a donor's
-- donations and nothing else.

CREATE INDEX person_charity_idx ON person (charity_id);

CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);

CREATE INDEX donation_time_idx ON donation (donated_at DESC);

CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);

CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);
