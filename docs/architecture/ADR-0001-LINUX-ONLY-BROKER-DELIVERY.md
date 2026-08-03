# ADR-0001: Linux-only Broker delivery

## Status

Accepted.

## Context

The Broker is delivered as a `linux/amd64` container to Cloud Run. Its
Dockerfile builds with `CGO_ENABLED=0`, `GOOS=linux`, and `GOARCH=amd64`,
then runs in a non-root Scratch image.

The Broker does not publish Windows or macOS binaries and contains no
platform-specific production path.

## Decision

The required production quality check is:

```text
Quality gates (linux-amd64)
```

It verifies the Go module, formatting, tests, complete statement coverage,
race safety, static analysis, Linux AMD64 binary build, embedded module
provenance, and fuzz smoke tests.

Container build, SBOM, provenance, signature, and attestation verification are
added only after the approved internal Go proxy, hermetic build image, and
artifact evidence registry exist.

Windows and macOS runtime checks are not required Broker release gates. They
do not validate the deployed Cloud Run artifact and must not remain stale
Ruleset requirements.

## Consequences

- Local developer tests remain allowed without production credentials.
- Production secrets, GitHub App keys, deployment identities, and Cloud Run
  mutation remain CI/CD-only.
- A future Windows or macOS Broker delivery contract requires a new ADR,
  platform artifacts, and matching required checks.
- The CI workflow and GitHub Rulesets must use identical emitted check names.
