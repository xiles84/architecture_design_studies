# Execution handoff

Goal: restore dependable queue execution and recovered-claim submission.

The SELECT worker reports that recovery leaves stale worktree context. Current main already saves that context and updates the live event pointer; the installed `localhost/ads-queue:1` image was built on 2026-09-22. Verify the discrepancy, add a source provenance gate for the image, rebuild in Podman with existing correctness tests, and exercise the recovered-claim lifecycle without production ref changes. Preserve all scientific evidence and other workers' files. Record meaningful decisions and every validation result. Merge completed changes locally and tag them. Only claim SELECT separately after its grace period and explicit end-of-work handoff have been verified.

NEXT MODEL: HIGH
