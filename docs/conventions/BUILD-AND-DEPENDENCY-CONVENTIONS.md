# Build and dependency conventions

## Go module contract

```text
go.mod and go.sum are reviewed build inputs
normal verification uses -mod=readonly
go mod tidy -diff detects unauthorized module metadata drift
go mod verify validates cached module checksums
```

## Toolchain contract

The current source-level contract requires:

```text
Go 1.26.5
GOTOOLCHAIN=local
GOFLAGS=-mod=readonly
GOVCS=*:off
```

`check-latest: true` and automatic toolchain download are prohibited in the
controlled CI lane.

## Linux artifact contract

The only current production artifact is:

```text
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
./cmd/broker
```

It is built as a non-root Cloud Run container. `.build/` is local generated
output and is not versioned.

## External Fortress boundary

The future hermetic lane must use an approved internal Go proxy, a
pre-provisioned signed build image, and an artifact evidence registry. Their
identifiers are platform-owned configuration, not repository literals or
developer credentials.
