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
go test ./...
go run ./cmd/check-coverage
```

## Container build

```powershell
docker build --tag git-governance-release-broker:dev .
```

## GCP deployment

The reviewed `develop` branch contains a manual, OIDC-authenticated deployment
workflow at `.github/workflows/gcp-deploy.yml`. It first bootstraps the Docker
Artifact Registry repository and then builds an immutable image digest before
deploying the private Cloud Run service.

Deployment identity setup, required non-secret GitHub variables, and the exact
runtime IAM boundary are documented in
[`deploy/gcp/README.md`](deploy/gcp/README.md).

For Cloud Run, deploy a pinned image digest with:

- Cloud Run IAM authentication required;
- `broker-runtime` as the service identity;
- the private key mounted from Secret Manager;
- minimum instances set to zero;
- a repository-specific `release-broker-invoker` as the only Cloud Run invoker.

## Operational limitations

This initial implementation supports `github.com` release automation. GitHub
Enterprise support requires an explicit allowlisted host-to-API mapping and
must be added through a separately reviewed change.