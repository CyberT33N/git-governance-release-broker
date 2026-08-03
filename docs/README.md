# Broker documentation

This directory defines the architecture, verification, operations, and
governance contracts for `git-governance-release-broker`.

The Broker is a Linux AMD64 Cloud Run service. It mints only short-lived,
repository-bound GitHub App installation tokens for authenticated workloads.
It does not create branches, pull requests, tags, releases, or production
deployments.

## Document map

- `architecture/ADR-0001-LINUX-ONLY-BROKER-DELIVERY.md` records why the
  production delivery contract is Linux AMD64 only.
- `development/VERIFICATION.md` defines local safe verification and the
  Linux production quality gate.
- `operations/CI-CD-SUPPLY-CHAIN.md` defines CI/CD identity, artifact, and
  evidence boundaries.
- `specification/CI-RULESET-CONTRACT.md` maps emitted GitHub checks to
  required Ruleset controls.
- `conventions/BUILD-AND-DEPENDENCY-CONVENTIONS.md` defines the Go module,
  toolchain, and build conventions.
- `TRACEABILITY.md` links tickets to implemented governance evidence.

The external approved Go proxy, hermetic build image, and artifact evidence
registry are not provisioned yet. Production delivery remains fail-closed
until those foundations exist.
