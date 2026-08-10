# git-governance-release-broker

`git-governance-release-broker` is a narrowly scoped credential broker for
governed GitHub release automation. It holds a GitHub App private key in a
server-side secret manager and returns only short-lived, repository-bound
installation tokens to an authenticated workload.

It does not create branches, tags, releases, or pull requests. Those actions
remain in `git-governance` and the protected GitHub Actions workflows.

## Security boundary

The broker accepts only:

```text
POST /v1/github/installations/token
Authorization: Bearer <workload-identity>
Content-Type: application/json
```

```json
{
  "host": "github.com",
  "owner": "CyberT33N",
  "repository": "git-governance"
}
```

The Cloud Run deployment must require IAM authentication. The service itself
enforces a configured repository allowlist before it requests a GitHub App
installation token.

The GitHub App private key must be mounted from Secret Manager at runtime. It
must never be committed, passed as a command-line argument, stored on a
developer workstation, or written to logs.

## Runtime configuration

- `BROKER_ALLOWED_REPOSITORIES` is required and contains a comma-separated
  `host/owner/repository` allowlist.
- `BROKER_APP_ID` is required and contains the numeric GitHub App ID.
- `BROKER_APP_INSTALLATION_ID` is required and contains the numeric approved
  repository installation ID.
- `BROKER_CREDENTIAL_PROFILE` selects a fixed server-side permission profile.
  It defaults to `release-automation`. Accepted values are:
  - `release-automation`: `actions: write`, `contents: read`, and
    `pull_requests: write`;
  - `reconciliation-publisher`: `contents: write` and
    `pull_requests: write`, without an Actions permission request.
  - `hotfix-propagation-publisher`: `contents: write` and
    `pull_requests: write`, without an Actions permission request. It uses a
    separate GitHub App installation and is limited to reviewed
    hotfix-propagation candidates.
  - `release-credential-verification`: `contents: read` only. It uses a
    separate GitHub App installation and is limited to read-only verification
    of the release credential boundary.
  - `hotfix-delivery`: `actions: read`, `contents: read`, and
    `pull_requests: read`. It uses a separate GitHub App installation and is
    limited to validated main- or support-hotfix delivery evidence.
  The HTTP request never selects a profile or GitHub permission.
- `BROKER_PRIVATE_KEY_PATH` is required and contains the mounted PEM file
  path.
- `PORT` is optional and defaults to `8080`.
- `BROKER_API_BASE_URL` is optional, must use HTTPS, and defaults to
  `https://api.github.com`.
- `BROKER_REQUEST_TIMEOUT` is optional and defaults to `10s`.
- `BROKER_MAX_REQUEST_BYTES` is optional and defaults to `4096`.
- `BROKER_MIN_TOKEN_LIFETIME` is optional and defaults to `2m`.

## Local development

The broker cannot mint credentials without a GitHub App key, but its test suite
does not use a production key.

```powershell
go run -mod=readonly ./cmd/build
```

The command runs the source-level Broker quality contract and builds a Linux
AMD64 broker binary under `.build/bin/`. It does not deploy Cloud Run or use a
production credential.

See [`docs/README.md`](docs/README.md) for the architecture, verification, CI,
and Supply-Chain conventions.

## Container build

```powershell
docker build --tag git-governance-release-broker:dev .
```

## GCP deployment

The deployment topology separates staging from production:

```text
develop
→ gcp-broker-staging.yml
→ isolated staging Broker

main
→ gcp-broker-production.yml
→ immutable production digest
→ release-automation Broker

main
→ gcp-reconciliation-publisher-production.yml
→ immutable production digest
→ reconciliation publisher Broker

main
→ gcp-hotfix-propagation-publisher-production.yml
→ immutable production digest
→ hotfix propagation publisher Broker

main
→ gcp-release-credential-verification-production.yml
→ immutable production digest
→ release credential verification Broker

main
→ gcp-hotfix-delivery-production.yml
→ immutable production digest
→ hotfix delivery Broker
```

Deployment identity setup, required non-secret GitHub variables, and the exact
runtime IAM boundaries are documented in
[`deploy/gcp/README.md`](deploy/gcp/README.md).

Production Cloud Run deployment requires:

- Cloud Run IAM authentication required;
- an environment-specific runtime service identity;
- the private key mounted from Secret Manager;
- minimum instances set to zero;
- an environment-specific service account as the only Cloud Run invoker;
- a full immutable `@sha256:` image reference;
- a `main`-bound protected GitHub Environment.

Before promotion and deployment, the workflow must verify a staging-digest
evidence package containing an SPDX SBOM, keyless Sigstore signature,
provenance attestation, and SBOM attestation. The package is immutable and
stored separately from the Docker image repository; a missing, mismatched, or
unverifiable package blocks production deployment.

## Protected release lines

`release/<semver>` and `support/<major.minor>` are created only through
`.github/workflows/create-protected-line.yml` on `main`. The workflow validates
its request, derives a release line from `origin/develop` or a support line from
`origin/main`, and creates the remote protected line server-side.

Before dispatching it, configure a protected GitHub Environment named `release`
that permits only `main`, requires an independent reviewer, prevents
self-review, and disallows administrator bypass. The workflow never accepts a
caller-selected source branch.

An unpromoted candidate is not a released version. If a protected
`release/<semver>` becomes fully contained in `main` before promotion and no
immutable tag, release, or delivery evidence exists, it must be audited as
`superseded-before-delivery`. It must not receive retrospective tags, artifacts,
signatures, attestations, or a GitHub Release. A later successor uses a new
release version and the superseded ref is retained until controlled cleanup.

## Operational limitations

This initial implementation supports `github.com` release automation. GitHub
Enterprise support requires an explicit allowlisted host-to-API mapping and
must be added through a separately reviewed change.