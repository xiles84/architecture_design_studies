-- The design's own account of the truth, read from the tables it actually
-- maintains rather than from the query the read benchmark just timed.

-- name: a_key_values
-- params: installed_product_id
SELECT key, value FROM config_entry WHERE installed_product_id = $1 ORDER BY key;

-- name: a_revision
-- params: installed_product_id
SELECT current_revision AS revision FROM installed_product WHERE id = $1;

-- name: a_counts
-- params: installed_product_id
SELECT count(*) AS entries FROM config_entry WHERE installed_product_id = $1;

-- name: a_rolldown_mismatches
-- INV-8: every copied parent key must still equal the parent it was copied from.
-- A rolldown that has drifted is not a cache miss, it is a wrong answer to r05.
SELECT count(*) AS mismatches
FROM config_entry c
JOIN installed_product ip ON ip.id = c.installed_product_id
WHERE c.product_definition_id IS DISTINCT FROM ip.product_definition_id
   OR c.environment_id        IS DISTINCT FROM ip.environment_id
   OR c.business_unit_id      IS DISTINCT FROM ip.business_unit_id;
