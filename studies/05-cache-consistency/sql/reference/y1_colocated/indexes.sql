-- y1_colocated indexes.
--
-- donation_id is no longer the primary key, so identity lookups need a secondary
-- index -- and in YugabyteDB that index is a distributed table of its own, so the
-- lookup is a two-hop operation. The charity-scoped index is range-sharded so a
-- descending scan over donated_at stays ordered inside the index.
--
-- The extra index over `donation_time_idx (donated_at DESC)` creates a hot write
-- tablet at the "now" end of the range: the classic monotonic-key hotspot. It is
-- kept because the flattened baseline has the same index, so the pair stays
-- controlled; the write numbers will show the hotspot on both cells.

CREATE INDEX person_charity_idx ON person (charity_id);

CREATE UNIQUE INDEX donation_id_idx ON donation (donation_id HASH);

CREATE INDEX donation_charity_time_idx ON donation (charity_id HASH, donated_at DESC);

CREATE INDEX donation_time_idx ON donation (donated_at DESC);
