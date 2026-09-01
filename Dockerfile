# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25.12

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/gokus-ssh ./cmd/gokus-ssh

RUN mkdir -p /out/data/ssh

FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="gokus-ssh" \
      org.opencontainers.image.description="Blokus Duo over SSH" \
      org.opencontainers.image.source="https://github.com/livensmi1e/gokus" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build --chown=65532:65532 /out/gokus-ssh /usr/local/bin/gokus-ssh
COPY --from=build --chown=65532:65532 /out/data /data

ENV GOKUS_SSH_ADDRESS=0.0.0.0:23234 \
    GOKUS_SSH_HOST_KEY_PATH=/data/ssh/gokus_ed25519

VOLUME ["/data"]
EXPOSE 23234

USER 65532:65532
ENTRYPOINT ["/usr/local/bin/gokus-ssh"]
