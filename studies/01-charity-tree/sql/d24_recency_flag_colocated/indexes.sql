CREATE INDEX person_charity_idx ON person (charity_id);

-- donation_id is no longer the primary key, so identity lookups need this.
-- In YugabyteDB a secondary index is a distributed table of its own: the lookup
-- hits the index tablet, then hops to the base-row tablet.
CREATE UNIQUE INDEX donation_id_idx ON donation (donation_id HASH);

-- Charity-scoped recency. Range-sharded so that a descending scan over
-- donated_at stays ordered within the index.
CREATE INDEX donation_charity_time_idx ON donation (charity_id HASH, donated_at DESC);

-- Global recency. A single range-sharded index over time; note this creates a
-- hot write tablet at the "now" end of the range, which is the classic
-- monotonic-key hotspot and is worth seeing in the write numbers.
CREATE INDEX donation_time_idx ON donation (donated_at DESC);

-- New in D24, byte-identical in shape to D20's: a PARTIAL index, since
-- is_last_donation is TRUE for exactly one row per donor. A global q13/q14/q16
-- still scatters across every tablet regardless of this index's own sharding --
-- the placement question this design isolates is what colocating each donor's
-- OWN donations (the primary key above) does to the flag's maintenance cost,
-- not to this index's shape.
CREATE INDEX donation_last_flag_idx ON donation (donated_at DESC) WHERE is_last_donation;
