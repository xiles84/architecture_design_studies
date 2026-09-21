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
