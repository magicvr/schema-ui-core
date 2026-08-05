---
id: E-009-r6-c62-repository-ownership
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-009 · R6 C6.2 repository ownership 迁出完成

## 已发生事实

- `core.auth-session`、`core.operationlog` 与 `admin.settings` 现分别在
  `apps/api/internal/modules/{authsession,operationlog,settings}/` 持有领域类型、sentinel、
  SQL 与 repository。admin.users/admin.roles 通过 auth-session repository 复用账号与
  RBAC 域；各 repository 只消费 `TxRunner.WithTx`，不取得底层 `*sql.DB`。
- Fx composition 构造并共享 auth-session、operation-log、settings repositories；auth、
  users、roles、settings、activity handler/provider 只消费窄 repository/reader/recorder
  边界。Settings 字段级 PATCH 使用 repository 内的单语句 upsert，未提交字段不会因并发
  read-modify-write 被覆盖。
- `apps/api/internal/store/` 的生产职责已收窄为 SQLite 连接生命周期、事务边界、migration
  runner/ledger、快照/完整性与 readiness；旧 users/roles/settings/operation-log 领域实现、
  类型、sentinel 和测试 ownership 已删除。迁移、升级、恢复类集成测试仍保留在 store 包，
  但业务断言改经 module repositories 执行。
- repository 行为测试已迁至 owner 包，覆盖账号/refresh、last-admin 与 RBAC grants、
  operation-log append/list/filter/get/failure seam、settings defaults/校验/缺行 upsert/字段级
  patch；旧 store 重复领域测试删除。

## 验证

- `apps/api: go test -count=1 ./...`：通过。
- `apps/api: go vet ./...`：通过。
- `git diff --check`：通过。
- 生产 `internal/store` 领域类型、CRUD 方法与 sentinel 扫描：零命中。
- 全仓 `store.Operation` / `store.User` / `store.Role` / `store.SiteSettings` / store-owned
  domain sentinel 扫描：零命中。
- 生产 `internal/store` 导入仅剩 composition（平台 runner 生命周期）、handler health
  （Ping）与 testsupport；未发现模块 repository 回流依赖具体 store 类型。

## Git checkpoint

- `281090e` · `refactor(api): move domain repositories out of store`
- owned paths：module repositories、composition/handler/provider 接线、owner 与迁移/恢复
  测试、旧 store 领域实现删除；未暂存既有 handler 测试换行状态噪音。

## 阻塞与下一步

- 实现与 self 核验已具备关闭 Root A-010 F-001 的证据；C6.2 的 `cross` 门禁仍等待 Grok
  Build independent audit。独立意见落盘并响应全部 required 前，不勾选 C6.2。
- C6.3 / R6-I003 尚未开始，本条不构成 Schema 字节 ContributionSet 或 VP exit #1-#7
  完整取证。
