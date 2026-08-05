# Traceability

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

## GOV-33: Controlled production artifact promotion

Status: in implementation.

Scope:

```text
main-bound promotion of a reviewed staging image by full commit SHA
dedicated release-automation and publisher promoter identities
source commit ancestry verification against main
source and target digest equality verification
immutable production image output for deployment workflows
workflow contract coverage and deployment documentation
```

External prerequisites:

```text
separate WIF providers and promoter service accounts
Artifact Registry Reader on the staging repository
Artifact Registry Writer on only the corresponding production repository
non-secret promotion environment variables
```
