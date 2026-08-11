# GCP deployment topology

## Deployment lanes

The Broker has six separate deployment lanes:

```text
develop
→ gcp-broker-staging
→ isolated staging image repository
→ isolated staging Cloud Run service

main
→ gcp-broker-production
→ immutable digest from the production image repository
→ production release-automation broker

main
→ gcp-release-credential-verification-deployment
→ immutable digest from the dedicated verification image repository
→ production release-credential-verification broker

main
→ gcp-hotfix-delivery-deployment
→ immutable digest from the dedicated hotfix delivery image repository
→ production hotfix-delivery broker

main
→ gcp-reconciliation-publisher-deployment
→ immutable digest from the publisher production repository
→ production reconciliation-publisher broker

main
→ gcp-hotfix-propagation-publisher-deployment
→ immutable digest from the dedicated publisher production repository
→ production hotfix-propagation-publisher broker
```

The former develop-bound `gcp-broker-deployment` workflow is retired. A
production deployment must never build and deploy unpromoted `develop` source.

## GitHub environments

Create these protected environments before dispatching the workflows:

```text
gcp-broker-staging
→ selected branch: develop
→ required reviewers
→ prevent self-review
→ administrator bypass disabled

gcp-broker-production
→ selected branch: main
→ required reviewers
→ prevent self-review
→ administrator bypass disabled

gcp-release-credential-verification-deployment
→ selected branch: main
→ required reviewers
→ prevent self-review
→ administrator bypass disabled

gcp-hotfix-delivery-deployment
→ selected branch: main
→ required reviewers
→ prevent self-review
→ administrator bypass disabled

gcp-reconciliation-publisher-deployment
→ selected branch: main
→ required reviewers
→ prevent self-review
→ administrator bypass disabled

gcp-hotfix-propagation-publisher-deployment
→ selected branch: main
→ required reviewers
→ prevent self-review
→ administrator bypass disabled
```

## Staging resources

The `gcp-broker-staging.yml` workflow requires environment variables:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_STAGING_WORKLOAD_IDENTITY_PROVIDER
GCP_STAGING_DEPLOYER_SERVICE_ACCOUNT
GCP_STAGING_ARTIFACT_REPOSITORY
GCP_STAGING_EVIDENCE_ARTIFACT_REPOSITORY
GCP_STAGING_BROKER_SERVICE
GCP_STAGING_RUNTIME_SERVICE_ACCOUNT
GCP_STAGING_INVOKER_SERVICE_ACCOUNT
GCP_STAGING_BROKER_SECRET
GCP_STAGING_BROKER_APP_ID
GCP_STAGING_BROKER_APP_INSTALLATION_ID
GCP_STAGING_BROKER_ALLOWED_REPOSITORIES
```

The staging deployer has Artifact Registry writer for the staging image and
staging evidence repositories, Cloud Run deployment, and Service Account User
permissions only for staging resources. The staging runtime identity reads only
its staging GitHub App key secret.

## Production release-automation resources

The `gcp-broker-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_PRODUCTION_WORKLOAD_IDENTITY_PROVIDER
GCP_PRODUCTION_DEPLOYER_SERVICE_ACCOUNT
GCP_PRODUCTION_ARTIFACT_REPOSITORY
GCP_PRODUCTION_EVIDENCE_ARTIFACT_REPOSITORY
GCP_PRODUCTION_BROKER_SERVICE
GCP_PRODUCTION_RUNTIME_SERVICE_ACCOUNT
GCP_PRODUCTION_INVOKER_SERVICE_ACCOUNT
GCP_PRODUCTION_BROKER_SECRET
GCP_PRODUCTION_BROKER_APP_ID
GCP_PRODUCTION_BROKER_APP_INSTALLATION_ID
GCP_PRODUCTION_BROKER_ALLOWED_REPOSITORIES
```

Production accepts only a full immutable image reference:

```text
<region>-docker.pkg.dev/<project>/<production-repository>/broker@sha256:<64-lowercase-hex-digest>
```

It never builds from `develop`, accepts a mutable tag, or creates a public
Cloud Run service.

## Release credential verification resources

The `gcp-release-credential-verification-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_RELEASE_CREDENTIAL_VERIFICATION_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER
GCP_RELEASE_CREDENTIAL_VERIFICATION_DEPLOYMENT_SERVICE_ACCOUNT
GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_REPOSITORY
GCP_RELEASE_CREDENTIAL_VERIFICATION_EVIDENCE_ARTIFACT_REPOSITORY
GCP_RELEASE_CREDENTIAL_VERIFICATION_BROKER_SERVICE
GCP_RELEASE_CREDENTIAL_VERIFICATION_RUNTIME_SERVICE_ACCOUNT
GCP_RELEASE_CREDENTIAL_VERIFICATION_INVOKER_SERVICE_ACCOUNT
GCP_RELEASE_CREDENTIAL_VERIFICATION_BROKER_SECRET
GCP_RELEASE_CREDENTIAL_VERIFICATION_BROKER_APP_ID
GCP_RELEASE_CREDENTIAL_VERIFICATION_BROKER_APP_INSTALLATION_ID
GCP_RELEASE_CREDENTIAL_VERIFICATION_BROKER_ALLOWED_REPOSITORIES
```

The broker always runs:

```text
BROKER_CREDENTIAL_PROFILE=release-credential-verification
```

Its GitHub App, Secret Manager key, Cloud Run service, WIF provider, runtime,
invoker, deployer, promoter, environment, and Artifact Registry repository are
separate from every mutating release, reconciliation, and hotfix lane. The App
may request only `contents: read` for the approved repository.

## Hotfix delivery resources

The `gcp-hotfix-delivery-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_HOTFIX_DELIVERY_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER
GCP_HOTFIX_DELIVERY_DEPLOYMENT_SERVICE_ACCOUNT
GCP_HOTFIX_DELIVERY_ARTIFACT_REPOSITORY
GCP_HOTFIX_DELIVERY_EVIDENCE_ARTIFACT_REPOSITORY
GCP_HOTFIX_DELIVERY_BROKER_SERVICE
GCP_HOTFIX_DELIVERY_RUNTIME_SERVICE_ACCOUNT
GCP_HOTFIX_DELIVERY_INVOKER_SERVICE_ACCOUNT
GCP_HOTFIX_DELIVERY_BROKER_SECRET
GCP_HOTFIX_DELIVERY_BROKER_APP_ID
GCP_HOTFIX_DELIVERY_BROKER_APP_INSTALLATION_ID
GCP_HOTFIX_DELIVERY_BROKER_ALLOWED_REPOSITORIES
```

The broker always runs:

```text
BROKER_CREDENTIAL_PROFILE=hotfix-delivery
```

Its GitHub App, Secret Manager key, Cloud Run service, WIF provider, runtime,
invoker, deployer, promoter, environment, and Artifact Registry repository are
separate from regular release delivery and all publisher lanes. The App may
request only `actions: read`, `contents: read`, and `pull_requests: read` for
the approved repository.

## Reconciliation publisher resources

The `gcp-reconciliation-publisher-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_RECONCILIATION_PUBLISHER_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER
GCP_RECONCILIATION_PUBLISHER_DEPLOYMENT_SERVICE_ACCOUNT
GCP_RECONCILIATION_PUBLISHER_ARTIFACT_REPOSITORY
GCP_RECONCILIATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY
GCP_RECONCILIATION_PUBLISHER_BROKER_SERVICE
GCP_RECONCILIATION_PUBLISHER_RUNTIME_SERVICE_ACCOUNT
GCP_RECONCILIATION_PUBLISHER_INVOKER_SERVICE_ACCOUNT
GCP_RECONCILIATION_PUBLISHER_BROKER_SECRET
GCP_RECONCILIATION_PUBLISHER_BROKER_APP_ID
GCP_RECONCILIATION_PUBLISHER_BROKER_APP_INSTALLATION_ID
GCP_RECONCILIATION_PUBLISHER_BROKER_ALLOWED_REPOSITORIES
```

The publisher broker always runs:

```text
BROKER_CREDENTIAL_PROFILE=reconciliation-publisher
```

Its GitHub App is separate from release automation and has only the repository
and permissions required to publish a provenance-validated reconciliation
candidate.

## Hotfix propagation publisher resources

The `gcp-hotfix-propagation-publisher-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_HOTFIX_PROPAGATION_PUBLISHER_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER
GCP_HOTFIX_PROPAGATION_PUBLISHER_DEPLOYMENT_SERVICE_ACCOUNT
GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_REPOSITORY
GCP_HOTFIX_PROPAGATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY
GCP_HOTFIX_PROPAGATION_PUBLISHER_BROKER_SERVICE
GCP_HOTFIX_PROPAGATION_PUBLISHER_RUNTIME_SERVICE_ACCOUNT
GCP_HOTFIX_PROPAGATION_PUBLISHER_INVOKER_SERVICE_ACCOUNT
GCP_HOTFIX_PROPAGATION_PUBLISHER_BROKER_SECRET
GCP_HOTFIX_PROPAGATION_PUBLISHER_BROKER_APP_ID
GCP_HOTFIX_PROPAGATION_PUBLISHER_BROKER_APP_INSTALLATION_ID
GCP_HOTFIX_PROPAGATION_PUBLISHER_BROKER_ALLOWED_REPOSITORIES
```

The broker always runs:

```text
BROKER_CREDENTIAL_PROFILE=hotfix-propagation-publisher
```

Its GitHub App, Secret Manager key, Cloud Run service, WIF provider, runtime,
invoker, deployer, promoter, environment, and Artifact Registry repository are
separate from release automation and reconciliation publication. The App may
request only `contents: write` and `pull_requests: write` for the approved
repository; it receives no Actions, Workflows, Administration, Secrets, or
Ruleset-bypass permission.

## Hotfix publisher artifact promotion

`gcp-hotfix-propagation-publisher-promotion.yml` is main-bound and copies a
reviewed staging digest into the dedicated hotfix publisher repository without
rebuilding it. It requires:

```text
GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_PROMOTION_WIF_PROVIDER
GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_PROMOTER_SERVICE_ACCOUNT
GCP_HOTFIX_PROPAGATION_PUBLISHER_SOURCE_ARTIFACT_REPOSITORY=release-broker-staging-images
GCP_HOTFIX_PROPAGATION_PUBLISHER_ARTIFACT_REPOSITORY
GCP_HOTFIX_PROPAGATION_PUBLISHER_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY
GCP_HOTFIX_PROPAGATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY
```

The promoter may read only the staging repository and write only the dedicated
hotfix publisher repository. It receives no Cloud Run, Secret Manager, Service
Account User, or GitHub App permission.

## Release credential verification artifact promotion

`gcp-release-credential-verification-promotion.yml` is main-bound and copies
a reviewed staging digest into the dedicated verification repository without
rebuilding it. It requires:

```text
GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_PROMOTION_WIF_PROVIDER
GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_PROMOTER_SERVICE_ACCOUNT
GCP_RELEASE_CREDENTIAL_VERIFICATION_SOURCE_ARTIFACT_REPOSITORY=release-broker-staging-images
GCP_RELEASE_CREDENTIAL_VERIFICATION_ARTIFACT_REPOSITORY
GCP_RELEASE_CREDENTIAL_VERIFICATION_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY
GCP_RELEASE_CREDENTIAL_VERIFICATION_EVIDENCE_ARTIFACT_REPOSITORY
```

The promoter may read only the staging repository and write only the dedicated
verification repository. It receives no Cloud Run, Secret Manager, Service
Account User, or GitHub App permission.

## Hotfix delivery artifact promotion

`gcp-hotfix-delivery-promotion.yml` is main-bound and copies a reviewed
staging digest into the dedicated hotfix delivery repository without rebuilding
it. It requires:

```text
GCP_HOTFIX_DELIVERY_ARTIFACT_PROMOTION_WIF_PROVIDER
GCP_HOTFIX_DELIVERY_ARTIFACT_PROMOTER_SERVICE_ACCOUNT
GCP_HOTFIX_DELIVERY_SOURCE_ARTIFACT_REPOSITORY=release-broker-staging-images
GCP_HOTFIX_DELIVERY_ARTIFACT_REPOSITORY
GCP_HOTFIX_DELIVERY_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY
GCP_HOTFIX_DELIVERY_EVIDENCE_ARTIFACT_REPOSITORY
```

The promoter may read only the staging repository and write only the dedicated
hotfix delivery repository. It receives no Cloud Run, Secret Manager, Service
Account User, or GitHub App permission.

## Release-automation artifact promotion

`gcp-broker-production-promotion.yml` copies a verified staging digest into the
release-automation production repository without rebuilding it. It requires:

```text
GCP_PRODUCTION_ARTIFACT_PROMOTION_WIF_PROVIDER
GCP_PRODUCTION_ARTIFACT_PROMOTER_SERVICE_ACCOUNT
GCP_PRODUCTION_SOURCE_ARTIFACT_REPOSITORY=release-broker-staging-images
GCP_PRODUCTION_ARTIFACT_REPOSITORY
GCP_PRODUCTION_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY
GCP_PRODUCTION_EVIDENCE_ARTIFACT_REPOSITORY
```

## Reconciliation publisher artifact promotion

`gcp-reconciliation-publisher-promotion.yml` copies a verified staging digest
into the dedicated reconciliation publisher repository without rebuilding it.
It requires:

```text
GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTION_WIF_PROVIDER
GCP_RECONCILIATION_PUBLISHER_ARTIFACT_PROMOTER_SERVICE_ACCOUNT
GCP_RECONCILIATION_PUBLISHER_SOURCE_ARTIFACT_REPOSITORY=release-broker-staging-images
GCP_RECONCILIATION_PUBLISHER_ARTIFACT_REPOSITORY
GCP_RECONCILIATION_PUBLISHER_SOURCE_EVIDENCE_ARTIFACT_REPOSITORY
GCP_RECONCILIATION_PUBLISHER_EVIDENCE_ARTIFACT_REPOSITORY
```

## Immutable evidence contract

The staging workflow signs the immutable image digest with GitHub OIDC,
generates SPDX SBOM and provenance attestations, and stores an evidence package
in a generic Artifact Registry repository. The package version is the
lower-case image digest without its `sha256:` prefix and contains:

```text
manifest.json
broker.spdx.json
signature.json
provenance.intoto.jsonl
sbom.intoto.jsonl
```

`manifest.json` is schema version 2 and models a subject graph rooted at the
immutable artifact digest:

```text
Source
Dependency Resolution
Build
Artifact
Promotion
Deployment
Operation
```

The staging graph materializes Source through Artifact bindings and marks
Promotion and Deployment as `not-recorded` plus Operation as `not-evaluated`.
Those states are explicit absence-of-evidence markers, not successful lifecycle
claims.

Promotion verifies the complete staging graph, signature, provenance, and SPDX
attestation before copying the image digest. It then writes a lane-specific
`promotion.json` that binds the staging manifest hash, source digest, target
digest, promotion workflow, promoter identity, and functional target lane into
the immutable target evidence package. Production deployment requires that
promotion subject and re-verifies the manifest, payload hashes, staging
signature, provenance, and SBOM attestation before Cloud Run mutation.

The evidence IAM model is:

```text
staging deployer
→ Artifact Registry Writer only on staging Docker and generic evidence repositories

lane promoter
→ Artifact Registry Reader only on staging Docker and generic evidence repositories
→ Artifact Registry Writer only on its lane Docker and generic evidence repositories
→ no Cloud Run, Secret Manager, or Service Account User permission

lane deployer
→ Artifact Registry Reader only on its lane Docker and generic evidence repositories
→ Artifact Registry Reader on staging Docker for signature and attestation verification
→ Cloud Run deployment and Service Account User only for its lane runtime

runtime and invoker
→ no Artifact Registry evidence access
```

## External Fortress prerequisites

The approved Go proxy, hermetic build image, production image promotion,
SBOM, provenance, signature, and attestation registry are external platform
controls. These workflows fail closed on missing environment variables but do
not claim those controls exist until the corresponding platform resources are
provisioned and verified.
