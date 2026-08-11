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

The staging package is an `evidence-graph/v1` collection of immutable,
separately signed subject documents. It records Source, Dependency Resolution,
Build, Artifact, Promotion, Deployment, and Operation without mutating a
previous subject document. Source through Artifact contain the materialized
staging evidence: source commit and tree, the complete resolved Go module
graph, hermetic module verification/test result, builder definition and
toolchain, immutable image digest, SBOM, registry signature, and GitHub
attestation bundles.

## Controlled local builder materialization

Before its isolated Go module and test consumer phase, the staging workflow
explicitly materializes the full builder digest declared by the Dockerfile. It
pulls that `linux/amd64` digest, verifies the local repository-digest and
platform binding, and tags the resulting local image under an ephemeral
digest-derived reference.

All offline Go commands use only that local reference with `--pull=never` and
`--network=none`; the final staging image build likewise uses
`--pull=false --network=none`. A missing, mismatched, or wrong-platform local
builder image fails before the final image build, image push, evidence upload,
or Cloud Run mutation. BuildKit, Docker, runner, and Node caches are not an
authority for builder availability or trust.

This local materialization establishes only the exact local image binding
needed by the isolated consumer. It does not make the builder a verified
builder artifact: separate builder signature, SBOM, provenance, policy, and
revocation evidence remain required before any artifact subject can become
`verified`.

The staging package contains at least:

```text
source.subject.json
dependency-resolution.subject.json
dependency-resolution.json
build.subject.json
artifact.subject.json
promotion.not-recorded.subject.json
operation.not-recorded.subject.json
<subject>.integrity.sigstore.json
broker.spdx.json
signature.json
provenance.intoto.jsonl
sbom.intoto.jsonl
test-result.json
```

Each materialized subject stores a canonical payload digest and a keyless
Sigstore bundle over that payload. Its `relations[]` array binds source to
dependency resolution, source and dependency resolution to build, and build to
the final artifact. A phase that has not occurred is explicitly
`not-recorded`; it never grants promotion, deployment, or runtime admission.

A main-bound promoter creates a new lane-specific `promotion.subject.json`
with its own canonical payload, signature bundle, target digest, environment
approval evidence, and relation to the verified artifact subject. A deployment
records a new `deployment.subject.json` only after it observes the exact
deployed Cloud Run revision, immutable digest, and ready condition. Production
deployment verifies the full upstream subject chain and verified promotion
subject before a Cloud Run mutation.

Cloud Run may canonicalize an OCI image reference while retaining its immutable
digest. Deployment evidence therefore compares the deployed and approved
digests, not the original reference string. Its fail-closed guards emit only
bounded validation codes for missing or invalid evidence, revision readiness,
and digest mismatch; they never print credentials, headers, or token values.
The Cloud Run v2 control plane is observed through the documented
`gcloud run revisions describe` compatibility projection: readiness requires
`status.conditions[?type="Ready"].status=True`, and the deployed image is read
from `spec.containers.image`. Raw v2 REST field paths are not interchangeable
with this CLI projection.

The approved Go proxy, dependency admission, immutable scan and quality
evidence, hermetic build image, operation-evidence writer, lane-specific
generic evidence repositories, and their IAM boundaries remain external
prerequisites. Until they are materialized and re-verified, artifact subjects
remain `pending`, promotions and production deployments fail closed, and no
workflow may claim a completed Supply-Chain-Fortress production delivery.

## Cosign v3 evidence contract

The staging signer and composite evidence verifier explicitly pin
`sigstore/cosign-installer` `v4.1.2` and the Cosign binary `v3.1.3`. The
installer action release and installed binary version are separate contracts;
both pins are required so the installer can verify the V3 keyless release
bundles before making the binary available to the workflow.

Cosign v3 uses bundle-based OCI evidence by default. The broker workflows use
that default and do not retain the former experimental OCI environment setting,
the sign-only referrer-mode override, or the removed verify discovery flag.
Any future Cosign upgrade requires a new governed compatibility review of the
signing, verification, SBOM, provenance, and immutable evidence-package
contracts.

## Registry authentication for attestations

GitHub's registry attestation action requires a static Docker `auths` entry; it
does not consume Google Cloud Docker credential helpers. The staging workflow
therefore exchanges its federated identity for a short-lived access token and
uses `docker/login-action` with `oauth2accesstoken` and password-stdin.

The composite evidence verifier creates an isolated temporary `DOCKER_CONFIG`,
logs in with a short-lived `gcloud auth print-access-token` value, then logs
out and removes that directory through its shell cleanup trap. No long-lived
credential, private key, or token value is written to the repository, workflow
output, or GitHub environment variables.

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
