-- The design's own account of the truth, read from the document it maintains.

-- name: a_key_values
-- params: installed_product_id
SELECT e.key, e.value
FROM installed_product ip
CROSS JOIN LATERAL jsonb_each_text(ip.doc) AS e(key, value)
WHERE ip.id = $1
ORDER BY e.key;

-- name: a_revision
-- params: installed_product_id
SELECT current_revision AS revision FROM installed_product WHERE id = $1;

-- name: a_counts
-- params: installed_product_id
SELECT (SELECT count(*) FROM jsonb_object_keys(ip.doc)) AS entries
FROM installed_product ip WHERE ip.id = $1;

-- name: a_document_revision_mismatches
-- The document's own guard revision and the parent's pointer must agree; for d2
-- they are two columns of the same row, so a gap would still mean a publication
-- had been acknowledged without the pointer moving.
SELECT count(*) AS mismatches
FROM installed_product ip
WHERE ip.doc_revision <> ip.current_revision;

-- name: a_content_hash_mismatches
-- INV-9 for this design: the stored hash must equal the hash of the stored
-- document. A hash that has drifted from its own document is a rollup that lies.
SELECT count(*) AS mismatches
FROM installed_product ip
WHERE ip.content_hash <> md5(ip.doc::text);
