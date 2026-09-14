-- L2 — section-document: indexes, built after the load.

-- The sweeper's index: sections whose earliest held claim has expired. ASC so
-- YugabyteDB range-shards it. (Sections with no held claim have NULL and are
-- never returned by next_expiry <= now().)
CREATE INDEX event_section_next_expiry_idx ON event_section (next_expiry ASC);

-- "My tickets", and the audit's tickets per seat. Same in every design.
CREATE INDEX ticket_customer_idx ON ticket (customer_id);

CREATE INDEX ticket_event_seat_idx ON ticket (event_id, seat_id);
