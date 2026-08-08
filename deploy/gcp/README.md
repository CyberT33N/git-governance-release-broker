# GCP deployment topology

## Deployment lanes

The Broker has four separate deployment lanes:

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
GCP_STAGING_BROKER_SERVICE
GCP_STAGING_RUNTIME_SERVICE_ACCOUNT
GCP_STAGING_INVOKER_SERVICE_ACCOUNT
GCP_STAGING_BROKER_SECRET
GCP_STAGING_BROKER_APP_ID
GCP_STAGING_BROKER_APP_INSTALLATION_ID
GCP_STAGING_BROKER_ALLOWED_REPOSITORIES
```

The staging deployer has Artifact Registry writer, Cloud Run deployment, and
Service Account User permissions only for staging resources. The staging
runtime identity reads only its staging GitHub App key secret.

## Production release-automation resources

The `gcp-broker-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_PRODUCTION_WORKLOAD_IDENTITY_PROVIDER
GCP_PRODUCTION_DEPLOYER_SERVICE_ACCOUNT
GCP_PRODUCTION_ARTIFACT_REPOSITORY
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

## Reconciliation publisher resources

The `gcp-reconciliation-publisher-production.yml` workflow requires:

```text
GCP_PROJECT_ID
GCP_REGION
GCP_SDK_VERSION
GCP_RECONCILIATION_PUBLISHER_DEPLOYMENT_WORKLOAD_IDENTITY_PROVIDER
GCP_RECONCILIATION_PUBLISHER_DEPLOYMENT_SERVICE_ACCOUNT
GCP_RECONCILIATION_PUBLISHER_ARTIFACT_REPOSITORY
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
```

The promoter may read only the staging repository and write only the dedicated
hotfix publisher repository. It receives no Cloud Run, Secret Manager, Service
Account User, or GitHub App permission.

## External Fortress prerequisites

The approved Go proxy, hermetic build image, production image promotion,
SBOM, provenance, signature, and attestation registry are external platform
controls. These workflows fail closed on missing environment variables but do
not claim those controls exist until the corresponding platform resources are
provisioned and verified.
