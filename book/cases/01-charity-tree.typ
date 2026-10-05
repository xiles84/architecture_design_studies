#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#pagebreak(weak: true)
== Charity donations <case-charity>
#marker("case", "Study 01 · charity → person → donation")

A charity site lists donors, their donations and the latest activity. Its dashboard repeatedly asks for totals; donation entry changes the facts underneath.

*Keep fixed:* the same donors, donations and answers. Change where information lives and how each question reaches it.

#figure-evidence(
  "../assets/pilot-00-overview.svg", "structure", "implemented design contract",
  "One business tree, several ways to serve it.",
  [Donations belong to people, who belong to charities. Variants change joins, stored totals or child storage. The overview's pair labels describe design intent; claim limits qualify confounding and unverified placement.],
  image-width: 82%,
)

=== Variation: find donations without copying everything

#block(breakable: false)[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2fr, 1.5fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Variant], [What changes], [Theory]),
    [D1 minimal], [Keep the normalized tree with minimal indexes.], [#link(<theory-normalized>)[Normalized facts]],
    [D2 indexed], [Add the query index; keep the tree.], [#link(<theory-indexes>)[Indexes]],
    [D3 flat+FK], [Copy the charity key onto donations; retain the foreign key.], [#link(<theory-rolldown>)[Rolldown]],
  )
]

*Observed:* D2 scores about 10-11x D1 on the blended read score; D3 reaches about 14-16x in two whole runs. Carry the ordering and bands, rather than the first run's point figures.
#claim-badge("v2-01-normalized-index-first")

*Observed:* D3's advantage over D2 survives fresh loads and reversed whole-run order. The repeated band is 1.59-1.84x; its insert-cost ratio remains imprecise.
#claim-badge("v2-05-d3-read-advantage-repeated")

*Thrives — inferred fit:* repeatable read paths that the index or copied key can answer directly.

*Dies — inferred boundary:* parent reassignment becomes costly if a copied key must fan out across many children. This is a maintenance risk, not a measured crossover here.

=== Variation: answer the dashboard from stored totals

*Change:* D4 uses triggers and D5 uses application code to maintain totals on parents. Readers no longer aggregate all contributing donations.

*Observed summary:* D4 scores about 61-85x D1 on aggregate reads; D5 about 64-83x. Inserts run at about a third to two-fifths of D1's rate.
#claim-badge("v2-02-rollup-family")

*Thrives — inferred fit:* frequently viewed totals with relatively few contributing writes.

*Dies — observed cost:* every insert pays maintenance. A write-heavy workload must justify that penalty; the study does not establish a universal read/write threshold.

Theory: #link(<theory-rollup>)[Rollup]. For a replacement-size effect, continue to #link(<case-config>)[Configuration portal].

=== Variation: put donations inside the donor

*Change:* D6 stores donation children in a JSON/array column on the donor.

*Observed summary:* reads score about 1.7-2.2x D1, below both indexed designs. Inserts run at about 0.57-0.60x D1's rate.
#claim-badge("v2-03-embedding-not-ahead")

*Dies — observed:* this embedding package adds write cost without beating the indexed read alternatives.

*Thrives — not established:* whole-document access may fit other workloads. That expectation is not a win measured in this charity case.

Theory: #link(<theory-embedding>)[Embedding].

=== Variation: ask who donated last in a period

*Define the answer first:* a donor's latest donation falling inside a period differs from their latest donation *within* that period. One filters a global latest value; the other computes a period-local latest value.

*Observed summary:* the rollup-column index beats the recency flag by 1.9-2.3x on q13 and 2.6-3.1x on q15 trailing. It is level on q14 and slightly slower on q16.
#claim-badge("v2-04-recency-maintained-index")

*Thrives or dies — observed:* the winning path depends on the actual question. A group headline hides q16's opposite direction.

Theory: #link(<theory-select-contract>)[Answer contracts], #link(<theory-select-aggregate>)[Aggregation] and #link(<theory-select-pagination>)[Pagination].

#pagebreak(weak: true)
=== Variation: join donations and filter them <case-charity-join-plan>

*Business question:* return each qualifying donation with its donor, one row per donation. The theory's customers and orders map to people and donations. The measured table holds 200,000 donations; the filter keeps 0.01% to 90% of them.

*Observed — the six forms:* the six #link(<theory-select-join-filter>)[join and filter forms] ran the same plan in all 120 engine, phase and selectivity pairs. No form is faster. Write the one that reads best. \
#claim-badge("v6-sel-01-inner-join-spellings-same-plan")

*Observed — a CTE fence:* forcing MATERIALIZED on a CTE that holds the join, with the filter outside, cost 84–101x on PostgreSQL and 34–49x on YugabyteDB at 0.01%. At 90% it cost 1.02–1.44x and 0.93–1.03x. With the filter inside the CTE, forcing it cost at most 1.61x on PostgreSQL. \
#claim-badge("v6-sel-02-materialized-joined-cte-penalty")

*Observed — optional donations:* LEFT JOIN with the donation filter in ON kept every donor. The same filter in WHERE dropped donors without a qualifying donation; that wrong answer was detected in all 12 runs. It is a different question, not a faster spelling. \
#claim-badge("v6-sel-01-inner-join-spellings-same-plan")

*Thrives — observed fit:* any readable join form with the engine's default CTE policy, or a filter written inside the CTE.

*Dies — observed:* a #link(<theory-select-cte>)[materialization hint] that keeps a selective outer filter away from the join.

=== Variation: reuse one expensive CTE <case-charity-cte-reuse>

*Business question:* several consumers read the same expensive per-donation token. Compute it once, or inline it into every consumer? One consumer selects 1% of donations; the other reads all of them.

*Observed — the crossing:* inlining (NOT MATERIALIZED) against MATERIALIZED crossed between one and two references, on PostgreSQL with the unselective consumer and on YugabyteDB with both. At two references inlining took 1.43–2.06x the materialized time; at sixteen, 3.9–12.1x. \
#claim-badge("v6-sel-03-cte-reuse-crossing")

*Observed — the exception:* on PostgreSQL with the selective consumer, inlining stayed faster at every count from 1 to 16: 0.17–0.27x at one reference, 0.68–0.81x at sixteen. Both engines' defaults materialize from two references, so forcing NOT MATERIALIZED was 1.2–4.8x faster there. \
#claim-badge("v6-sel-04-cte-reuse-postgresql-selective-no-crossing")

*Thrives or dies — observed:* the default policy fits three of the four measured regimes. An inlining hint helps only the selective PostgreSQL consumer. The crossing is bracketed between one and two references, not located.

#pagebreak(weak: true)
=== Variation: who has, or has not, donated <case-charity-membership>

*Business question:* list donors with a qualifying donation, or donors with none, each donor once. The generated SELECT fixture of Study 01 stands in for donors and donations: 10,000 donations over 10,000, 1,000 or 100 donors. That gives fanout 1, 10 and 100.

*Observed — membership:* JOIN plus DISTINCT took 1.40–1.69x the EXISTS time at fanout 1 and 3.64–4.23x at fanout 100 on PostgreSQL. On YugabyteDB it was 1.10–1.16x at fanout 1 and 0.95–1.01x above it. Fanout and donor count change together here. \
#claim-badge("v6-sel-05-membership-join-distinct-penalty")

*Observed — IN or EXISTS:* neither engine showed a resolved difference; PostgreSQL ranged 0.86–1.30 with mixed signs. Choose the form that reads better. \
#claim-badge("v6-sel-06-membership-in-exists-no-difference")

*Observed — absence:* NOT IN over a donation key containing a NULL returned no donors in all 18 checks. NOT EXISTS, LEFT JOIN with IS NULL, and NOT IN over a NOT NULL subquery returned the right donors. \
#claim-badge("v6-sel-07-nullable-not-in-empty")

*Observed — anti-join cost:* at fanout 100 on PostgreSQL, NOT EXISTS took 0.21–0.24 ms. LEFT JOIN with IS NULL was 2.71–3.27x slower and NOT-NULL NOT IN 2.85–3.16x. On YugabyteDB the three stayed within 0.997–1.097x. \
#claim-badge("v6-sel-08-anti-join-not-exists-postgresql")

*Thrives or dies — observed:* #link(<theory-select-membership>)[EXISTS and NOT EXISTS] fit both questions. JOIN plus DISTINCT dies at high fanout on PostgreSQL. A NOT IN that can meet a NULL always dies.

=== Variation: donor totals and each donor's top donation <case-charity-aggregate>

*Business question:* one total per requested donor, or each donor's largest donation. The request asks for 1 to 10,000 donors on the same generated fixture, two donations per donor.

*Observed — totals on PostgreSQL:* a correlated SUM per donor took 0.47–0.90x the join-then-GROUP-BY time for one donor. It took 1.26–1.72x at 100 donors and 2.49–2.70x at 10,000. \
#claim-badge("v6-sel-09-aggregate-correlated-postgresql-crossover")

*Observed — totals on YugabyteDB:* the correlated SUM never won. It took 1.07–1.23x at one donor and 398–541x at 10,000: 6.7–20.0 s per query, against 16–17 ms on PostgreSQL. \
#claim-badge("v6-sel-10-aggregate-yugabytedb-no-reversal")

*Observed — the top donation:* a LATERAL probe per donor beat `ROW_NUMBER()` only for one donor on YugabyteDB, at 0.54–0.70x. It lost from 10 donors, reaching 257–289x at 10,000. On PostgreSQL it lost from 100 donors. \
#claim-badge("v6-sel-11-topn-lateral-probe-crossover")

*Observed — DISTINCT ON:* on PostgreSQL it took 0.61–0.64x the `ROW_NUMBER()` time at 10,000 donors. At 10 and 100 donors the sign was mixed. It answers only the single top row. \
#claim-badge("v6-sel-12-topn-distinct-on-postgresql")

*Thrives or dies — observed:* #link(<theory-select-aggregate>)[per-donor probes] fit a handful of donors. Join-then-aggregate and `ROW_NUMBER()` with a tiebreak fit large requests. A double-counting join and a tiebreak-free `ROW_NUMBER()` were detected as wrong answers in every run.

#pagebreak(weak: true)
=== Variation: page through donations, or match two ranges <case-charity-pagination>

*Business question:* show the next five donations in amount order, or the donations inside either of two amount ranges. The generated fixture holds 12,000 donations.

*Observed — pages:* on PostgreSQL, keyset took 1.06–1.29x the OFFSET time at depths 1–100 and 0.21–0.24x at depth 10,000. Both returned the same page at every depth under a fixed snapshot. The crossing lies between depth 100 and 1,000. \
#claim-badge("v6-sel-15-keyset-offset-crossover")

*Limit — the cursor:* this keyset form had no composite (amount, id) index. The flat keyset cost that #link(<theory-select-pagination>)[the theory] expects with a matching index was not tested.

*Observed — OR or UNION:* rewriting the two-range OR as UNION ALL or UNION never beat OR. On PostgreSQL's broad ranges they took 1.23–1.55x and 1.25–1.50x the OR time. \
#claim-badge("v6-sel-13-or-not-beaten-by-union")

*Observed — an expression index:* with the query unchanged, an index matching #box[`(amount + 0)`] cut PostgreSQL's selective expression OR to 0.28–0.33x. The plan switched to that index. \
#claim-badge("v6-sel-14-expression-index-postgresql")

*Thrives or dies — observed:* OFFSET fits shallow pages and keyset deep ones. OR fits a disjunction, with an index that matches the expression. UNION ALL over overlapping ranges returned duplicates; the control caught it in every run.

#gap[
  YugabyteDB's membership, absence, pagination and range figures come from scan-only plans: its default hash-sharded indexes could not serve these predicates. These three families also lack a cluster arm. \
  #claim-badge("v6-sel-gap-01-yugabytedb-scan-only-and-single-node-families")
]

=== Variation: place the joined donations on three nodes <case-charity-placement>

*Business question:* on a three-node YugabyteDB cluster, should the join's tables share one tablet, keep one tablet each, or split into three? Tablets and leaders were checked before and after every run.

*Observed — splitting:* three tablets with leaders on three nodes made a 2,000-row join 0.73x and a 180,000-row join 0.54x the single-tablet layout (isolated phase). The 20-row join was 1.22x slower. \
#claim-badge("v6-sel-16-placement-splitting-effect")

*Observed — colocation:* a colocated database took 0.71x the split layout's time at a 20-row join and 0.78x at 2,000 rows. At 180,000 rows it was 1.21x slower, and 1.54x with a concurrent writer. \
#claim-badge("v6-sel-17-placement-colocated-vs-distributed")

*Thrives or dies — observed:* split data fits broad reads; colocation fits point and medium reads. All nodes are containers on one laptop, so this is data placement without a network.

Theory: #link(<theory-placement>)[Colocation and distributed access].

#pagebreak(weak: true)
=== Choosing a SELECT template <case-charity-select-guide>

Start from the answer, then the readable form. Switch only where the observed boundary applies.

#let ids(..xs) = xs.pos().map(x => box(raw(x))).join(", ")
#block(breakable: false)[
  #set text(size: 9.5pt)
  #set par(justify: false)
  #table(columns: (1.1fr, 1fr, 2fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Required answer], [Start with], [Switch or avoid — observed]),
    [Donations with a filter], [Any readable join form], [Avoid MATERIALIZED on a CTE with a selective outer filter. #ids("v6-sel-01", "v6-sel-02")],
    [A CTE used several times], [The engine's default], [Inline only for a selective PostgreSQL consumer. #ids("v6-sel-03", "v6-sel-04")],
    [Donors with optional donations], [LEFT JOIN, filter in ON], [Never move the child filter to WHERE. #ids("v6-sel-01")],
    [Donors who donated], [EXISTS or IN], [Avoid JOIN plus DISTINCT at high fanout on PostgreSQL. #ids("v6-sel-05", "v6-sel-06")],
    [Donors who never donated], [NOT EXISTS], [Never NOT IN over a nullable key. #ids("v6-sel-07", "v6-sel-08")],
    [One total per donor], [Join, then GROUP BY], [Correlated SUM only for very few donors, on PostgreSQL. #ids("v6-sel-09", "v6-sel-10")],
    [Top donation per donor], [`ROW_NUMBER()` with a tiebreak], [LATERAL for a single donor; DISTINCT ON at scale on PostgreSQL. #ids("v6-sel-11", "v6-sel-12")],
    [The next page], [OFFSET for shallow pages], [Keyset past about depth 1,000 on PostgreSQL. #ids("v6-sel-15")],
    [Either of two ranges], [OR], [UNION ALL never won; index the expression you filter on. #ids("v6-sel-13", "v6-sel-14")],
  )
]

The full claim behind each id is linked from the badges above. Every row comes from one laptop, one sequential reader and named engines; no row is a universal ranking.

#gap[
  Native document, key-value, graph and wide-column engines were not measured. SQL JSON storage cannot establish their performance.
  #claim-badge("v2-gap-06-native-datastore-families")
]

*Reproduce:* `studies/01-charity-tree`; the SELECT experiment starts at `studies/01-charity-tree/protocols/select-query-templates.md`. Each badge links to the exact claims, reports, run tags and limits. The blended score combines queries; it is not the speed of every query.
