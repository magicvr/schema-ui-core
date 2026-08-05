# syntax=docker/dockerfile:1
# GOAL-008 S2 / I-008-001 §3 — web multi-stage image (node build → nginx).
#
# Build context is the REPO ROOT (compose: context: ., dockerfile: apps/web/Dockerfile)
# because Vite/tsc resolve the @schemas alias to `../../docs/schemas` (repo root),
# which lives outside apps/web. The Vite/tsc aliases are:
#   @schemas/* -> ../../docs/schemas/*  (relative to apps/web)
# With WORKDIR /app holding the web app, that resolves to /docs/schemas.
FROM node:22 AS build
WORKDIR /app
# npm deps first for layer caching (no root package.json; web is self-contained).
COPY apps/web/package.json apps/web/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
# Web sources + the pinned protocol schemas the renderer imports at build time.
COPY apps/web .
COPY docs/schemas /docs/schemas
RUN --mount=type=cache,target=/root/.npm npm run build \
    && test ! -e dist/.well-known/schema-ui/app-manifest.json

FROM nginx:1.27-alpine
# nginx.conf serves dist/ with SPA fallback and proxies /api to the api service.
# (Context is the repo root, so the file lives under apps/web/.)
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
