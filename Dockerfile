# syntax=docker/dockerfile:1

# golang:1.26.5-alpine3.23 (linux/amd64)
FROM --platform=linux/amd64 golang@sha256:2005724102f45917a63e9d092fc0e4ea56ea575048ce147caad5f5f61502c365 AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/broker ./cmd/broker

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/broker /broker

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/broker"]