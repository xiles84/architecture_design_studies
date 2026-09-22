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

CREATE TABLE installed_product (
  id                    bigint PRIMARY KEY,
  product_definition_id int NOT NULL REFERENCES product_definition(id),
  environment_id        int NOT NULL REFERENCES environment(id),
  business_unit_id      int REFERENCES business_unit(id),
  display_name          text NOT NULL,
  current_revision      bigint NOT NULL DEFAULT 0,
  updated_at            timestamptz NOT NULL DEFAULT now()
);

-- Staging for the bulk load. The loader COPYs every (installed product, key,
-- value) here and then runs the design's w_load_* statement once, so the design
-- under test does not pay a per-row load cost that the other designs avoid.
CREATE TABLE load_entry (
  installed_product_id bigint NOT NULL,
  key                  text NOT NULL,
  value                text NOT NULL
);
-- One row per installed product, holding the whole effective configuration as a
-- JSON document. A publication replaces the entire document, so the write cost
-- grows with the number of keys while the read cost does not: this is the other
-- end of the axis n1_rows_indexed sits at.
--
-- `revision` lives here rather than only on the parent because it is also this
-- design's optimistic guard: two publications of the same document must not
-- overwrite each other silently. The harness keeps it equal to
-- installed_product.current_revision in the same transaction, and the audit
-- checks that they agree.
CREATE TABLE config_document (
  installed_product_id bigint PRIMARY KEY REFERENCES installed_product(id) ON DELETE CASCADE,
  revision             bigint NOT NULL DEFAULT 0,
  doc                  jsonb  NOT NULL DEFAULT '{}'::jsonb,
  content_hash         text   NOT NULL DEFAULT '',
  bytes                int    NOT NULL DEFAULT 0
);
