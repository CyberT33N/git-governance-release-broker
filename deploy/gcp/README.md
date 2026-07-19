# GCP deployment

The deployment workflow provisions the Artifact Registry repository, builds an
immutable Linux AMD64 broker image from `develop`, and deploys Cloud Run without
making the broker publicly invokable.

## Prerequisites

Before dispatching `.github/workflows/gcp-deploy.yml`, create a dedicated
`broker-deployer` Google service account. It is distinct from:

- `broker-runtime`, which reads only the GitHub App private key at runtime;
- `release-broker-invoker`, which invokes the running private broker.

The GitHub Workload Identity Federation principal for this repository may
impersonate `broker-deployer`. The deployer needs only these narrowly scoped
permissions:

- Artifact Registry writer on `release-broker-images`;
- Cloud Run administrator for the broker service;
- Service Account User on `broker-runtime`;
- permission to attach the existing Secret Manager secret to the Cloud Run
  service.

Do not give `broker-deployer` Secret Manager Secret Accessor. The running
`broker-runtime` identity is the only identity that reads secret payloads.

## Required GitHub repository variables

Create the following GitHub Actions repository variables. They are deployment
identifiers, not secrets:

- `GCP_PROJECT_ID`: the GCP project ID.
- `GCP_REGION`: one region for Artifact Registry and Cloud Run, for example
  `europe-west3`.
- `GCP_WORKLOAD_IDENTITY_PROVIDER`: the complete Workload Identity provider
  resource name.
- `GCP_DEPLOYER_SERVICE_ACCOUNT`: `broker-deployer` service-account email.
- `GCP_SDK_VERSION`: an explicitly approved Google Cloud SDK version.
- `GITHUB_RELEASE_APP_ID`: the numeric release GitHub App ID.
- `GITHUB_RELEASE_APP_INSTALLATION_ID`: the numeric installation ID.
- `BROKER_ALLOWED_REPOSITORIES`: the exact comma-separated allowlist, for
  example `github.com/CyberT33N/git-governance`.

Create the protected GitHub Environment `gcp-broker-deployment` before
dispatching the workflow. Restrict it to `develop`; require a release
maintainer approval when the selected GitHub plan supports that control.

## Workflow operations

The workflow accepts two explicit operations:

- `bootstrap`: creates the Docker Artifact Registry repository if it is absent.
- `deploy`: builds an image from the selected `develop` commit, pushes it with
  an immutable commit tag, resolves its digest, deploys Cloud Run, mounts the
  existing secret as a file, and grants the narrow invoker identity access.

Run `bootstrap` once. Run `deploy` only from `develop` after a reviewed merge.

## Runtime resource contract

The deployment uses:

```text
Artifact Registry repository: release-broker-images
Cloud Run service: git-governance-release-broker
Runtime service account: broker-runtime
Invoker service account: release-broker-invoker
Secret: github-release-automation-private-key
Secret mount: /var/run/secrets/github-app/private-key.pem
```

Cloud Run receives `--no-allow-unauthenticated`, one maximum instance for the
test environment, and a zero minimum-instance count. The service is reachable
over HTTPS only after Cloud Run IAM validates an ID token from the approved
invoker identity.
