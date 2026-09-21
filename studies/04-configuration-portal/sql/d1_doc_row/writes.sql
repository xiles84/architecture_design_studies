-- The write catalogue. The document is read, changed in the application, and
-- written back guarded by its revision, so a racing publication is rejected
-- rather than silently overwriting (INV-3) -- the same guarantee c1 obtains with
-- a row version, obtained here with a document version.

-- name: w_revision_bump
-- params: installed_product_id
UPDATE installed_product
SET current_revision = current_revision + 1, updated_at = now()
WHERE id = $1;

-- name: wd1_read_doc
-- params: installed_product_id
SELECT doc, revision FROM config_document WHERE installed_product_id = $1;

-- name: wd1_write_doc
-- params: installed_product_id, doc, bytes, expected_revision
-- Whole-document replacement. The cost that grows with cardinality is here: every
-- publication rewrites every key, and a large document is TOASTed and rewritten
-- as a unit.
--
-- The content hash is computed by the database from the document it is about to
-- store rather than sent by the application: jsonb normalises key order, so a
-- hash computed in Go over the marshalled text would disagree with the stored
-- document for reasons that have nothing to do with this study. The audit that
-- compares the two columns therefore checks that the statement is internally
-- consistent; the study's correctness gate is the ledger audit, which compares
-- the whole document against what the client was told.
UPDATE config_document
SET doc = $2::jsonb, revision = $4 + 1, content_hash = md5(($2::jsonb)::text), bytes = $3
WHERE installed_product_id = $1 AND revision = $4;

-- name: w06_update_metadata
-- params: installed_product_id, display_name
-- The reason this design keeps the document on its own row rather than on the
-- parent: ordinary metadata is a configuration-free update.
UPDATE installed_product SET display_name = $2, updated_at = now() WHERE id = $1;

-- name: w_load_documents
INSERT INTO config_document (installed_product_id, revision, doc, content_hash, bytes)
SELECT s.ip, 0, s.doc, md5(s.doc::text), length(s.doc::text)
FROM (
  SELECT le.installed_product_id AS ip, jsonb_object_agg(le.key, le.value) AS doc
  FROM load_entry le
  GROUP BY le.installed_product_id
) s;
