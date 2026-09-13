-- Kept identical to C1's indexes even though the booking path no longer counts
-- through ticket_event_idx: every created-ticket design maintains the same
-- indexes on each insert, so a difference in booking cost is the concurrency
-- strategy's, not an index's. (Listing an event's tickets needs it anyway.)
CREATE INDEX ticket_event_idx ON ticket (event_id);

CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX event_band_idx ON event (band_id, starts_at);
