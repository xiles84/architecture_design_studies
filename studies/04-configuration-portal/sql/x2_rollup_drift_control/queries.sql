-- The read catalogue. Every design answers these five questions with the same
-- statement names and the same result shapes, so a difference between designs is
-- a difference in mechanism rather than in what was asked.

-- name: r01_effective_config
-- params: installed_product_id
-- The complete effective configuration. One statement, so the answer comes from a
-- single snapshot and can never mix two revisions (INV-2). The revision comes back
-- beside the rows so the harness knows which revision it read.
SELECT c.key, c.value, ip.current_revision AS revision
FROM config_entry c
JOIN installed_product ip ON ip.id = c.installed_product_id
WHERE c.installed_product_id = $1
ORDER BY c.key;

-- name: r02_read_key
-- params: installed_product_id, key
SELECT c.value, ip.current_revision AS revision
FROM config_entry c
JOIN installed_product ip ON ip.id = c.installed_product_id
WHERE c.installed_product_id = $1 AND c.key = $2;

-- name: r03_revision_check
-- params: installed_product_id, known_revision
-- "Has anything changed since revision N?" is answerable from the parent row
-- alone. A client must not have to download a configuration to learn that there
-- is nothing new.
SELECT ip.current_revision AS revision,
       (ip.current_revision > $2) AS newer
FROM installed_product ip
WHERE ip.id = $1;

-- name: r04_list_installations
-- params: none
-- The portal overview. Configuration count, current revision and last
-- modification come from data this design must aggregate -- the aggregate the
-- rollup designs (n3, n4) buy back.
-- No aggregate and no GROUP BY: the four numbers are already on the parent row.
SELECT ip.id,
       pd.name AS product_definition,
       e.name  AS environment,
       bu.name AS business_unit,
       ip.config_count,
       ip.current_revision,
       ip.updated_at AS last_modified
FROM installed_product ip
JOIN product_definition pd ON pd.id = ip.product_definition_id
JOIN environment e         ON e.id  = ip.environment_id
LEFT JOIN business_unit bu ON bu.id = ip.business_unit_id
ORDER BY ip.id;

-- name: r05_search_key
-- params: scope_kind, scope_id, key
-- One key across installations, scoped by product definition, environment or
-- business unit. The scope is a parameter rather than three statements, so every
-- design answers the same three questions through one plan shape.
SELECT c.installed_product_id AS installed_product_id, c.value
FROM config_entry c
JOIN installed_product ip ON ip.id = c.installed_product_id
WHERE c.key = $3
  AND (   ($1 = 'product_definition' AND ip.product_definition_id = $2)
       OR ($1 = 'environment'        AND ip.environment_id = $2)
       OR ($1 = 'business_unit'      AND ip.business_unit_id = $2))
ORDER BY c.installed_product_id;
