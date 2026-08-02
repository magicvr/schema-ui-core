# syntax=docker/dockerfile:1
# GOAL-008 S2 / I-008-001 §3 — api multi-stage image.
# modernc.org/sqlite is a pure-Go driver, so CGO_ENABLED=0 static build works
# and no gcc/libc toolchain is needed at runtime.
#
# BuildKit cache mounts keep the Go module cache across builds, so a transient
# proxy.golang.org download failure just re-runs `go mod download` from cache.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
# Retry once: a transient proxy.golang.org drop then completes from the cache.
RUN --mount=type=cache,target=/go/pkg/mod go mod download || go mod download
COPY . .
ARG VERSION=0.1.0
ARG COMMIT=unknown
ARG BUILT_AT=unknown
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -trimpath \
    -ldflags "-X github.com/magicvr/schema-ui-core/apps/api/pkg/version.Version=${VERSION} \
              -X github.com/magicvr/schema-ui-core/apps/api/pkg/version.Commit=${COMMIT} \
              -X github.com/magicvr/schema-ui-core/apps/api/pkg/version.BuiltAt=${BUILT_AT}" \
    -o /out/schema-ui-core-api ./cmd/server

# Minimal alpine runtime; busybox provides wget for the compose healthcheck.
FROM alpine:3.20
RUN addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app/data \
    && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/schema-ui-core-api /app/schema-ui-core-api
USER app
ENV DB_PATH=/app/data/schema-ui.db \
    APP_ENV=production \
    AUTH_DEV_SESSION_ENABLED=false
EXPOSE 8080
ENTRYPOINT ["/app/schema-ui-core-api"]
