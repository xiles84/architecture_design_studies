-- Common tables. Identical in every design: this study compares how
-- configuration is stored, not how products are catalogued.
CREATE TABLE product_definition (
  id              int PRIMARY KEY,
  name            text NOT NULL,
  product_version text NOT NULL
);

CREATE TABLE environment (
  id   int PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE business_unit (
  id   int PRIMARY KEY,
  name text NOT NULL
);

-- The whole effective configuration is a JSON column ON THE PARENT ROW. d1 keeps
-- the document in a separate one-to-one table; d2 removes that table, so this
-- design isolates embedding the document in the parent. The expected cost is not
-- the publication (both rewrite the whole document) but the OTHER write: an
-- ordinary metadata update (w06_update_metadata) now rewrites the parent row,
-- including the TOASTed document, where d1's metadata update touches only a short
-- row. The pair d1 <-> d2 therefore isolates the embedding decision.
--
-- `doc_revision` is the document's own optimistic guard (d1 keeps it on the
-- document row); `current_revision` is the parent pointer the rest of the schema
-- uses. The harness keeps them equal in one transaction and the audit checks it,
-- exactly as for d1.
CREATE TABLE installed_product (
  id                    bigint PRIMARY KEY,
  product_definition_id int NOT NULL REFERENCES product_definition(id),
  environment_id        int NOT NULL REFERENCES environment(id),
  business_unit_id      int REFERENCES business_unit(id),
  display_name          text NOT NULL,
  current_revision      bigint NOT NULL DEFAULT 0,
  updated_at            timestamptz NOT NULL DEFAULT now(),
  doc                   jsonb  NOT NULL DEFAULT '{}'::jsonb,
  doc_revision          bigint NOT NULL DEFAULT 0,
  content_hash          text   NOT NULL DEFAULT '',
  bytes                 int    NOT NULL DEFAULT 0
);

-- Staging for the bulk load. The loader COPYs every (installed product, key,
-- value) here and then runs the design's w_load_* statement once, so the design
-- under test does not pay a per-row load cost that the other designs avoid.
CREATE TABLE load_entry (
  installed_product_id bigint NOT NULL,
  key                  text NOT NULL,
  value                text NOT NULL
);
