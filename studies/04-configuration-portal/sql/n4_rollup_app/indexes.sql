-- Built after the load.
--
-- No trigger. The application recomputes the rollup with w_recompute_rollup in
-- the same transaction as the change, which is the other half of the pair.
CREATE INDEX config_entry_by_key ON config_entry (key, installed_product_id);
