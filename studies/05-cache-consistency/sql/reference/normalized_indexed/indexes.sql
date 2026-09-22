-- normalized-indexed indexes.
--
-- The D3 indexes that still apply, with the two charity-keyed ones replaced by the
-- join's input: person(charity_id). The flattened model does not need
-- person_charity_idx for the charity feed, because the copied donation.charity_id
-- answers it directly -- that index is present in both cells only because the
-- portal's own reads and the person-level aggregates use it.

CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
