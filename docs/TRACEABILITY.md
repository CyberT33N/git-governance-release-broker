# Traceability

## GOV-52: Release credential verification and hotfix delivery source gates

Status: in implementation.

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
