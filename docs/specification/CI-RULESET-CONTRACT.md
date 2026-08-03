# CI and Ruleset contract

## Required check names

The Broker CI emits these server-side controls:

```text
Quality gates (linux-amd64)
CodeQL (go)
Dependency admission review
```

The `develop` and `main` Rulesets must reference the same emitted check names.
A Ruleset must not require stale Windows or macOS check contexts for this
Linux-only Cloud Run service.

## Linux quality gate

`Quality gates (linux-amd64)` runs on Ubuntu and verifies:

```text
exact Go 1.26.5
GOTOOLCHAIN=local
GOFLAGS=-mod=readonly
GOVCS=*:off
source-level build contract
configuration fuzz smoke
HTTP boundary fuzz smoke
Linux AMD64 source binary build
```

The hermetic container build, SBOM, provenance, signature, and attestation
lane remains externally blocked until the approved proxy, build image, and
artifact evidence registry are provisioned.

## Code scanning and dependency admission

`CodeQL (go)` uploads Go code-scanning results. `Dependency admission review`
reviews dependency changes for pull requests targeting protected integration or
production lines.

## Ruleset restrictions

Rulesets remain fail-closed:

```text
required review
required check contexts
no force push
no deletion
no Ruleset bypass
no stale or non-emitted required contexts
```
