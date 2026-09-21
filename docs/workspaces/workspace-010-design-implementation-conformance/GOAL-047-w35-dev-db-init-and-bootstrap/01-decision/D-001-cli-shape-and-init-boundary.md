---
id: D-001-cli-shape-and-init-boundary
doc: decision-entry
status: accepted
parent: GOAL-047-w35-dev-db-init-and-bootstrap
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-001 · 初始化快捷方式 CLI 形态与自举边界

## 决定的来源

- 用户 2026-09-21 P-004 裁决：在 workspace-010 新增子目标，不新开 VP；范围选**完整**（快捷方式、自举修复、文档、启动错误可操作化）；slug `GOAL-047-w35-dev-db-init-and-bootstrap` 已确认。
- 用户报告与编排器复现：目标库不存在时 API ping 失败 SQLSTATE `3D000`；既有 `e2e-pgset` 也因维护连接依赖 `DB_NAME` 而无法创建第一个库。

## 选择

- 新增仓库根入口 **`dev.cmd init-db`**。
- 实际实现入口 **`apps/api/cmd/dbsetup`**：默认幂等确保 dev (`DB_NAME`) + test (`PG_TEST_DB`)；支持 `--dev-only`、`--test-only` 与额外安全库名。
- 创建/清理/验证能力抽在 `apps/api/internal/pgsetup`：按角色分离 `DB_*` 与 `PG_TEST_*` namespace；维护连接候选顺序为「目标库 → `postgres` → `template1`」；已存在返回 `exists`，不会删除或覆盖。
- `cmd/e2e-pgset` 复用同一自举包，保留 e2e 专用 `schema_ui_e2e_*` 生命周期；`create` 幂等、`drop` 缺失库安全返回 `absent`，不并入 dev/test 默认集。
- PG SQLSTATE `3D000` 启动错误追加 actionable hint：`dev.cmd init-db` / `go run ./cmd/dbsetup`；原始 error chain 与分类不变，密码不进入提示。

## 约束

- 应用启动**不自动 `CREATE DATABASE`**：避免拼错库名/连错实例时静默创建；init 是显式用户动作。
- SQLite 不需要该命令，文件由首次启动创建。
- 迁移由 API/test 自身应用，dbsetup 只建库。
- e2e 库继续由 e2e 工具自管；dev/test database 分别使用各自 namespace 的连接参数。

## 未选方案

- 把 `CREATE DATABASE` 放进 API startup：未采用，违反既有启动身份边界。
- 仅修 e2e helper、不提供普通开发快捷方式：未采用。
- 只写文档、不提供命令入口：未采用。
