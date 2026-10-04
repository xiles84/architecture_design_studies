#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#pagebreak(weak: true)
== Configuration publication <case-config>
#marker("case", "Study 04 · installed product → configuration entries")

A portal publishes configuration for installed products. Each installation owns its entries; installations of the same product definition must remain separate.

Readers request the effective configuration, an individual key or the installation overview. Writers edit keys, publish a batch or replace the whole configuration.

#block(inset: 8pt, fill: palette.surface)[
  *Conceptual maintenance path:* change entries → maintain installation summary → commit publication.
  \
  A trigger repeats work as rows change. An application maintainer can collect the batch and maintain its summary together.
]

=== Variation: store entries as rows or as a document

#block[
  #set text(size: 9.5pt)
  #table(columns: (1.2fr, 2fr, 1.3fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Variant], [Representation], [Theory]),
    [n1_rows_indexed], [An indexed row for each key: the reference.], [#link(<theory-normalized>)[Normalized facts]],
    [n2_rows_rolldown], [Parent keys copied onto entry rows.], [#link(<theory-rolldown>)[Rolldown]],
    [d1_doc_row], [One whole configuration document per installation.], [#link(<theory-embedding>)[Embedding]],
  )
]

*Thrives — inferred fit:* entry rows suit independent key access; a document suits fetching the whole configuration together.

*Dies — inferred boundary:* a whole document can amplify small edits; copied parent keys add maintenance paths. These are selection considerations, not numeric wins established by the claims below.

=== Variation: choose who maintains the overview

*Change:* n3_rollup_trigger and n4_rollup_app maintain the same overview facts. The comparison changes whether database triggers or application code do the maintenance.

*Observed summary:* whole-configuration replacement costs 4.0-5.7x more with the trigger. The slowest 1% is 22-46x worse; repeated measurements disagree substantially, so use ranges.
#claim-badge("v2-13-trigger-rollup-cost")

#block[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Replacement size], [Application throughput relative to trigger]),
    [10 entries], [Roughly 1.1-3.2x.],
    [60 entries], [Roughly 2.5-9.0x.],
    [240 entries], [Roughly 10-24x.],
  )
]
#claim-badge("v2-18-maintenance-size-interaction")

*Thrives — inferred fit:* triggers centralise maintenance when multiple write paths must obey it.

*Dies — observed cost:* large replacements multiply the trigger penalty. The application variant is faster in every measured size arrangement; neither variant records a rollup mismatch.

Theory: #link(<theory-rollup>)[Rollup]. The experiment cannot separate trigger firing frequency from the work each firing performs.

=== Variation: many writers edit one hot installation

*Change:* c1_optimistic_version checks a version and retries bounded conflicts. c2_pessimistic_lock acquires the parent lock before reading and updating.

*Observed summary:* under one hot key at 16 writers, lock-before-update delivers about 1.6-1.8x the version-check throughput. Removing the retry budget does not explain that advantage.
#claim-badge("v2-12-lock-vs-version-upper-bound")

*Observed crossing:* the two strategies cross between one and two writers; repeated runs preserve that boundary and the divergence above it. An unexplained one-writer excursion remains recorded.
#claim-badge("v2-17-optimistic-pessimistic-crossover")

*Thrives — observed regime:* version checking leads at one writer in the reproduced crossing; locking leads as the hot-key contention increases.

*Dies — observed correctness:* unchecked read-modify-write loses acknowledged updates. Both guarded designs lose none; correctness precedes the speed comparison.

Theory: #link(<theory-optimistic>)[Optimistic writes], #link(<theory-pessimistic>)[Pessimistic writes] and #link(<theory-hot-rows>)[Hot rows]. The crossing describes this hot-key workload, not every configuration portal.

#gap[
  Eight of the eighteen planned designs are unimplemented. Cardinality, cadence/churn, equal-total resources, placement and YugabyteDB coverage remain incomplete; the active gap card also retains the earlier repeated-randomized-trials gap.
  #claim-badge("v2-gap-03-study04-missing-designs")
]

*Reproduce:* `studies/04-configuration-portal`. Badges name the exact replacement and concurrency experiments. Treat the later repeated claims as targeted coverage, rather than completion of the whole study plan.
