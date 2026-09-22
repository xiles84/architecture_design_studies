-- Indexes for the owned model.
--
-- The five D3 indexes are byte-identical to the legacy model's, so the index
-- maintenance cost on every write is held constant and the version/outbox work is
-- the only write-path difference between the pair. One index is added on the
-- outbox, because a change log without an index on (person_id, seq) is not a
-- usable change log -- and pretending it is would make the owned model look
-- cheaper than any real deployment.

CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);
CREATE INDEX donation_charity_amount_idx ON donation (charity_id) INCLUDE (amount_cents);

CREATE INDEX cache_outbox_person_seq_idx ON cache_outbox (person_id, seq);
