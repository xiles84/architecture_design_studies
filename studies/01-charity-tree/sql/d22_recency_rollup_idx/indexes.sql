CREATE INDEX person_charity_idx ON person (charity_id);
CREATE INDEX donation_person_time_idx ON donation (person_id, donated_at DESC);
CREATE INDEX donation_time_idx ON donation (donated_at DESC);
CREATE INDEX donation_charity_time_idx ON donation (charity_id, donated_at DESC);

-- "Who donates the most" becomes an index range scan over the rollup column
-- instead of a GROUP BY over every donation.
CREATE INDEX person_top_donor_idx ON person (charity_id, total_donated_cents DESC);

-- "Who donated most recently" likewise.
CREATE INDEX person_last_donation_idx ON person (charity_id, last_donation_at DESC);

-- New in D22: D4 already STORES last_donation_at, but the index above is
-- charity-scoped. q13/q14/q16 (RECENCY.md) ask across every charity, and the
-- charity-scoped index gives them no useful access path. This one line is
-- D22's entire change over D4 -- see D20 -> D22 in RECENCY.md section 4 for
-- what it costs against materialising the fact on the child instead.
CREATE INDEX person_last_donation_global_idx ON person (last_donation_at DESC);
