-- S1r — confirm-retry: indexes, built after the load.

-- A section's seat map, and the block a buyer holds or confirms. On YugabyteDB
-- the first column is hash-sharded, so all of one event's seats live together;
-- L3 changes exactly that.
CREATE INDEX event_seat_section_idx ON event_seat (event_id, section_no, seat_id);

-- The sweeper's index: holds in expiry order. Partial, so sold and available
-- seats cost it nothing. ASC is spelled out so YugabyteDB range-shards it: a
-- hash index cannot answer "hold_expires_at <= now()".
CREATE INDEX event_seat_expiry_idx ON event_seat (hold_expires_at ASC) WHERE status = 'held';

-- "My tickets", and the audit's tickets per seat. Same in every design.
CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX ticket_event_seat_idx ON ticket (event_id, seat_id);
