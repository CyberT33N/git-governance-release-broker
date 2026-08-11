# syntax=docker/dockerfile:1

ARG BUILDER_IMAGE
FROM --platform=linux/amd64 ${BUILDER_IMAGE} AS build

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