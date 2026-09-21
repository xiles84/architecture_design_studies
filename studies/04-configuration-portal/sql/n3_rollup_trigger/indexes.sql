-- Built after the load, trigger included.
--
-- The trigger is attached here rather than in schema.sql on purpose: study 01
-- learned that loading with per-row triggers measures the loader, not the design.
-- The trigger's cost is what the *write* benchmark is for.
CREATE INDEX config_entry_by_key ON config_entry (key, installed_product_id);

-- One trigger function, recomputing the whole rollup for the affected installed
-- product. It fires per row, so a batch of fifty changes recomputes fifty times;
-- that is the honest cost of maintaining the aggregate in the database, and it is
-- what n4 is compared against.
CREATE OR REPLACE FUNCTION refresh_ip_rollup() RETURNS trigger AS $$
DECLARE target bigint;
BEGIN
  target := COALESCE(NEW.installed_product_id, OLD.installed_product_id);
  UPDATE installed_product ip SET
    config_count = (SELECT count(*) FROM config_entry c WHERE c.installed_product_id = target),
    content_hash = COALESCE((SELECT md5(string_agg(c.key || '=' || c.value, ',' ORDER BY c.key))
                             FROM config_entry c WHERE c.installed_product_id = target), ''),
    updated_at   = now()
  WHERE ip.id = target;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER config_entry_rollup
AFTER INSERT OR UPDATE OR DELETE ON config_entry
FOR EACH ROW EXECUTE FUNCTION refresh_ip_rollup();
