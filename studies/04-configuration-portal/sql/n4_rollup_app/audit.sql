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

-- name: a_rollup_mismatches
-- INV-9. Recomputed from config_entry, never from the stored columns: a rollup
-- checked against itself would always agree.
SELECT count(*) AS mismatches
FROM installed_product ip
WHERE ip.config_count <> (SELECT count(*) FROM config_entry c WHERE c.installed_product_id = ip.id)
   OR ip.content_hash <> COALESCE((SELECT md5(string_agg(c.key || '=' || c.value, ',' ORDER BY c.key))
                                   FROM config_entry c WHERE c.installed_product_id = ip.id), '');
