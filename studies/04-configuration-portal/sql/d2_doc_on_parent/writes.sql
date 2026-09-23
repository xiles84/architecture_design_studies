-- The write catalogue. The document is read, changed in the application, and
-- written back guarded by `doc_revision`, so a racing publication is rejected
-- rather than silently overwriting (INV-3) -- the same guarantee d1 obtains,
-- on the same row as the rest of the parent's state.

-- name: w_revision_bump
-- params: installed_product_id
UPDATE installed_product
SET current_revision = current_revision + 1, updated_at = now()
WHERE id = $1;

-- name: wd1_read_doc
-- params: installed_product_id
SELECT doc, doc_revision FROM installed_product WHERE id = $1;

-- name: wd1_write_doc
-- params: installed_product_id, doc, bytes, expected_revision
-- Whole-document replacement on the parent row. The statement is the same shape
-- as d1's; the difference is that the parent's short columns live in the same
-- tuple, so a metadata-only update (w06) pays for this document's TOASTed bytes
-- too. That amplification is the decision d2 exists to measure.
UPDATE installed_product
SET doc = $2::jsonb, doc_revision = $4 + 1, content_hash = md5(($2::jsonb)::text), bytes = $3
WHERE id = $1 AND doc_revision = $4;

-- name: w06_update_metadata
-- params: installed_product_id, display_name
-- An ordinary metadata update. In d1 this touches a short row; here it rewrites
-- the parent row, including the document, which is the cost being isolated.
UPDATE installed_product SET display_name = $2, updated_at = now() WHERE id = $1;

-- name: w_load_documents
-- params: none
UPDATE installed_product ip
SET doc = s.doc, bytes = length(s.doc::text), content_hash = md5(s.doc::text)
FROM (
  SELECT le.installed_product_id AS ip, jsonb_object_agg(le.key, le.value) AS doc
  FROM load_entry le
  GROUP BY le.installed_product_id
) s
WHERE ip.id = s.ip;
