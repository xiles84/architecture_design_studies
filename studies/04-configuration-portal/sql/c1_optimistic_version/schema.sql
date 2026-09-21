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

-- The concurrency family measures arbitration, so its "modify" operation is a
-- read-modify-write whose new value is *derived from the value it read* (the
-- harness increments a counter carried in the configuration). That is what makes
-- a lost update observable at all: a design that writes an explicit value cannot
-- lose one, it can only be overwritten, which is correct behaviour.

-- This design reads the row's version and writes back only if the version has
-- not moved. Zero rows affected means someone else published first, and the
-- harness retries within its deadline. The version column already exists on
-- config_entry in every row-per-key design; here it is used as a guard rather
-- than as bookkeeping.
