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
