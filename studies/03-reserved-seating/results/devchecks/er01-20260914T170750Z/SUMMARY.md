# ER-01 diagnosis — 20260914T170750Z

- repo_commit: 87fa9fac22da7379a068f243da59e3ff20f98c2a (dirty=false); design s1_conditional_update; tiny; race 10 000 seats; 8 trials per configuration

| Configuration | Early rejections (sum over trials) | Races with at least one | Trials run | Errors | Seats held/s per trial |
|---|---:|---:|---:|---:|---|
| D-log-restarts | 1 | 1 | 8 | 0 | 25.8 27.1 27.0 25.5 32.4 24.1 26.6 28.0  |
| E-no-internal-retries | 1 | 1 | 8 | 2520 | 22.2 31.1 27.2 30.2 23.0 24.0 24.2 22.1  |

Re-issued statement inside the refusing transaction (round 2 onwards):

- D-log-restarts:       2 issued again in this transaction matched 2 of 2
- E-no-internal-retries:       2 issued again in this transaction matched 6 of 6
