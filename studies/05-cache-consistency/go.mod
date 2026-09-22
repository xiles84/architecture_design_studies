module cacheconsistency

go 1.24

require adsplatform v0.0.0

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

// The shared measurement core, ports and adapters live in the repository's
// platform module. A replace directive (not a published version) keeps the study
// pinned to the platform code in the same commit, which is what the run tags and
// the repository version recorded in every result refer to.
replace adsplatform => ../../platform
