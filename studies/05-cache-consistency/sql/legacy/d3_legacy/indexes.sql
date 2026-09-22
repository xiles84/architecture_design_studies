-- Indexes for the legacy D3 model. Applied AFTER the bulk load, so the index
-- maintenance cost is not folded into the loader's time (study 01's convention).
--
-- These are D3's own indexes and nothing is added for the cache: the legacy model
-- cannot be changed for caching, and adding an index here would be a change.

CREATE INDEX person_charity_idx ON person (charity_id);

-- The portal's recent-20 read and its aggregate both start from (person_id, time).
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);

CREATE INDEX donation_time_idx ON donation (donated_at DESC);

-- Rolldown's payoff: a charity-scoped recency read never joins to person.
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);

-- Covering index so "total donated to a charity" can be an index-only scan.
CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);
