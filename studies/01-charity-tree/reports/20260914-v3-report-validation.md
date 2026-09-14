# Study 01 v3 — reporting and matched-control validation

## TL;DR

- Containerized Go race tests and `go vet` passed, including the new test that unavailable
  YugabyteDB storage must not be displayed as a measured zero.
- The reporter now shows growth read rates, arrival retries and acknowledged counts,
  first write errors, and warmup/final audit counts beside each trial.
- The matched memory group repeats the same medium/history-multiplier=2 dataset and
  mutation sequence with standard and constrained memory configurations.

Environment: `host-zenbook-ux5406sa`. Author: GPT-6 via Codex desktop, 2026-09-14.
All Go execution used `golang:1.26-bookworm` in Podman, with the shared benchmark lock,
2 CPUs and 2 GiB for the check container. No database was started by these checks.

```text
gofmt -w harness/report_enhancements.go harness/experiment_test.go
go test -race ./...
go vet ./...
bash -n run-enhancements.sh run-study.sh
```

The first host invocation selected the WSL shell shim and failed before starting a
container. Re-running through the explicit Windows Git Bash executable succeeded.
The matched-control cell addition received a separate final shell syntax check.
The runner now installs a lock-release trap immediately after acquisition, so a dirty
tree or build failure before database startup cannot leave this session's lock behind.
These checks establish tooling behavior; database evidence remains in tagged run reports.
