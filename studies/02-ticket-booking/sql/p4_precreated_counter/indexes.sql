-- The natural key of a pre-created seat. It does not prevent overbooking (the
-- rows already exist; selling updates them), it prevents publishing a seat twice.
CREATE UNIQUE INDEX ticket_event_seat_uq ON ticket (event_id, seat_no);

-- The booking path's index: only unsold seats, in seat order within an event.
-- Partial, so a nearly sold-out 100 000-seat event has a small index to search
-- -- in principle. In practice a sold seat leaves a dead index entry behind until
-- vacuum removes it, so the scan for "the next available seat" walks over the
-- entries of every seat sold since the last vacuum. That is the queue-table
-- problem, and the race measures it.
CREATE INDEX ticket_available_idx ON ticket (event_id, seat_no) WHERE status = 'available';

CREATE INDEX ticket_customer_idx ON ticket (customer_id) WHERE status = 'sold';

CREATE INDEX event_band_idx ON event (band_id, starts_at);
