CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);

-- New in D3: charity-scoped recency without a join to person.
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);

-- Covering index so "total donated to a charity" can be an index-only scan:
-- the amount travels in the index leaf, so the heap is never visited.
CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);

-- New in D20: a PARTIAL index, since is_last_donation is TRUE for exactly one
-- row per donor -- this index stays small (one row per donor) no matter how
-- long the study's history grows, which is the whole appeal of the flag over
-- a full index on donated_at.
CREATE INDEX donation_last_flag_idx ON donation (donated_at DESC) WHERE is_last_donation;
