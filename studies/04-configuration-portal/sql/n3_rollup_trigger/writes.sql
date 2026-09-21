-- The write catalogue. The harness composes these into transactions; which
-- statements share a transaction is decided in harness/designs.go, and the SQL
-- decides only what each statement does.
--
-- Every publication ends with w_revision_bump inside the same transaction, so a
-- change and its revision advance are atomic together (INV-3, INV-12). That is
-- what makes "the client was told it published" and "the database holds it" the
-- same event.

-- name: w_revision_bump
-- params: installed_product_id
UPDATE installed_product
SET current_revision = current_revision + 1, updated_at = now()
WHERE id = $1;

-- name: w01_modify_key
-- params: installed_product_id, key, value
UPDATE config_entry SET value = $3, version = version + 1, updated_at = now()
WHERE installed_product_id = $1 AND key = $2;

-- name: w02_add_key
-- params: installed_product_id, key, value
INSERT INTO config_entry (installed_product_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (installed_product_id, key)
DO UPDATE SET value = EXCLUDED.value, version = config_entry.version + 1, updated_at = now();

-- name: w03_delete_key
-- params: installed_product_id, key
DELETE FROM config_entry WHERE installed_product_id = $1 AND key = $2;

-- name: w04_batch_upsert
-- params: installed_product_id, keys, values
INSERT INTO config_entry (installed_product_id, key, value)
SELECT $1, t.k, t.v FROM unnest($2::text[], $3::text[]) AS t(k, v)
ON CONFLICT (installed_product_id, key)
DO UPDATE SET value = EXCLUDED.value, version = config_entry.version + 1, updated_at = now();

-- name: w04_batch_delete
-- params: installed_product_id, delete_keys
DELETE FROM config_entry WHERE installed_product_id = $1 AND key = ANY($2::text[]);

-- name: w05_replace_clear
-- params: installed_product_id
DELETE FROM config_entry WHERE installed_product_id = $1;

-- name: w05_replace_fill
-- params: installed_product_id, keys, values
INSERT INTO config_entry (installed_product_id, key, value)
SELECT $1, t.k, t.v FROM unnest($2::text[], $3::text[]) AS t(k, v);

-- name: w06_update_metadata
-- params: installed_product_id, display_name
-- Ordinary metadata, deliberately not a publication: it moves no revision, and on
-- a design that embeds configuration in this row it must not rewrite it either.
UPDATE installed_product SET display_name = $2, updated_at = now() WHERE id = $1;

-- name: w_load_entries
INSERT INTO config_entry (installed_product_id, key, value)
SELECT installed_product_id, key, value FROM load_entry;

-- There is deliberately no rollup statement here. The trigger below owns the
-- rollup, so the harness cannot forget it -- and cannot choose to skip it either,
-- which is exactly the trade being measured.
