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
