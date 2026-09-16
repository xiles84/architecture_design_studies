-- Query catalogue for D6 (embedded-jsonb).
--
-- Each query is written the *best* way the embedded design allows, not the
-- naive way, so that the comparison is fair. Two exploits are used repeatedly:
--   * the array is appended in chronological order, so donations->0 is the
--     first donation and donations->-1 is the last, with no unnesting;
--   * WITH ORDINALITY lets "most recent N" be an array-position slice rather
--     than a timestamp sort.
-- Where a query still has to unnest every array, that is a genuine property of
-- the design, not a strawman.

-- name: q01_last_donation_global
-- params: none
-- One element inspected per person, but every person row must still be read.
SELECT (e ->> 'i')::BIGINT        AS donation_id,
       p.person_id,
       (e ->> 'a')::BIGINT        AS amount_cents,
       e ->> 'c'                  AS currency,
       (e ->> 't')::TIMESTAMPTZ   AS donated_at
  FROM person p
 CROSS JOIN LATERAL (SELECT p.donations -> -1 AS e) x
 WHERE jsonb_array_length(p.donations) > 0
 ORDER BY (e ->> 't')::TIMESTAMPTZ DESC
 LIMIT 1;

-- name: q02_last_donation_charity
-- params: charity_id
SELECT (e ->> 'i')::BIGINT        AS donation_id,
       p.person_id,
       (e ->> 'a')::BIGINT        AS amount_cents,
       e ->> 'c'                  AS currency,
       (e ->> 't')::TIMESTAMPTZ   AS donated_at
  FROM person p
 CROSS JOIN LATERAL (SELECT p.donations -> -1 AS e) x
 WHERE p.charity_id = $1
   AND jsonb_array_length(p.donations) > 0
 ORDER BY (e ->> 't')::TIMESTAMPTZ DESC
 LIMIT 1;

-- name: q03_top_donor_charity
-- params: charity_id
-- Every array in the charity must be fully unnested and summed. This is the
-- shape of query the embedded design is worst at.
SELECT p.person_id, p.full_name, t.total_cents
  FROM person p
 CROSS JOIN LATERAL (
        SELECT COALESCE(SUM((e ->> 'a')::BIGINT), 0) AS total_cents
          FROM jsonb_array_elements(p.donations) e
      ) t
 WHERE p.charity_id = $1
 ORDER BY t.total_cents DESC
 LIMIT 1;

-- name: q04_top_donors_leaderboard
-- params: charity_id
SELECT p.person_id, p.full_name, t.total_cents
  FROM person p
 CROSS JOIN LATERAL (
        SELECT COALESCE(SUM((e ->> 'a')::BIGINT), 0) AS total_cents
          FROM jsonb_array_elements(p.donations) e
      ) t
 WHERE p.charity_id = $1
 ORDER BY t.total_cents DESC
 LIMIT 10;

-- name: q05_last_donor_charity
-- params: charity_id
SELECT p.person_id, p.full_name, (e ->> 't')::TIMESTAMPTZ AS donated_at
  FROM person p
 CROSS JOIN LATERAL (SELECT p.donations -> -1 AS e) x
 WHERE p.charity_id = $1
   AND jsonb_array_length(p.donations) > 0
 ORDER BY (e ->> 't')::TIMESTAMPTZ DESC
 LIMIT 1;

-- name: q06_person_first_last
-- params: person_id
-- The embedded design's best case: no unnesting, no second table, no index
-- lookup. Read one row, index into both ends of the array.
SELECT (donations -> 0  ->> 't')::TIMESTAMPTZ AS first_at,
       (donations -> -1 ->> 't')::TIMESTAMPTZ AS last_at
  FROM person
 WHERE person_id = $1;

-- name: q07_total_donated_global
-- params: none
-- Full unnest of every donation in the database.
SELECT COALESCE(SUM((e ->> 'a')::BIGINT), 0) AS total_cents
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.donations) e;

-- name: q08_total_donated_charity
-- params: charity_id
SELECT COALESCE(SUM((e ->> 'a')::BIGINT), 0) AS total_cents
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.donations) e
 WHERE p.charity_id = $1;

-- name: q09_person_recent_donations
-- params: person_id
-- Array-position slice: no timestamp parsing, no sort, one row read.
SELECT e.value ->> 'i'                   AS donation_id,
       (e.value ->> 'a')::BIGINT         AS amount_cents,
       e.value ->> 'c'                   AS currency,
       (e.value ->> 't')::TIMESTAMPTZ    AS donated_at,
       e.value ->> 'n'                   AS note
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.donations) WITH ORDINALITY AS e(value, ord)
 WHERE p.person_id = $1
 ORDER BY e.ord DESC
 LIMIT 20;

-- name: q10_donation_by_id
-- params: donation_id
-- Identity lookup is the embedded design's structural weak point: donation_id
-- is buried inside a document. The GIN index turns this from a sequential scan
-- into a containment probe, but it still requires a recheck and an unnest of
-- the matching person's whole array.
SELECT (e ->> 'i')::BIGINT       AS donation_id,
       p.person_id,
       (e ->> 'a')::BIGINT       AS amount_cents,
       e ->> 'c'                 AS currency,
       (e ->> 't')::TIMESTAMPTZ  AS donated_at,
       e ->> 'n'                 AS note
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.donations) e
 WHERE p.donations @> jsonb_build_array(jsonb_build_object('i', $1::BIGINT))
   AND (e ->> 'i')::BIGINT = $1::BIGINT;

-- name: q11_person_donation_count
-- params: person_id
-- Stored in the document header by JSONB itself; no scan at all.
SELECT jsonb_array_length(donations) AS donation_count
  FROM person
 WHERE person_id = $1;

-- name: q12_charity_recent_feed
-- params: charity_id
-- A merge of many per-person arrays. There is no index that can order across
-- documents, so this unnests the charity and sorts.
SELECT (e ->> 'i')::BIGINT       AS donation_id,
       (e ->> 'a')::BIGINT       AS amount_cents,
       (e ->> 't')::TIMESTAMPTZ  AS donated_at,
       p.full_name
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.donations) e
 WHERE p.charity_id = $1
 ORDER BY (e ->> 't')::TIMESTAMPTZ DESC
 LIMIT 50;

-- ---------------------------------------------------------------------------
-- Recency-window questions (study 01 v4, RECENCY.md). See d1's queries.sql
-- for the full explanation of the two window regimes.
--
-- D6 has no separate donation table to aggregate: the whole history lives in
-- person.donations. Answering "who last gave in this window" means unnesting
-- every donor's ENTIRE document to find its maximum timestamp -- this is the
-- honest cost of taking embedding to its conclusion, and it is measured here
-- rather than excused. There is no useful index for this: an expression index
-- on the array's timestamp text is either impossible (a STABLE cast is
-- rejected by PostgreSQL for an index) or unsound (indexing raw text makes
-- correctness depend on JSON rendering). See RECENCY.md section 3.
-- ---------------------------------------------------------------------------

-- name: q13_donors_last_gift_window
-- params: since, until
SELECT p.person_id, p.full_name, l.last_at
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX((e->>'t')::timestamptz) AS last_at
                       FROM jsonb_array_elements(p.donations) e) l
 WHERE l.last_at >= $1 AND l.last_at < $2
 ORDER BY l.last_at DESC
 LIMIT 100;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX((e->>'t')::timestamptz) AS last_at
                       FROM jsonb_array_elements(p.donations) e) l
 WHERE l.last_at >= $1 AND l.last_at < $2;

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT p.person_id, p.full_name, l.last_at
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX((e->>'t')::timestamptz) AS last_at
                       FROM jsonb_array_elements(p.donations) e) l
 WHERE p.charity_id = $1 AND l.last_at >= $2 AND l.last_at < $3
 ORDER BY l.last_at DESC
 LIMIT 100;

-- name: q16_lapsed_donors_count
-- params: since
SELECT COUNT(*) AS donor_count
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX((e->>'t')::timestamptz) AS last_at
                       FROM jsonb_array_elements(p.donations) e) l
 WHERE l.last_at < $1;
