-- The index the capacity check counts through. Its cost per booking is
-- proportional to the tickets already sold for the event: cheap for a 10-seat
-- gig, a 100 000-entry index scan for the last ticket of a stadium.
CREATE INDEX ticket_event_idx ON ticket (event_id);

CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX event_band_idx ON event (band_id, starts_at);
