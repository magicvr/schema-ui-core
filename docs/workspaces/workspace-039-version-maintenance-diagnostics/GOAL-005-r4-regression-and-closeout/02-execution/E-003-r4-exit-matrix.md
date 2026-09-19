---
id: E-003-r4-exit-matrix
doc: execution-entry
parent: GOAL-005-r4-regression-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-003 · R4 退出矩阵（当前事实）

| VP-039 判据 | 状态 | 证据 |
|--------------|------|------|
| 1. 分母与契约冻结 | **达成** | GOAL-002 `done · 4/4`; D-001; two freeze matrices; A-001/A-002 pass; required 0 |
| 2. 维护提示可感知 | **达成** | GOAL-003 `done · 4/4`; Go handler full pass; RuntimeBanner tests; R2 cross pass; existing operational allowlist unchanged |
| 3. 版本提示可核对 | **达成** | GOAL-004 `done · 4/4`; VersionChip permission gate; status version source; QUICKSTART URL; R3 cross pass |
| 4. 诊断摘要只读 | **达成** | existing `admin.system-monitoring` status contract; no new diagnostic page; full Go/Vitest pass |
| 5. 体验与权限 | **达成（可运行证据）** | zh/en catalog; light/dark shell tests; mvp/admin browser slices; monitoring.read/no-permission unit tests; full Vitest 123/1489 |
| 6. 范围保持 | **达成** | R4 boundary scan: no Profile default/pinned upstream/VP-012 gate changes; no VP-040/Redis/MQ/multi-instance/search |
| 7. 证据与审计 + 用户确认 | **待完成** | R4 C4 self + grok independent pending; Root/VP close requires user written confirmation |

### Browser residual

Default mvp full Playwright had 17 passed / 5 skipped / 1 failed because the existing shared fresh-seed sign-in ordering contract caused a later list-visual login helper timeout. The same test isolated passed (1/1), and admin core slice passed 7/1/0. This remains the pre-existing bounded harness residual already registered in `docs/vision/roadmap.md`; it is not a VP-039 product finding.
