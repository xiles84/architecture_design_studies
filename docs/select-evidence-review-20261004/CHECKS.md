# SELECT independent HIGH evidence check

Generated facts only. No database was started; timings below are read from saved result cells.

## Input manifest

SHA-256 hashes normalize CRLF to LF for checkout portability. Corpus digest frames sorted repository-relative paths and normalized byte lengths followed by file bytes.

| run | JSON SHA-256 | producing commit | run tag | cells | CPU ids | MemTotal kB |
|---|---|---|---|---|---|---|
| select-access-pg-1-20261003 | 2bb893ca5c35fe65680e7f58e0a929c78dab56966406e2e7f8c79e35cc49dc43 | 524ce328afa3e738f095b71214051240bcc6f1f0 | run/study-01/select-access-pg-1-20261003 | 34 | 0-7 | 16109780 |
| select-access-pg-2-20261003 | 24ceae2fe2933f2a15f44d6f7d85f180d5c90491917635a5e76d79ba5bd2d691 | 524ce328afa3e738f095b71214051240bcc6f1f0 | run/study-01/select-access-pg-2-20261003 | 34 | 0-7 | 16109780 |
| select-access-pg-3-20261003 | 5e5325b883ecca3c00553c693725e682bb7d7f0500c0fb620a868c2e6d8adc03 | 524ce328afa3e738f095b71214051240bcc6f1f0 | run/study-01/select-access-pg-3-20261003 | 34 | 0-7 | 16109780 |
| select-access-pg-pos-20261003 | 8747043ad53b9991914f1e57664f06e957377a5978f165648a2893145a62bb65 | 6ad01ef65352b0be194007dab460b297ae94a541 | run/study-01/select-access-pg-pos-20261003 | 35 | 0-7 | 16109780 |
| select-access-yb-1-20261003 | ce2139aceea0cd56889486f4484b112e43ab26f92797c53a034d3d9f617468a3 | 6ad01ef65352b0be194007dab460b297ae94a541 | run/study-01/select-access-yb-1-20261003 | 34 | 0-7 | 16109780 |
| select-access-yb-2-20261003 | bf5210d941f968030b6f832c59c9bd2585eebf68eb00695e2a18bf113b1056dd | 6ad01ef65352b0be194007dab460b297ae94a541 | run/study-01/select-access-yb-2-20261003 | 35 | 0-7 | 16109780 |
| select-access-yb-3-20261003 | 13e9fc78aae1a024d810a2237f697e3749d006be26d454df50463902b1fc69ec | 6ad01ef65352b0be194007dab460b297ae94a541 | run/study-01/select-access-yb-3-20261003 | 34 | 0-7 | 16109780 |
| select-agg-pg-1b-20261003 | e21889935810d969a5aa30eb1372626417e25929be34437dc01b7a2f2e46f54f | 08a3c55396e01e07799a1014ee4c99b4b59f964c | run/study-01/select-agg-pg-1b-20261003 | 30 | 0-7 | 16109780 |
| select-agg-pg-2-20261003 | 41eb6645c6ec0c3c5c9cb01d01a64e4bb4dbba2db087844f9ef9b8fb3d1a6671 | 33b769017da9648bef92cb1a804697e42c957e90 | run/study-01/select-agg-pg-2-20261003 | 30 | 0-7 | 16109780 |
| select-agg-pg-3-20261003 | ada5ed4dbaf2ef18f4cdfee3de54a55d1bc7d440246f53f850b63795a16a35a9 | 33b769017da9648bef92cb1a804697e42c957e90 | run/study-01/select-agg-pg-3-20261003 | 31 | 0-7 | 16109780 |
| select-agg-yb-1-20261003 | 54068ef4024eeae7adf5213110c0ed9b7a74edf23d6241a859756b7ecc6333c8 | 33b769017da9648bef92cb1a804697e42c957e90 | run/study-01/select-agg-yb-1-20261003 | 30 | 0-7 | 16109780 |
| select-agg-yb-2-20261003 | 6fef86fb346299c135482acc9dfd37e6fc8ed0116a42b8d6a42256153c7f9251 | 33b769017da9648bef92cb1a804697e42c957e90 | run/study-01/select-agg-yb-2-20261003 | 31 | 0-7 | 16109780 |
| select-agg-yb-3-20261003 | 9e91e3fa58f3e19db0a9fc1bbb6368558f922e9d5c6977eed5ce10d724bd73a0 | 33b769017da9648bef92cb1a804697e42c957e90 | run/study-01/select-agg-yb-3-20261003 | 30 | 0-7 | 16109780 |
| select-exist-pg-1b-20261004 | b28df54d3fe6122d6452085fe2f4c9861f9e8d1e120535b6532dadfe10280051 | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-pg-1b-20261004 | 21 | 0-7 | 16109780 |
| select-exist-pg-2-20261004 | 5c768bd584f1d6e7d06f2124a40c0b909a1213c19fc08a4264fa807a4fbd6145 | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-pg-2-20261004 | 22 | 0-7 | 16109780 |
| select-exist-pg-3-20261004 | 4530bb3a280d9febe7efa3f0b8acb638e1c3a3bec8ef28b53c28f1c1bea947bf | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-pg-3-20261004 | 21 | 0-7 | 16109780 |
| select-exist-yb-1-20261004 | 228927f1b8b87df702b52ddcb6658f239228012729216b82eb4d4771d125e724 | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-yb-1-20261004 | 21 | 0-7 | 16109780 |
| select-exist-yb-2-20261004 | 819c269cda2684002b61e68833e882ad6e17992b21d3cb2afbdb2424e88d32b8 | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-yb-2-20261004 | 22 | 0-7 | 16109780 |
| select-exist-yb-3-20261004 | 531d48754b8121f4645fcb54acfa2bdb12193430afc7d11088634697a49c4c87 | 805b5b85828387ee6a74912e815dbbfe56105efa | run/study-01/select-exist-yb-3-20261004 | 21 | 0-7 | 16109780 |

Corpus inputs_digest: 02e75129215498d7066b870c1e339d388e7b54efbe72c8a2d7fc0a1e7b3b0b60

19 runs; 550 valid timed cells; 2750 saved drained-query samples.

## Ratios read directly from cells

Ratio below 1 means numerator has lower saved median; this is not a confidence test.

| run | pair/level | ratio |
|---|---|---|
| select-access-pg-1-20261003 | keyset/OFFSET depth=1 | 1.145311 |
| select-access-pg-1-20261003 | keyset/OFFSET depth=10 | 1.091467 |
| select-access-pg-1-20261003 | keyset/OFFSET depth=100 | 1.101705 |
| select-access-pg-1-20261003 | keyset/OFFSET depth=1000 | 0.832145 |
| select-access-pg-1-20261003 | keyset/OFFSET depth=10000 | 0.228863 |
| select-access-pg-2-20261003 | keyset/OFFSET depth=1 | 1.191034 |
| select-access-pg-2-20261003 | keyset/OFFSET depth=10 | 1.164121 |
| select-access-pg-2-20261003 | keyset/OFFSET depth=100 | 1.061632 |
| select-access-pg-2-20261003 | keyset/OFFSET depth=1000 | 0.706340 |
| select-access-pg-2-20261003 | keyset/OFFSET depth=10000 | 0.196873 |
| select-access-pg-3-20261003 | keyset/OFFSET depth=1 | 1.203252 |
| select-access-pg-3-20261003 | keyset/OFFSET depth=10 | 1.171813 |
| select-access-pg-3-20261003 | keyset/OFFSET depth=100 | 0.970294 |
| select-access-pg-3-20261003 | keyset/OFFSET depth=1000 | 0.732055 |
| select-access-pg-3-20261003 | keyset/OFFSET depth=10000 | 0.233641 |
| select-access-pg-pos-20261003 | keyset/OFFSET depth=1 | 1.147624 |
| select-access-pg-pos-20261003 | keyset/OFFSET depth=10 | 1.135761 |
| select-access-pg-pos-20261003 | keyset/OFFSET depth=100 | 0.988493 |
| select-access-pg-pos-20261003 | keyset/OFFSET depth=1000 | 0.749109 |
| select-access-pg-pos-20261003 | keyset/OFFSET depth=10000 | 0.230635 |
| select-access-yb-1-20261003 | keyset/OFFSET depth=1 | 1.260553 |
| select-access-yb-1-20261003 | keyset/OFFSET depth=10 | 1.271760 |
| select-access-yb-1-20261003 | keyset/OFFSET depth=100 | 1.122058 |
| select-access-yb-1-20261003 | keyset/OFFSET depth=1000 | 1.431872 |
| select-access-yb-1-20261003 | keyset/OFFSET depth=10000 | 0.538512 |
| select-access-yb-2-20261003 | keyset/OFFSET depth=1 | 1.177539 |
| select-access-yb-2-20261003 | keyset/OFFSET depth=10 | 1.258110 |
| select-access-yb-2-20261003 | keyset/OFFSET depth=100 | 1.101538 |
| select-access-yb-2-20261003 | keyset/OFFSET depth=1000 | 1.168841 |
| select-access-yb-2-20261003 | keyset/OFFSET depth=10000 | 0.673710 |
| select-access-yb-3-20261003 | keyset/OFFSET depth=1 | 1.176862 |
| select-access-yb-3-20261003 | keyset/OFFSET depth=10 | 1.182503 |
| select-access-yb-3-20261003 | keyset/OFFSET depth=100 | 1.249952 |
| select-access-yb-3-20261003 | keyset/OFFSET depth=1000 | 1.331411 |
| select-access-yb-3-20261003 | keyset/OFFSET depth=10000 | 0.691768 |
| select-agg-pg-1b-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 0.648007 |
| select-agg-pg-1b-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 0.811336 |
| select-agg-pg-1b-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.810626 |
| select-agg-pg-1b-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.653396 |
| select-agg-pg-1b-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 0.700854 |
| select-agg-pg-2-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 0.969459 |
| select-agg-pg-2-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 1.102281 |
| select-agg-pg-2-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.806195 |
| select-agg-pg-2-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.709316 |
| select-agg-pg-2-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 0.630136 |
| select-agg-pg-3-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 1.237880 |
| select-agg-pg-3-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 0.982812 |
| select-agg-pg-3-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.777042 |
| select-agg-pg-3-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.730651 |
| select-agg-pg-3-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 0.641957 |
| select-agg-yb-1-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 0.953399 |
| select-agg-yb-1-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 1.002720 |
| select-agg-yb-1-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.933781 |
| select-agg-yb-1-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.894380 |
| select-agg-yb-1-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 0.818487 |
| select-agg-yb-2-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 0.913451 |
| select-agg-yb-2-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 1.045632 |
| select-agg-yb-2-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.885301 |
| select-agg-yb-2-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.856465 |
| select-agg-yb-2-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 0.831846 |
| select-agg-yb-3-20261003 | DISTINCT ON/ROW_NUMBER parents=1 | 1.012456 |
| select-agg-yb-3-20261003 | DISTINCT ON/ROW_NUMBER parents=10 | 0.964455 |
| select-agg-yb-3-20261003 | DISTINCT ON/ROW_NUMBER parents=100 | 0.990092 |
| select-agg-yb-3-20261003 | DISTINCT ON/ROW_NUMBER parents=1000 | 0.922614 |
| select-agg-yb-3-20261003 | DISTINCT ON/ROW_NUMBER parents=10000 | 1.015397 |

## Recorded execution sequences (position repeats excluded)

| family/engine | distinct sequences | runs |
|---|---|---|
| access/postgres | 1 | 4 |
| access/yugabyte | 1 | 3 |
| aggregate/postgres | 1 | 3 |
| aggregate/yugabyte | 1 | 3 |
| existence/postgres | 1 | 3 |
| existence/yugabyte | 1 | 3 |
