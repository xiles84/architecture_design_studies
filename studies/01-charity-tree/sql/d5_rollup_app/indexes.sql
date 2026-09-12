CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);

-- "Who donates the most" becomes an index range scan over the rollup column
-- instead of a GROUP BY over every donation.
CREATE INDEX person_top_donor_idx ON person (charity_id, total_donated_cents DESC);

-- "Who donated most recently" likewise.
CREATE INDEX person_last_donation_idx ON person (charity_id, last_donation_at DESC);
