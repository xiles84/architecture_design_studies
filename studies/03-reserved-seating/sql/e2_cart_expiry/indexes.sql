-- E2 — cart-expiry: indexes, built after the load.

-- A section's seat map, and the block a buyer holds or confirms (as S1).
CREATE INDEX event_seat_section_idx ON event_seat (event_id, section_no, seat_id);

-- The seats of one cart: the sweeper releases an expired cart's seats through it.
CREATE INDEX event_seat_hold_idx ON event_seat (hold_id) WHERE status = 'held';

-- The sweeper's index: carts in expiry order. ASC so YugabyteDB range-shards it.
CREATE INDEX hold_expiry_idx ON hold (expires_at ASC) WHERE status = 'held';

-- "My tickets", and the audit's tickets per seat. Same in every design.
CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX ticket_event_seat_idx ON ticket (event_id, seat_id);
