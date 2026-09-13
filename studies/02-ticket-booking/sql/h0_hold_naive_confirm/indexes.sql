-- Same ticket and event indexes as every created-ticket design.
CREATE INDEX ticket_event_idx ON ticket (event_id);

CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX event_band_idx ON event (band_id, starts_at);

-- The sweeper's index: live holds in expiry order. ASC is spelled out because
-- YugabyteDB would otherwise hash-shard the first column, and a hash index
-- cannot answer "expires_at <= now()".
CREATE INDEX reservation_expiry_idx ON reservation (expires_at ASC) WHERE status = 'held';
