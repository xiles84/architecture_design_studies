// Command queue is the durable AI work queue CLI described in
// docs/ai-work/WORKFLOW.md and docs/ai-work/SCHEMA.md.
//
// It runs from the pinned Podman image built by tools/queue/Containerfile (the
// host wrapper tools/queue/queue starts it); the only requirement on the host
// is Podman. Live coordination lives in Git refs, the permanent archive in
// committed event files.
package main

import (
	"os"

	"adsqueue/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))
}
