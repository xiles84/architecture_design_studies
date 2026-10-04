# Independent correction-input check — GPT-6, 2026-10-04

TL;DR: immutable inputs checked without a database run.

| run | cells | samples | producing tag | normalized SHA-256 |
|---|---:|---:|---|---|
| select-order-access-postgres-1-20261004 | 35 | 175 | `run/study-01/select-order-access-postgres-1-20261004` | `0783debd960ffb0e76f9cf0771a11a3c065e433a6e1efa11519f06447204a0f4` |
| select-order-access-postgres-2-20261004 | 34 | 170 | `run/study-01/select-order-access-postgres-2-20261004` | `357620cf7fb2687698f082b5d1a7294fe72f41fb8c30d93def9d9156a43c12f5` |
| select-order-access-postgres-3-20261004 | 34 | 170 | `run/study-01/select-order-access-postgres-3-20261004` | `f0b0f2c1bf41c51187b6a1fe78d73d735bd6d0a5fcd98ac849bfc37fc0c56aa0` |
| select-order-access-yugabyte-1-20261004 | 35 | 175 | `run/study-01/select-order-access-yugabyte-1-20261004` | `9a3b8b777b99388a1cb983808a53bb0f976aad693cb14e7dd4f4fe8da8b99d7f` |
| select-order-access-yugabyte-2-20261004 | 34 | 170 | `run/study-01/select-order-access-yugabyte-2-20261004` | `818eaca45753f59f49626bb1135e68dec19cf9359c611f6048632cbb478e4b9d` |
| select-order-access-yugabyte-3-20261004 | 34 | 170 | `run/study-01/select-order-access-yugabyte-3-20261004` | `ee27b23999ddd2405173c950ed2128986d7f723807e30232c02190a5532c2e52` |
| select-order-aggregate-postgres-1-20261004 | 31 | 155 | `run/study-01/select-order-aggregate-postgres-1-20261004` | `76268c647e965cd579b85330b32ccdf716c959230c56fe0d71216ce23ab39dc2` |
| select-order-aggregate-postgres-2-20261004 | 30 | 150 | `run/study-01/select-order-aggregate-postgres-2-20261004` | `b742e387f0cd26e3a730227db33212c88c0535251e947928e0e6c480cbce0b35` |
| select-order-aggregate-postgres-3-20261004 | 30 | 150 | `run/study-01/select-order-aggregate-postgres-3-20261004` | `5285f21e2330e7118cb9fd7d3fc06576ad2fcc8d6b5c0619825f0b3726cd4823` |
| select-order-aggregate-yugabyte-1-20261004 | 31 | 155 | `run/study-01/select-order-aggregate-yugabyte-1-20261004` | `21ad333d1e4ff8a9da998c1a748ec21beca168a49d122da6c11658dcb9ced477` |
| select-order-aggregate-yugabyte-2-20261004 | 30 | 150 | `run/study-01/select-order-aggregate-yugabyte-2-20261004` | `e21df42d7c564cc456eba1e11a5d0ec43eaed59d7849ccdfd21b1be390eb2ab7` |
| select-order-aggregate-yugabyte-3-20261004 | 30 | 150 | `run/study-01/select-order-aggregate-yugabyte-3-20261004` | `971d1d0df2e057b6aaae034b0834b4a706559e22842bdec96ecb822cddd03014` |
| select-order-existence-postgres-1-20261004 | 22 | 110 | `run/study-01/select-order-existence-postgres-1-20261004` | `a1cd2e20e7fd50961ed5562a099cd46e5b993e69847f8822daf4e3696d74f042` |
| select-order-existence-postgres-2-20261004 | 21 | 105 | `run/study-01/select-order-existence-postgres-2-20261004` | `5a12edf5d2b7a1df53338b019e7e77cec89e64d80302731659073759283f4f39` |
| select-order-existence-postgres-3-20261004 | 21 | 105 | `run/study-01/select-order-existence-postgres-3-20261004` | `b2cce0c8c6ad02cc39a45e93ddfe29dcd26f23988af6ebc3700a63965286eafd` |
| select-order-existence-yugabyte-1-20261004 | 22 | 110 | `run/study-01/select-order-existence-yugabyte-1-20261004` | `b8e15f58e4e81a72c956ddd1e23bb656a932cf3a33ce8b41f59111be6fe9a82c` |
| select-order-existence-yugabyte-2-20261004 | 21 | 105 | `run/study-01/select-order-existence-yugabyte-2-20261004` | `703b92c3644821ed9ddb95fb5a5e4566ea8d342cc63a998e02f5be3a11dc6749` |
| select-order-existence-yugabyte-3-20261004 | 21 | 105 | `run/study-01/select-order-existence-yugabyte-3-20261004` | `62979fd70cec76ac2cfc6f2441d73310d61cd0108485bb3abae839bbf451f262` |

18 runs; 516 valid cells; 2580 verified samples; 3 distinct recorded orders for each family/engine.

access corpus: `781fcb68ac4907615da0301344d5225445ebec91a8573687ae0ed6eb03c4a14e`.

aggregate corpus: `2df8f4499412a81f23a7f85d7952f25262e70abad40b69ccddac2c30300bde59`.

existence corpus: `e6b469b428b00501e0dcd64dd277aa72d7121e5a2e248c9dea1473f9ca246e0f`.
