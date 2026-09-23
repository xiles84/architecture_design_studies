-- The read catalogue, same names and same result shapes as every other design.
-- The document is unrolled into key/value rows so that a difference between
-- designs is a difference in mechanism, not in the answer's shape.
--
-- Every read here touches only `installed_product`: the document is a column on
-- that row, so there is no join to a separate document table. That is the point
-- of the design -- and the reason its metadata write is expensive instead.

-- name: r01_effective_config
-- params: installed_product_id
-- One row read, unrolled. One statement, so the answer comes from a single
-- snapshot and cannot mix two revisions (INV-2).
SELECT e.key, e.value, ip.current_revision AS revision
FROM installed_product ip
CROSS JOIN LATERAL jsonb_each_text(ip.doc) AS e(key, value)
WHERE ip.id = $1
ORDER BY e.key;

-- name: r02_read_key
-- params: installed_product_id, key
SELECT ip.doc ->> $2 AS value, ip.current_revision AS revision
FROM installed_product ip
WHERE ip.id = $1;

-- name: r03_revision_check
-- params: installed_product_id, known_revision
SELECT ip.current_revision AS revision,
       (ip.current_revision > $2) AS newer
FROM installed_product ip
WHERE ip.id = $1;

-- name: r04_list_installations
-- params: none
-- Counting the document's keys is a scan of the JSON object, exactly as in d1;
-- the difference from d1 is only where the object is stored.
SELECT ip.id,
       pd.name AS product_definition,
       e.name  AS environment,
       bu.name AS business_unit,
       (SELECT count(*) FROM jsonb_object_keys(ip.doc)) AS config_count,
       ip.current_revision,
       ip.updated_at AS last_modified
FROM installed_product ip
JOIN product_definition pd ON pd.id = ip.product_definition_id
JOIN environment e         ON e.id  = ip.environment_id
LEFT JOIN business_unit bu ON bu.id = ip.business_unit_id
ORDER BY ip.id;

-- name: r05_search_key
-- params: scope_kind, scope_id, key
-- A containment test over every parent row in scope. There is no index that can
-- answer it, so this is the design's worst read and is kept precisely because it
-- must be shown rather than assumed.
SELECT ip.id AS installed_product_id, ip.doc ->> $3 AS value
FROM installed_product ip
WHERE ip.doc ? $3
  AND (   ($1 = 'product_definition' AND ip.product_definition_id = $2)
       OR ($1 = 'environment'        AND ip.environment_id = $2)
       OR ($1 = 'business_unit'      AND ip.business_unit_id = $2))
ORDER BY ip.id;
