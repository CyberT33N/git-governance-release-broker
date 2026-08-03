# Verification contract

## Local safe verification

Developers may run source verification without production credentials:

```text
go run -mod=readonly ./cmd/build
go test ./...
go run ./cmd/check-coverage
```

The controlled build command verifies:

```text
Go formatting
go mod verify
go mod tidy -diff
go test -mod=readonly ./...
100% statement coverage
go test -race
go vet
Linux AMD64 broker build
embedded module provenance
```

Local verification must not use production GitHub App private keys, Cloud Run
invoker identities, deployment credentials, or a mutable production service.

## CI production-quality verification

The required CI check is:

```text
Quality gates (linux-amd64)
```

It additionally runs deterministic fuzz smoke tests for Broker configuration
and HTTP request boundaries, then verifies the Linux AMD64 source binary.

## External Fortress prerequisites

The following production controls are intentionally fail-closed until their
platform infrastructure is provisioned:

```text
approved internal Go proxy
hermetic pre-provisioned Go 1.26.5 build image
artifact registry for SBOM, provenance, signatures, and attestations
```

No CI or release workflow may substitute public-network fallback, automatic
toolchain download, or unapproved runtime tooling for those controls.
