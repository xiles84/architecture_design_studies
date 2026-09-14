# ER-01 diagnosis — 20260914T164638Z

- repo_commit: f25032d5e30e4fea3f5587b4a8210dc76f012ff3 (dirty=false); design s1_conditional_update; tiny; race 10 000 seats; 8 trials per configuration

| Configuration | Early rejections (sum over trials) | Races with at least one | Trials run | Errors | Seats held/s per trial |
|---|---:|---:|---:|---:|---|
| A-default | 2 | 2 | 8 | 0 | 19.0 34.3 31.6 23.9 24.6 36.8 28.5 29.1  |
| B-no-pushdown | 2 | 2 | 8 | 0 | 27.0 28.9 24.3 28.3 25.0 27.6 30.9 25.6  |
| C-no-wait-queues | 2 | 2 | 8 | 0 | 36.2 33.1 27.0 29.0 29.4 29.9 31.4 29.4  |
