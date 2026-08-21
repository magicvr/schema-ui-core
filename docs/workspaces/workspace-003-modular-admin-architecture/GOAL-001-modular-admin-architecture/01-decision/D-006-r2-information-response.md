---
id: D-006-r2-information-response
title: R2 I-004/I-005 information response
status: accepted
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# D-006 · R2 I-004/I-005 information response

## Decision

Close the Root R2 information gates I-004 and I-005 as `verified` based on
GOAL-003 C1-C4 evidence, the local verification snapshot, A-002 self review,
and Grok Build A-003 independent re-audit.

## I-004 Profile contract

The exact compiled sets are:

- `mvp`: core.server-registration, core.auth-session, core.manifest-route,
  core.navigation-capability, core.schema-render, core.operationlog,
  admin.users, admin.roles.
- `admin`: the `mvp` set plus admin.settings and admin.activity.
- `custom`: explicit modules are required.

Explicit `APP_MODULES_ENABLED` replaces the compiled Profile default. The
recorded precedence is compiled-profile-default, then modules.enabled, then
environment. Unknown profiles, missing custom modules, missing dependencies,
capability gaps, contribution conflicts, and incompatible module API ranges
fail closed.

## I-005 Manifest contract

Fragments are sorted deterministically by ModuleID and aggregate conflicts for
app identity, protocol version, pages, navigation, capabilities, and
contribution ownership fail closed. `ForModules` projects only selected Admin
modules. The API publishes a login-free exact-path endpoint with deterministic
ETag and 304 handling. Vite and Nginx proxy the endpoint to API; the
production image removes the static manifest so it cannot be a silent fallback.

## Boundaries

This decision closes information readiness for R2. It does not close I-006,
does not mark Root R2 complete, and does not establish CI, deployment, or
release acceptance. Any change to the Profile contract or aggregation boundary
requires a new decision and evidence.
