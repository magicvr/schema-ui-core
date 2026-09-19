---
id: E-001-r1-recon
doc: execution-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R1 只读侦察（I-039-001～003）

## 事实

2026-09-19 立项 `GOAL-002` 并完成只读侦察，未改 `apps/**`。

三份报告：

- `attachments/r1-recon-i-039-001-version-identity.md`
- `attachments/r1-recon-i-039-002-mode-projection.md`
- `attachments/r1-recon-i-039-003-diagnostics.md`

关键发现（非决策）：`runtime.mode=maintenance` 时 Host bootstrap 在 availability-gate 终态为 `MAINTENANCE`，Admin Shell **不会加载**；Shell 横幅只能覆盖 Shell 仍加载的模式（`degraded` / `read-only` → Host `degraded` → `READY_DEGRADED`）。

## 非事实

尚未冻结 C1～C3；P-004 待用户裁决。
