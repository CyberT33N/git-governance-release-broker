# ADR-0002: Main-bound Broker deployment topology

## Status

Accepted.

## Context

The original Broker deployment workflow built and deployed directly from
`develop` through one environment, one Artifact Registry repository, and one
Cloud Run service. That makes unpromoted integration source a production
deployment authority.

The Broker includes five materially different production identities:

```text
release-automation
reconciliation-publisher
hotfix-propagation-publisher
release-credential-verification
hotfix-delivery
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
→ protected release credential verification environment
→ immutable verification digest
→ release-credential-verification Broker service

main
→ protected hotfix delivery environment
→ immutable hotfix delivery digest
→ hotfix-delivery Broker service

main
→ protected reconciliation publisher environment
→ immutable publisher digest
→ reconciliation publisher Broker service

main
→ protected hotfix propagation publisher environment
→ immutable publisher digest
→ hotfix propagation publisher Broker service
```

Production workflows accept only full `@sha256:` image references located in
their configured production repository. They never build or deploy `develop`
source.

## Consequences

- The former develop-bound deployment workflow is removed.
- Staging and production require distinct WIF, deployer, runtime, invoker,
  Artifact Registry, Cloud Run, Secret Manager, and GitHub Environment
  boundaries.
- Staging signs and attests immutable image digests, then stores an immutable
  evidence package. Promotion and deployment re-verify that evidence before
  production mutation.
- Lane-specific generic evidence repositories, the approved Go proxy, and the
  hermetic build image remain external fail-closed prerequisites until
  provisioned.
- A GitHub Environment cannot be treated as a production boundary until
  required reviewers, self-review prevention, branch restrictions, and
  administrator-bypass policy are all verified.
