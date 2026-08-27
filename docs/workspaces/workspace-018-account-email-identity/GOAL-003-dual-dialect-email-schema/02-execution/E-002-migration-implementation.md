---
id: E-002
doc: execution-entry
goal: GOAL-003-dual-dialect-email-schema
status: recorded
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · 迁移 0054 实现与验证（2026-08-24）

## 已发生事实

- `authsession/migration/migration.go`：新增 **0054 account_email_identity**（三条可移植 DDL：`email TEXT` / `email_status TEXT CHECK('pending','verified')` / `CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))`）；`ApplyPostgres: nil`（双方言同一文本）；checksum = `f9a0bc65…2b0b`（transform tag `0054:account-email-identity:v1`）。既有 DDL 字符串零改动。
- 黄金断言六处同步：identity.go head 53→54；identity_test lockedHeadExtraTables[54]；migrate_test 尾链+计数 ×3；operations_test、restart_test 尾断言；migrate_test 冻结 checksum 目录追加 0054 行。
- 新增 `store/migrate_0054_test.go`（升级路径 unbound 落地 / 全新库对象形状 / 大小写唯一拒绝 / 多 NULL 共存 / CHECK 拒绝）。
- 修复一轮测试自身键名笔误（PRAGMA 列名小写）；无产品代码返工。

## 验证矩阵

| 面 | 结果 |
|----|------|
| `go build ./...` | 干净 |
| `go test ./internal/store/ -count=1`（SQLite 全量含黄金断言） | ok ~32s |
| `go test ./internal/composition/ -count=1` | ok ~17s |
| authsession / kernel / systemdata | ok |
| PostgreSQL 17 集成（configs/.env → LAN PG）：`-run Postgres -v` | **15/15 PASS**，含 TestFullCatalogPostgresBootstrapIntegration（全 catalog 含 0054 实跑）与 TestAuthsessionPostgresApplyIntegration |
| 复跑稳定性 | EXIT=0 · PASS 15 · 无 FAIL（首轮一次瞬时非复现失败，复跑两轮干净） |

## 证据

| 主张 | 路径 |
|------|------|
| 迁移实现 | `apps/api/internal/modules/authsession/migration/migration.go` |
| 冻结目录行 | `apps/api/internal/store/migrate_test.go`（want 表 0054 行） |
| 专项测试 | `apps/api/internal/store/migrate_0054_test.go` |

## 未做

- 未写绑定流仓储语义（R3）；未动 Web 面；independent 审计待执行（实现后门禁）。
