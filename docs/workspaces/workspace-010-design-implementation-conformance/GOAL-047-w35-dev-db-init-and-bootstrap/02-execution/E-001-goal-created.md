---
id: E-001-goal-created
doc: execution-entry
status: active
parent: GOAL-047-w35-dev-db-init-and-bootstrap
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-001 · W35 立项

- **来源**：用户 2026-09-21 报告（原文见 `00-meta.md` §概述）+ P-004 裁决：
  - **结构**：在 workspace-010 **新加一个子目标**（不新开 VP）；
  - **范围**：**完整**（初始化快捷方式 + `e2e-pgset` 自举修复 + README/QUICKSTART 章节 + 启动错误可操作化）；
  - **编号与 slug**（用户确认）：`GOAL-047-w35-dev-db-init-and-bootstrap`。
- **登记出处**：`docs/vision/roadmap.md` §未决项统一登记 §一（commit `934c60f5`）。
- **复现事实（编排器一手）**：
  - API 启动：`FATAL: database "schema_ui_dev" does not exist (SQLSTATE 3D000)` → `LIFECYCLE_START_FAILED [core.auth-session]`；
  - `go run ./cmd/e2e-pgset list` → 同源 3D000（`maintenanceDSN()` 连 `DB_NAME`，无法自举）；
  - `DB_NAME=postgres go run ./cmd/e2e-pgset create schema_ui_dev|schema_ui_test` → 成功；
  - 建库后 `dev.cmd start` → `API ready (/readyz 200)` + `Web listening`，`dev.cmd stop` 干净收尾；PG 测试路径恢复（`internal/backup` / `internal/store` PG 用例 ok）。
- **产物**：本目标五件套 + 三个 ledger 目录。
- **进度评估**：`progress: 0/3`；尚未落码初始化快捷方式。
