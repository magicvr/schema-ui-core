---
id: r3-c2-c3-v1-v4-evidence
doc: evidence
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: recorded
---

# R3 C2/C3 V-1～V-4 证据包

## 代码定位

- Plan 过滤和 HTTP 注册：`apps/api/internal/kernel/module.go:240-247`、
  `apps/api/internal/handler/health.go:25-38`、
  `apps/api/internal/composition/composition.go:89-104`。
- Module-owned registration/schema：
  `apps/api/internal/modules/settings/module.go:14-16`、
  `apps/api/internal/modules/activity/module.go:14-16`、
  `apps/api/internal/handler/schema.go:26-72`。
- API source/change headers：`apps/api/internal/handler/manifest.go:31-36`、
  `apps/api/internal/handler/settings.go:36-38,155-160`。
- Host bridge and Branding reload：`apps/web/src/app/config-events.ts:10-20`、
  `apps/web/src/main.tsx:34-35`、
  `apps/web/src/app/App.integration.test.tsx:221-278`。
- Development fallback warning：`apps/web/src/protocol/app-manifest.ts:824-835`。
- Production static-file boundary：`apps/web/Dockerfile:17-18`、
  `apps/web/nginx.conf:10-25`。
- Recovery drill and profile matrix assertions：
  `apps/api/internal/composition/composition_test.go:86-176,177-321`。

## Result map

| Gate | Result | Evidence |
|------|--------|----------|
| V-1 | pass | API/Web tests; same-build landing and Manifest 200 |
| V-2 | pass | Plan-driven MVP/Admin route, Schema, nav and Manifest projections |
| V-3 | pass | one Web image, two API Profiles; page and endpoint matrix |
| V-4 | pass | config header/event/Branding integration; snapshot restore and data retention |

## Runtime output

Final Web image ID: `sha256:fff440c5afb29e20d220dfb008e5c1fdf5dc6af1559fd7de3ed09e3e9606916c`.

MVP: pages `data-table,form-controls,form-with-reactions,overview,roles,
search-form-table,users`; Settings/Activity Schema `404/404`; Settings/Operations
`404/404`; `readyz=200`; landing `200`.

Admin: pages `activity,data-table,form-controls,form-with-reactions,overview,roles,
search-form-table,settings,users`; Settings/Activity Schema `200/200`;
Settings/Operations `401/401`; `readyz=200`; landing `200`.

Both manifests carried `X-Schema-UI-Manifest-Source: api`. The production image
contained no static Manifest and passed `nginx -t` with the temporary `api` DNS
alias. Temporary containers/network were removed after the check.

## Boundary

This is local dirty-snapshot evidence. It does not establish CI, clean revision,
deployment, release, or VP close-out.
