# Study 01 v3 harness validation

## TL;DR

- Go tests with the race detector and `go vet` passed in the pinned Go container.
- The catalogue test first found an accidentally copied documentation placeholder;
  after correcting it, all intended query deltas passed.
- Seven new PlantUML diagrams rendered successfully in a container.
- Database correctness and performance have not yet been established by these checks.

Environment: `host-zenbook-ux5406sa`. Implementation branch:
`study-01/measurement-enhancements`; milestone `study-01/v3-harness`.

The tests exercise scheduled overload (demand conservation, visible rejected arrivals,
queue-inclusive latency), unchanged SQL between mechanism variants, preservation of
independent trials and errors in reporting, and deterministic longer histories.

Commands run inside Podman using `docker.io/library/golang:1.26-bookworm`:

```text
go test -race ./...
go vet ./...
bash -n studies/01-charity-tree/run-enhancements.sh
```

The syntax check mounts the repository root; Go commands mount the study as `/src`.
Initial Go compilation also passed before the tests were added. The initial diagram
generation failed because PowerShell flattened a constructed line array; parenthesizing
the concatenated elements restored line breaks and all seven SVGs then rendered.
No database timing is inferred from these build and tooling checks.
