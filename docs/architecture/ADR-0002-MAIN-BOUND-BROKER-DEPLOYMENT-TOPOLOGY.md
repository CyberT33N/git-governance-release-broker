# ADR-0002: Main-bound Broker deployment topology

## Status

Accepted.

## Context

The original Broker deployment workflow built and deployed directly from
`develop` through one environment, one Artifact Registry repository, and one
Cloud Run service. That makes unpromoted integration source a production
deployment authority.

The Broker includes two materially different production identities:

```text
release-automation
reconciliation-publisher
```

They must not share an App key, runtime identity, deployment identity, or
environment approval boundary.

## Decision

The deployment topology is:

```text
develop
→ protected staging environment
→ staging image repository
→ staging Broker service

main
→ protected production environment
→ immutable production digest
→ release-automation Broker service

main
→ protected reconciliation publisher environment
→ immutable publisher digest
→ reconciliation publisher Broker service
```

Production workflows accept only full `@sha256:` image references located in
their configured production repository. They never build or deploy `develop`
source.

## Consequences

- The former develop-bound deployment workflow is removed.
- Staging and production require distinct WIF, deployer, runtime, invoker,
  Artifact Registry, Cloud Run, Secret Manager, and GitHub Environment
  boundaries.
- Production image promotion, SBOM, provenance, signatures, and attestations
  remain external fail-closed prerequisites until the platform evidence lane is
  provisioned.
- A GitHub Environment cannot be treated as a production boundary until
  required reviewers, self-review prevention, branch restrictions, and
  administrator-bypass policy are all verified.
