-- The read catalogue, same names and same result shapes as every other design.
-- The document is unrolled into key/value rows so that a difference between
-- designs is a difference in mechanism, not in the answer's shape.

-- name: r01_effective_config
-- params: installed_product_id
-- One row read, unrolled. One statement, so the answer comes from a single
-- snapshot and cannot mix two revisions (INV-2).
SELECT e.key, e.value, ip.current_revision AS revision
FROM config_document d
JOIN installed_product ip ON ip.id = d.installed_product_id
CROSS JOIN LATERAL jsonb_each_text(d.doc) AS e(key, value)
WHERE d.installed_product_id = $1
ORDER BY e.key;

-- name: r02_read_key
-- params: installed_product_id, key
SELECT d.doc ->> $2 AS value, ip.current_revision AS revision
FROM config_document d
JOIN installed_product ip ON ip.id = d.installed_product_id
WHERE d.installed_product_id = $1;

-- name: r03_revision_check
-- params: installed_product_id, known_revision
SELECT ip.current_revision AS revision,
       (ip.current_revision > $2) AS newer
FROM installed_product ip
WHERE ip.id = $1;

-- name: r04_list_installations
-- params: none
-- Counting the document's keys is not free here: it is a scan of the JSON object,
-- where the row-per-key design counts rows and the rollup designs read a column.
SELECT ip.id,
       pd.name AS product_definition,
       e.name  AS environment,
       bu.name AS business_unit,
       (SELECT count(*) FROM jsonb_object_keys(d.doc)) AS config_count,
       ip.current_revision,
       ip.updated_at AS last_modified
FROM installed_product ip
JOIN product_definition pd ON pd.id = ip.product_definition_id
JOIN environment e         ON e.id  = ip.environment_id
LEFT JOIN business_unit bu ON bu.id = ip.business_unit_id
LEFT JOIN config_document d ON d.installed_product_id = ip.id
ORDER BY ip.id;

-- name: r05_search_key
-- params: scope_kind, scope_id, key
-- A containment test over every document in scope. There is no index that can
-- answer it, so this is the design's worst read and is kept in the catalogue
-- precisely because it must be shown rather than assumed.
SELECT d.installed_product_id AS installed_product_id, d.doc ->> $3 AS value
FROM config_document d
JOIN installed_product ip ON ip.id = d.installed_product_id
WHERE d.doc ? $3
  AND (   ($1 = 'product_definition' AND ip.product_definition_id = $2)
       OR ($1 = 'environment'        AND ip.environment_id = $2)
       OR ($1 = 'business_unit'      AND ip.business_unit_id = $2))
ORDER BY d.installed_product_id;
