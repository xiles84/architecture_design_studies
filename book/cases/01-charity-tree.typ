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

Theory: #link(<theory-select-contract>)[Answer contracts], #link(<theory-select-aggregate>)[Aggregation] and #link(<theory-select-pagination>)[Pagination]. SELECT-family timing promotion remains separate; no new SELECT ratios are asserted here.

#gap[
  Native document, key-value, graph and wide-column engines were not measured. SQL JSON storage cannot establish their performance.
  #claim-badge("v2-gap-06-native-datastore-families")
]

*Reproduce:* `studies/01-charity-tree`. Each badge links to the exact claims, reports, run tags and limits. The blended score combines queries; it is not the speed of every query.
