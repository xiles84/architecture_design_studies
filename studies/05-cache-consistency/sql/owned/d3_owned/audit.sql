-- Owned audit catalogue.
--
-- The legacy audits, unchanged, plus the two invariants only a version token and
-- an outbox can express. If the owned model cannot show these two holding, then
-- the token it sells to the cache is not a token at all and no strict owned cell
-- may be believed.

-- name: a_person_aggregate
-- params: person_id
SELECT count(*)                                  AS donation_count,
       COALESCE(sum(amount_cents), 0)            AS donation_total_cents
  FROM donation
 WHERE person_id = $1;

-- name: a_donation_ids
-- params: person_id
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donation_id;

-- name: a_recent_ids
-- params: person_id
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC, donation_id DESC
 LIMIT 20;

-- name: a_rolldown_drift
-- params: none
SELECT count(*)
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id <> p.charity_id;

-- name: a_version_monotonic
-- params: person_id
-- Counts the ways this person's version and its change log can disagree:
--   * the version is smaller than the number of committed outbox events -- the
--     token was rolled back while the event survived, or the bump was skipped;
--   * the version is greater than the event count -- a bump with no event, which
--     is precisely a cache invalidation with no record of why.
-- It must be zero. This is the audit that makes "the version commits atomically
-- with the mutation" a measured claim rather than a comment.
SELECT count(*)
  FROM (
    SELECT $1::bigint AS person_id,
           (SELECT cache_version FROM person WHERE person_id = $1::bigint) AS version,
           (SELECT count(*) FROM cache_outbox o WHERE o.person_id = $1::bigint) AS events
  ) x
 WHERE x.version <> x.events;

-- name: a_outbox_orphans
-- params: none
-- Counts outbox events naming a person that does not exist. A transactional outbox
-- cannot produce one; a bug that writes the event outside the transaction can. It
-- must be zero.
SELECT count(*)
  FROM cache_outbox o
 WHERE NOT EXISTS (SELECT 1 FROM person p WHERE p.person_id = o.person_id);
