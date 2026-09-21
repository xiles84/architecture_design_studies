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

-- One row per configuration key, keyed by the installed product. The primary key
-- is the lookup path for a complete configuration; the secondary index on (key)
-- exists only to answer the cross-installation search r05, which is precisely
-- what the index-control pair removes.
--
-- `version` is carried by every row-per-key design and bumped on every write.
-- On the reference and the rollup designs it is bookkeeping, not arbitration:
-- these designs write an explicit value, so last-writer-wins is correct and no
-- guard is needed. c1 is what uses it as a guard, and x1 is what happens when
-- the read-modify-write it protects has no guard at all.
CREATE TABLE config_entry (
  installed_product_id bigint NOT NULL REFERENCES installed_product(id) ON DELETE CASCADE,
  key                  text NOT NULL,
  value                text NOT NULL,
  version              bigint NOT NULL DEFAULT 0,
  updated_at           timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (installed_product_id, key)
);

-- The parent rollups, maintained by the application rather than by a trigger.
-- The columns are identical to n3_rollup_trigger's on purpose: rollup-trigger and
-- rollup-app are a controlled pair, and the only decision between them is who
-- does the maintenance. count, revision and last modification are the three
-- numbers the portal overview (r04) needs; content_hash is the fourth, and the
-- one that is expensive: it can only be recomputed by reading every entry.
ALTER TABLE installed_product ADD COLUMN config_count int  NOT NULL DEFAULT 0;
ALTER TABLE installed_product ADD COLUMN content_hash text NOT NULL DEFAULT '';
