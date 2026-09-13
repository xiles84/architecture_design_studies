-- The arbiter. Built after the bulk load like every other index; the loaded
-- data already satisfies it, and the build would fail loudly if it did not.
-- It also serves the availability COUNT and the MAX(seat_no) lookup, so C5
-- maintains one fewer index than C1 on every insert.
CREATE UNIQUE INDEX ticket_event_seat_uq ON ticket (event_id, seat_no);

CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX event_band_idx ON event (band_id, starts_at);
