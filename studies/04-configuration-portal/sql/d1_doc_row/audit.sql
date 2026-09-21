-- The design's own account of the truth, read from the document it maintains.

-- name: a_key_values
-- params: installed_product_id
SELECT e.key, e.value
FROM config_document d
CROSS JOIN LATERAL jsonb_each_text(d.doc) AS e(key, value)
WHERE d.installed_product_id = $1
ORDER BY e.key;

-- name: a_revision
-- params: installed_product_id
SELECT current_revision AS revision FROM installed_product WHERE id = $1;

-- name: a_counts
-- params: installed_product_id
SELECT (SELECT count(*) FROM jsonb_object_keys(d.doc)) AS entries
FROM config_document d WHERE d.installed_product_id = $1;

-- name: a_document_revision_mismatches
-- The document's own guard revision and the parent's must agree: the harness keeps
-- them equal in one transaction, and a gap would mean a publication had been
-- acknowledged without the pointer moving.
SELECT count(*) AS mismatches
FROM config_document d
JOIN installed_product ip ON ip.id = d.installed_product_id
WHERE d.revision <> ip.current_revision;

-- name: a_content_hash_mismatches
-- INV-9 for this design: the stored hash must equal the hash of the stored
-- document. A hash that has drifted from its own document is a rollup that lies.
SELECT count(*) AS mismatches
FROM config_document d
WHERE d.content_hash <> md5(d.doc::text);
