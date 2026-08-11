# Traceability

## GOV-69: Normalize Cloud Run revision evidence

Status: in implementation.

Scope:

```text
Read the Cloud Run revision once as JSON and normalize the CLI-compatible and
raw-v2 condition and container shapes locally. Require exactly one Ready
condition with True or CONDITION_SUCCEEDED and exactly one immutable deployed
image before recording a deployment subject. Do not use a fragile gcloud
--format list-filter projection for runtime admission.
```

## GOV-68: Project gcloud revision readiness

Status: superseded by GOV-69.

Scope:

```text
The documented CLI list-filter projection still produced an empty condition in
the live staging workflow. GOV-69 replaces single-shape projections with a
JSON-normalized, fail-closed interpretation.
```

## GOV-67: Project Cloud Run v2 revision evidence

Status: superseded by GOV-68.

Scope:

```text
The initial raw Cloud Run v2 REST projection did not match the response model
of gcloud run revisions describe. GOV-68 replaces it with the documented CLI
projection without weakening immutable digest binding.
```

## GOV-66: Check Cloud Run v2 deployment readiness

Status: superseded by GOV-68.

Scope:

```text
The initial state-based readiness check did not match the gcloud CLI projection.
GOV-68 retains the fail-closed readiness and digest guards using the documented
status-based compatibility fields.
```

## GOV-65: Record digest-bound staging deployment evidence

Status: in implementation.

Scope:

```text
Record a deployment subject only after the observed Cloud Run revision is
ready and its immutable digest matches the approved staging digest. Accept
provider canonicalization of the OCI reference string without weakening the
digest binding. Replace silent fail-closed guards with bounded, token-free
validation codes for artifact-subject, promotion-subject, revision-readiness,
and deployed-digest failures.
```

## GOV-64: Materialize the local broker builder image

Status: in implementation.

Scope:

```text
Materialize the full digest declared by the Dockerfile before isolated Go
module and test consumption. Verify the local repository-digest and
linux/amd64 platform binding, then use only a digest-derived local builder
reference with --pull=never and --network=none. Build the final staging image
only after that consumer phase succeeds, with --pull=false and --network=none.

The local materialization is an availability and digest-binding control, not a
verified builder-artifact claim. Builder signature, SBOM, provenance, policy,
and revocation evidence remain required before an Artifact Subject can become
verified. Missing or mismatched local builder materialization must fail before
image publication, evidence upload, or Cloud Run mutation.
```

## GOV-63: Conform broker evidence subjects to evidence-graph/v1

Status: in implementation.

Scope:

```text
Replace the flat staging manifest with immutable evidence-graph/v1 subject
documents for Source, Dependency Resolution, Build, Artifact, Promotion,
Deployment, and Operation. Each materialized subject has a canonical-payload
digest and a keyless Sigstore signature bundle. Source, dependency, build, and
artifact relationships are explicitly typed; promotion and deployment create
new append-only documents instead of rewriting the artifact document.

The staging lane records a complete resolved Go module graph and hermetic
module verification/test evidence. The artifact remains pending until the
approved proxy, admission, quality, scan, policy, and approval evidence are
actually materialized. Promotion and production deployment fail closed unless
the upstream artifact chain is verified. Deployment evidence is persisted only
after the deployed Cloud Run revision, immutable digest, and readiness are
observed.
```

## GOV-62: Bind staging registry authentication to deploy job

Status: in implementation.

Scope:

```text
Emit the short-lived Artifact Registry access token in the deploy job that
consumes it. Keep bootstrap authenticated for registry provisioning without a
deploy-only token output, and lock the job-local output and ordering contract
with same-package workflow tests.
```

## GOV-61: Authenticate registry attestations

Status: integrated into develop; staging deployment remediation continues in GOV-62.

Scope:

```text
Use a short-lived Workload Identity Federation access token to create static
Docker registry auth for GitHub attestation uploads. Keep the token out of
logs, use password-stdin, isolate composite verification auth in a temporary
Docker configuration, and remove it after use. Lock the contract with
same-package workflow tests.
```

## GOV-60: Upgrade Cosign installer for v3 bundles

Status: integrated into develop through PR #33.

Scope:

```text
Pin both evidence workflow installation points to the Cosign installer release
that verifies Cosign v3 keyless bundles. Keep the explicit Cosign v3.1.3
binary pin and lock the compatible installer and binary pair with same-package
workflow contract tests.
```

## GOV-59: Upgrade Cosign v3 evidence workflows

Status: integrated into develop through PR #32.

Scope:

```text
Pin both evidence workflow installation points to Cosign v3.1.3. Adopt the
v3 bundle defaults for image signing and verification by removing obsolete
experimental OCI mode and referrer-discovery flags. Lock the v3 command
contract with same-package workflow tests.
```

## GOV-58: Enable Cosign verify OCI discovery

Status: integrated into develop through PR #31.

Scope:

```text
Use the supported --experimental-oci11 discovery flag on read-only cosign
verify invocations that validate OCI 1.1 referrer signatures. Retain the
sign-only registry-referrers mode, certificate identity, and issuer checks.
Lock the compatible signing and verification syntax with same-package workflow
contract tests.
```

## GOV-57: Preserve Cosign verify compatibility

Status: integrated into develop through PR #30.

Scope:

```text
Keep OCI 1.1 referrer mode scoped to immutable image signing. Retain
experimental mode only within the affected sign and verification steps. Remove
the unsupported registry-referrers flag from all read-only cosign verify
invocations while retaining certificate identity and issuer verification. Lock
the supported sign and verify syntax with same-package workflow contract tests.
```

## GOV-56: Enable Cosign OCI referrers

Status: integrated into develop through PR #29.

Scope:

```text
Enable COSIGN_EXPERIMENTAL=1 only in the staging signing step and the
read-only composite evidence verifier that use OCI 1.1 referrers.
Keep the experimental mode out of unrelated workflow steps and GitHub
Environment variables.
Lock both scoped settings with same-package workflow contract tests.
```

## GOV-55: Prepare staging evidence workspace

Status: in implementation.

Scope:

```text
Create .build/evidence before the staging SBOM action writes broker.spdx.json.
Keep .build/evidence as a temporary runner workspace.
Persist final SBOM, signature, attestations, and manifest only through the
immutable evidence package and digest-bound registry attachments.
Lock the ordering with a same-package workflow contract test.
```

## GOV-54: Broker container evidence gate

Status: in implementation.

Scope:

```text
SPDX SBOM for every immutable staging image digest
keyless Sigstore container signature
GitHub provenance and SBOM attestations attached to the digest
immutable generic evidence package
promotion-time signature, provenance, and SBOM verification
deployment-time evidence verification before Cloud Run mutation
lane-specific evidence repositories and IAM documentation
```

The evidence package is versioned by the image digest and is immutable. A
promotion or production deployment fails closed when the package, signature,
provenance, or SPDX attestation is missing, mismatched, or unverifiable.

## GOV-53: Superseded broker release candidate 1.0.0

Status: review pending.

The protected candidate `release/1.0.0` is a pre-delivery candidate, not a
released version:

```text
release ref:
release/1.0.0

candidate tip and current merge base:
f626600f3ab128e1a7dab271de8b1eab3bb071d3

main-only commits:
GOV-31 main baseline
GOV-33 controlled artifact promotion

published GitHub releases and v* tags:
none
```

The candidate is `superseded-before-delivery`, not `not-required`. The latter
is reserved for a post-delivery reconciliation with no effective delta.

No retrospective promotion, tag, artifact publication, signature, attestation,
or GitHub Release may be created for `release/1.0.0`. Its successor will be a
new release line with a new version; the existing ref remains retained until
the successor completes delivery and its retention and controlled cleanup
conditions are recorded.

## GOV-52: Release credential verification and hotfix delivery source gates

Status: integrated into develop through PR #22.

Scope:

```text
fixed release-credential-verification credential profile
fixed hotfix-delivery credential profile
main-bound immutable production deployment workflows
main-bound immutable artifact-promotion workflows
same-package profile and runtime-wiring whitebox tests
workflow contracts and lane-specific GCP documentation
```

The profiles are server-side only. A caller cannot request or extend GitHub
App permissions. The new verification and hotfix delivery App, Secret Manager,
WIF, runtime, invoker, deployer, promoter, Artifact Registry, Cloud Run, and
GitHub Environment boundaries remain external prerequisites until provisioned
and verified.

## GOV-43: Hotfix propagation publisher boundary

Status: review pending.

Scope:

```text
fixed hotfix-propagation-publisher credential profile
separate main-bound publisher deployment workflow
separate immutable staging-to-publisher artifact promotion workflow
workflow-contract tests
separate App, Secret, WIF, runtime, invoker, deployer, promoter,
environment, Artifact Registry, and Cloud Run boundaries
```

The publisher may create only provenance-validated hotfix-propagation
candidates and their pull requests. It receives no Actions, Workflows,
Administration, Secrets, or Ruleset-bypass permission.

## GOV-28: Reconciliation publisher credential profiles

Status: review pending.

Evidence:

```text
fixed server-side credential profiles
reconciliation publisher requests no Actions permission
whitebox tests
100% statement coverage
PR #4
```

## GOV-29: Broker CI and Supply-Chain contract

Status: in implementation.

Scope:

```text
Linux AMD64 required quality gate
CodeQL and dependency admission workflows
Dependabot intake for Go modules, GitHub Actions, and Docker base images
broker-specific source quality command
Windows/macOS required runtime check removal
Ruleset check-name alignment
architecture and operations documentation
```

External blockers:

```text
approved internal Go proxy
hermetic Go 1.26.5 build image
artifact evidence registry and production identity
```

Those blockers prevent a full Production Supply-Chain-Fortress claim. They do
not authorize public-network fallback or unverified release delivery.

## GOV-30: Main-bound Broker deployment topology

Status: Scratch exploration.

Scope:

```text
separate staging environment and resources from develop
main-bound release-automation production deployment
main-bound reconciliation publisher production deployment
immutable production digest-only deployment
retirement of the develop-bound direct deployment workflow
```

External prerequisites:

```text
GitHub environments with reviewer and branch restrictions
separate WIF, deployer, runtime, and invoker identities
separate Artifact Registry repositories
separate Cloud Run services and Secret Manager boundaries
approved production image evidence lane
```
