-- Built after the load. A per-row index maintained during the COPY measures the
-- loader, not the design.
CREATE INDEX config_entry_by_key ON config_entry (key, installed_product_id);
