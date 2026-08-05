# CI/CD and Supply-Chain operations

## Trust boundaries

```text
Developer workstation
→ source tests only
→ no production secret, deployment, or GitHub App authority

CI quality lane
→ Linux AMD64 source and container verification
→ no production deployment

Production deployment lane
→ immutable image digest
→ dedicated workload identity
→ private Cloud Run service
→ runtime-only Secret Manager access
```

## Required CI controls

```text
Quality gates (linux-amd64)
CodeQL (go)
Dependency admission review
```

GitHub Rulesets must require only checks that the repository actually emits.
Windows and macOS runtime check contexts are not part of this Broker contract.

## Automated dependency intake

Dependabot opens bounded daily pull requests against `develop` for:

```text
Go modules
GitHub Actions
Docker base images
```

The `develop` target applies to version updates. GitHub Dependabot security
updates always target the repository default branch, `main`. They are intake
signals only and must be triaged into the governed hotfix process; they do not
authorize a direct main merge.

Dependabot does not approve, bypass, merge, or deploy changes. Every update
remains subject to dependency admission review, CodeQL, Linux quality gates,
required review, and protected-branch policy.

## Artifact evidence

The target production artifact flow is:

```text
reviewed source revision
→ approved Go module graph
→ hermetic Linux AMD64 build
→ immutable container digest
→ SBOM
→ vulnerability and policy evidence
→ provenance, signature, and attestation
→ protected production deployment
```

The main-bound artifact-promotion workflows resolve the staging image from the
reviewed `sha-<commit>` tag, require that commit to be reachable from `main`,
copy the image without rebuilding it, compare source and destination digests,
and emit only the destination `broker@sha256:<digest>` reference. Separate
promoter identities may read the staging repository and write only their own
production repository; they cannot deploy Cloud Run services or read runtime
secrets.

The approved Go proxy, hermetic build image, and evidence registry have not
yet been provisioned. Until they exist, no workflow may claim a completed
Supply-Chain-Fortress production delivery.

## Deployment topology

The staging deployment workflow is develop-bound and builds only the isolated
staging Broker. Production workflows are main-bound and accept only immutable
image digests from their configured production repository.

```text
gcp-broker-staging
→ staging release-automation profile

gcp-broker-production
→ production release-automation profile

gcp-reconciliation-publisher-deployment
→ production reconciliation-publisher profile
```

The three environments, service identities, secrets, image repositories, and
Cloud Run services must remain separate. No deployer, runtime, or invoker
identity receives permissions across those boundaries.

## Incident handling

If a module, toolchain, build image, App identity, or artifact evidence source
is suspected to be compromised:

```text
stop promotion
revoke or rotate affected identities
identify affected image digests and module graphs
rebuild from approved inputs
publish only new immutable evidence
document the recovery and exception decision
```
