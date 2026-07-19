# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang@sha256:73f9732658b30852522ee5ebe698daa27e1829add9a70ff4f4a828409f8d0a99 AS build # golang:1.26.5-alpine3.23
WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/broker ./cmd/broker

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/broker /broker

USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/broker"]
