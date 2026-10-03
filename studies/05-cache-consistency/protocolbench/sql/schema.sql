-- A bounded keyed projection isolates coherence from the donor portal's query cost.
-- The retained identity makes physical deletion/recreation monotonic instead of ABA.
DROP TABLE IF EXISTS protocol_outbox;
DROP TABLE IF EXISTS protocol_object;
DROP TABLE IF EXISTS protocol_generation;
CREATE TABLE protocol_generation (key_id bigint PRIMARY KEY, generation bigint NOT NULL);
CREATE TABLE protocol_object (key_id bigint PRIMARY KEY, generation bigint NOT NULL, value bigint NOT NULL, deleted boolean NOT NULL);
CREATE TABLE protocol_outbox (seq bigserial PRIMARY KEY, key_id bigint NOT NULL, generation bigint NOT NULL, delivered boolean NOT NULL DEFAULT false);
CREATE INDEX protocol_outbox_pending ON protocol_outbox(seq) WHERE NOT delivered;
