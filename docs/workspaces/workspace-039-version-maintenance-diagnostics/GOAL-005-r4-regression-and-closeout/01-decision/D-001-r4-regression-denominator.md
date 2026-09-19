---
id: D-001-r4-regression-denominator
doc: decision-entry
parent: GOAL-005-r4-regression-and-closeout
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · R4 回归分母与关门边界

## 纳入

- Existing R2/R3 targeted tests and full Go/Web suites.
- `mvp` / `admin` / `demo` profile evidence where a surface exists; custom profile boundary via existing Profile/Manifest tests.
- `monitoring.read` vs absent permission for VersionChip fetch/render.
- runtime modes normal/maintenance/degraded/read-only for RuntimeBanner and bootstrap producer tests.
- zh-CN/en-US and light/dark component/theme tests.
- Existing Playwright suites that can run against the managed local app; command/result recorded, including skips or environment blockers.

## Exclude

- No new feature implementation in R4.
- No Redis/MQ/multi-instance/search/timestamptz/file scanning.
- No change to Host consumer protocol fixtures or VP-012 operational gate semantics.
- No silent Root/VP close: user written confirmation remains required.
