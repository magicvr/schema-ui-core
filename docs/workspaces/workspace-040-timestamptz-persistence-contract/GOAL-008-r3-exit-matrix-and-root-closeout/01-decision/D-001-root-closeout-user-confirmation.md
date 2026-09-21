---
id: D-001-root-closeout-user-confirmation
doc: decision-entry
status: accepted
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-001 · 检查点 C：用户确认关门（判据 6）

## 决定的来源

- 用户 **2026-09-21** 在确认包（`03-audit/A-006` §4）中书面裁决：**确认关闭 Root**；并对 dev 环境指向选择**改为专用 dev 库**。
- 前置：退出判据矩阵（`attachments/r3d-root-exit-criteria-matrix-v0.1.md`）；self `A-001`；independent 关门审计 `A-002`（`conditional`/开放 required = 0）与定向复审 `A-005`（**`pass`**/开放 required = 0）；`F-I-101`（required，用户报告）→ `fixed` 并独立复现。

## 1. 用户裁决

| 事项 | 用户选择 | 处置 |
|------|----------|------|
| `I-041-010`（Root 关门确认，required） | **确认关闭 Root** | Root `D-020` 落盘；本目标 → `done · 3/3`；Root → `done · 3/3`；勾选六条成功标准；同步 `goal-tree` / `docs/vision/roadmap.md` / `workspace.md` |
| dev 环境指向（non-blocking） | **改为专用 dev 库** | `apps/api/configs/.env`（gitignored）`DB_NAME` → `schema_ui_dev`；在常驻实例创建该库并迁移；`PG_TEST_*` 等其余键不变 |

## 2. 关门依据（不放行未闭合项）

- 判据 1–5：满足（逐条证据见矩阵，含真实 PG/SQLite 恢复与升级路径）。
- 判据 6：**本条用户确认**。
- 跨目标开放 required = **0**；未闭合 recommended = **0**（`F-I-102`/`F-I-103` 已在 `A-006` 修复；`F-I-006` 随本条同步投影）。
- residual：无未闭合项（R1 `D-021` 的 `F-I-005` 已 `fixed`）。

## 3. 边界

- **VP-040 的波次关闭 / Vision Review 不在本决策范围**（决策层 `/vision` 动作）；本条只更新 `docs/vision/roadmap.md` 的投影事实。
- 不重开任何冻结物；不把容器矩阵当作生产就绪证据。

## 未选方案

- **暂不关门**：用户未选；`F-I-101` 已修复并独立验证，无阻断项。
- **保持 `.env` 指向共享 `postgres` 库**：用户未选（见 Root `D-020` §2）。
- **由编排器自行标 `done`**：违反判据 6 与 P-004，未采用。
