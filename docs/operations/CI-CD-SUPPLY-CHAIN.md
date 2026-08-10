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
→ SPDX SBOM
→ keyless Sigstore signature
→ GitHub provenance and SBOM attestations
→ immutable generic evidence package
→ promotion-time evidence verification
→ deployment-time evidence verification
→ protected production deployment
```

The source contract signs and attests a staging digest, then copies an
immutable evidence package with that digest to each production lane. The
approved Go proxy, hermetic build image, lane-specific generic evidence
repositories, and their IAM boundaries remain external prerequisites. Until
they exist, no workflow may claim a completed Supply-Chain-Fortress production
delivery.

## Cosign v3 evidence contract

The staging signer and composite evidence verifier explicitly install Cosign
`v3.1.3`. The installer action release and the installed Cosign binary version
are separate contracts; the workflow must set `cosign-release` explicitly.

Cosign v3 uses bundle-based OCI evidence by default. The broker workflows use
that default and do not retain the former experimental OCI environment setting,
the sign-only referrer-mode override, or the removed verify discovery flag.
Any future Cosign upgrade requires a new governed compatibility review of the
signing, verification, SBOM, provenance, and immutable evidence-package
contracts.

## Superseded pre-delivery candidates

A protected release candidate is not delivered merely because its ref exists.
If a candidate becomes an ancestor of `main` before its promotion and has no
immutable tag, published release, or delivery evidence, it is recorded as
`superseded-before-delivery`.

Such a candidate must not receive retrospective artifacts, SBOMs, signatures,
attestations, tags, or GitHub Releases. Its ref remains available for audit
until a successor release has completed delivery and a controlled retention
decision permits cleanup.

## Deployment topology

The staging deployment workflow is develop-bound and builds only the isolated
staging Broker. Production workflows are main-bound and accept only immutable
image digests from their configured production repository.

```text
gcp-broker-staging
→ staging release-automation profile

gcp-broker-production
→ production release-automation profile

gcp-release-credential-verification-deployment
→ production release-credential-verification profile

gcp-hotfix-delivery-deployment
→ production hotfix-delivery profile

gcp-reconciliation-publisher-deployment
→ production reconciliation-publisher profile

gcp-hotfix-propagation-publisher-deployment
→ production hotfix-propagation-publisher profile
```

The six environments, service identities, secrets, image repositories, and
Cloud Run services must remain separate. No deployer, runtime, or invoker
identity receives permissions across those boundaries.

Every production lane also owns a generic evidence repository. A promoter may
read only the staging image and staging evidence repositories, then write only
the lane's image and evidence repositories. A deployer may read only its
lane's evidence repository and the staging image repository needed to
re-validate the source signature and attestations.

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
