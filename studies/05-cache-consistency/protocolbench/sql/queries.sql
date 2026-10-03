-- name: snapshot
-- params: key
-- One statement observes the retained generation and its committed payload together.
SELECT g.generation, coalesce(o.value,0), coalesce(o.deleted,true)
FROM protocol_generation g LEFT JOIN protocol_object o USING(key_id) WHERE g.key_id=$1;

-- name: version
-- params: key
SELECT generation FROM protocol_generation WHERE key_id=$1;

-- name: lock
-- params: key
SELECT generation FROM protocol_generation WHERE key_id=$1 FOR UPDATE;

-- name: cas
-- params: key, expected
UPDATE protocol_generation SET generation=generation+1 WHERE key_id=$1 AND generation=$2 RETURNING generation;

-- name: bump
-- params: key
UPDATE protocol_generation SET generation=generation+1 WHERE key_id=$1 RETURNING generation;

-- name: upsert
-- params: key, generation, value, deleted
INSERT INTO protocol_object VALUES ($1,$2,$3,$4) ON CONFLICT(key_id) DO UPDATE SET generation=excluded.generation,value=excluded.value,deleted=excluded.deleted;

-- name: delete
-- params: key
DELETE FROM protocol_object WHERE key_id=$1;

-- name: event
-- params: key, generation
-- Every candidate pays this same transactional source-side cost.
INSERT INTO protocol_outbox(key_id,generation) VALUES($1,$2);

-- name: pending
-- params: none
-- A sequence cursor can miss a transaction committing after a later sequence.
SELECT seq,key_id,generation FROM protocol_outbox WHERE NOT delivered ORDER BY seq LIMIT 32 FOR UPDATE SKIP LOCKED;

-- name: delivered
-- params: sequence
UPDATE protocol_outbox SET delivered=true WHERE seq=$1;
