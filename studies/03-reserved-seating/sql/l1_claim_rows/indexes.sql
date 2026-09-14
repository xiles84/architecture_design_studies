-- L1 — claim-rows: indexes, built after the load.

-- A section's seat map: the claims of one section.
CREATE INDEX seat_claim_section_idx ON seat_claim (event_id, section_no, seat_id);

-- The sweeper's index: claims in expiry order ('infinity' sorts last and is never
-- reached). ASC so YugabyteDB range-shards it.
CREATE INDEX seat_claim_expiry_idx ON seat_claim (expires_at ASC);

-- "My tickets". The audit's tickets per seat use the UNIQUE constraint's index.
CREATE INDEX ticket_customer_idx ON ticket (customer_id);
