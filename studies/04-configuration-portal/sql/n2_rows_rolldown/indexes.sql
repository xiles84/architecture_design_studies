-- Built after the load.
--
-- Three indexes, not one. The copied keys only save the join if the search is
-- indexed on the column it now filters, so a rolldown trades one index on (key)
-- for three scope+key indexes and pays all three on every insert. That cost is
-- the point of the comparison, not an implementation detail.
CREATE INDEX config_entry_by_pd  ON config_entry (product_definition_id, key, installed_product_id);
CREATE INDEX config_entry_by_env ON config_entry (environment_id, key, installed_product_id);
CREATE INDEX config_entry_by_bu  ON config_entry (business_unit_id, key, installed_product_id);
