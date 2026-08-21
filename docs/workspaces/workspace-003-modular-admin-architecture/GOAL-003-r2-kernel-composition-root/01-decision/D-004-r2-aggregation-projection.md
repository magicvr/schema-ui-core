---
id: D-004
title: R2 migration and Profile manifest projection
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-003-r2-kernel-composition-root
version: 0.1.0
---

# D-004 R2 migration and Profile manifest projection

The migration collector validates the global ordered migration plan and is
called by the existing SQLite runner. A single persistence module may own
multiple migration versions; version and name remain globally unique.

The API publishes a deterministic same-origin Manifest assembled from the
resolved module IDs. Core pages remain in the core fragment; users, roles,
settings, and activity pages are projected only when their modules are
enabled. The exact Vite and Nginx path proxies to API, and the production Web
image removes the static fixture after build so it cannot be a silent fallback.

Schema, permission, reconcile, and full business-module migration remain
later-stage work; this R2 decision records the aggregation skeleton only.
