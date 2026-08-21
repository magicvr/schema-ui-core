---
title: I-007-004 · S6 重启保持与端到端验收协议
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-007-r4-schema-crud
version: 0.1.0
related_info: I-007-004
related_decision: D-007
---

# I-007-004 · S6 重启、迁移与端到端回归验收协议

> **结论**：本附件与 D-007 关闭 `I-007-004`。以下协议定义 S6 的机器可重复验收方式：**如何隔离数据库、固定操作序列、断言持久化结果并清理测试状态**，以及迁移/seed 重跑、失败路径与 API/Web 回归的覆盖口径。它是 S6 验收输入，**不是**已完成的验收证据；验收证据落在测试与 02-execution。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`；事实来自现有 store/handler 重启与迁移测试、`cmd/server` 启动路径与 Playwright E2E 配置。

## 1. 目标与验收范围

S6 成功标准（00-meta）：**自动证明 create/update/delete 后服务重启，list/detail 结果符合持久化预期；覆盖 migration/seed 重复执行、关键失败路径及 API/Web 回归。**

本协议把「服务重启」界定为两层，均需机器可重复证据：

| 层 | 定义 | 证据载体 |
|----|------|----------|
| **L1 · HTTP 层重启** | 同一 SQLite 文件上，完整 handler/auth 栈的 store 关闭→重开（进程重启的持久化边界），全 HTTP CRUD→重启→list/detail | `apps/api/internal/handler/records_restart_test.go`（本轮新增） |
| **L2 · 进程级重启** | 真实 `cmd/server` OS 进程终止→以同一 `DB_PATH` 重启，跨进程验证 | `apps/api/cmd/server/server_restart_test.go`（本轮新增） |

L2 是 A-003/A-004 明确标记的缺口（此前仅 store close/reopen 单测级）；本轮以真实子进程关闭并补齐。既有 store 级证据继续保留：`TestRecordsPersistAcrossRestart`（records 仓库层）、`TestRestartPersistence`（身份/refresh/RBAC/菜单/迁移台账）、`TestRecordsSeedIdempotentAcrossOpens`（seed 幂等）。

## 2. 数据库隔离

- 每个 S6 运行使用**全新的临时 DB 路径**（Go `t.TempDir()`），**绝不**复用默认 `./data/schema-ui.db` 或任何共享文件。
- 进程级测试通过环境变量 `DB_PATH=<temp>/restart.db` 注入；HTTP 层测试直接向 `store.Open(path,…)` 传 temp 路径。
- 端口隔离：进程级测试用 `net.Listen("127.0.0.1:0")` 取空闲端口，两个阶段各取一个，避免与既有服务/测试冲突。
- 环境显式固定：`ADMIN_INITIAL_PASSWORD=admin`、`AUTH_JWT_SECRET=test-secret`、`AUTH_DEV_SESSION_ENABLED=false`（真实会话路径，非 dev-session 兜底）、`APP_ENV=development`。

## 3. 固定操作序列（每轮）

1. **Phase 1（server #1）**：临时 DB 上启动（迁移应用到 ledger `{1,2,3}`，空表 seed admin + 8 条 records `rec-1…rec-8`）。
2. 以 admin 登录（`POST /api/auth/login`，admin / seed 密码）获得 Bearer access token。
3. 固定 CRUD 序列（I-007-001 契约）：
   - **create**：`POST /api/records` `{name,status,owner}` → 201，记录新 `id`；
   - **update**：`PATCH /api/records/rec-1`（改名）→ 200，记录响应 `updatedAt`（毫秒 RFC3339）；
   - **delete**：`DELETE /api/records/rec-2` → 204。
4. 关闭 Phase 1（HTTP 层 = store `Close()`；进程级 = 终止 OS 进程）。
5. **Phase 2（server #2，同一 DB_PATH / 同一文件）**：重新启动并等待 `/healthz`。
6. 重新登录，执行：
   - `GET /api/records`（pageSize 覆盖全表）→ 断言 create 的新行**存在**、update 的 rec-1 名称**已变**、delete 的 rec-2 **不存在**；总数 = 8（8 seed − 1 delete + 1 create，证明未重 seed 复活）；
   - `GET /api/records/{newID}` 与 `GET /api/records/rec-1` → 断言 detail 字段与 `updatedAt` 与 Phase 1 记录的**毫秒精确一致**（持久化往返）。

## 4. 持久化断言（对照期望态）

| 断言 | 期望 |
|------|------|
| create 新行 | Phase 2 list/detail 含该 id；`name`/`status`/`owner` 与 POST body 一致 |
| update 生效 | Phase 2 rec-1 `name` = 新名；`updatedAt` = Phase 1 PATCH 响应值（毫秒格式，`2006-01-02T15:04:05.000Z07:00`） |
| delete 生效 | Phase 2 无 rec-2（不复活） |
| 总数 | `total = 8`（非空表 seedRecords 跳过，未补回已删行） |
| 迁移台账 | ledger 仍为 `{1,2,3}`（不重跑；store 级 `TestRestartPersistence` 承担） |
| seed 幂等 | `user_roles` admin 行数不变（store 级承担）；records 表非空不重 seed（本轮 HTTP 总数断言） |

## 5. 迁移/seed 重跑断言

- 迁移 ledger 不重跑、seed 不重复：由既有 store 级 `TestRestartPersistence`（ledger `{1,2,3}` + `user_roles` 计数 + admin 密码不被覆盖）与 `TestRecordsSeedIdempotentAcrossOpens` 承担；HTTP 层以「总数不变、已删行不复活」补足 records 侧。
- 空表才 seed：`TestSeedRecordsEmptyTable` / `TestSeedRecordsSkipsNonEmpty` 已覆盖（S3）。

## 6. 关键失败路径

| 路径 | 覆盖 |
|------|------|
| checksum 漂移 fail closed | `TestMigrateFailClosedRecordsChecksumDrift`（T-DB-04） |
| 非空库升级前一致性快照 | `snapshotBeforePending` + `TestMigrateExistingV2ToV3` / `TestRestorePreV0002Snapshot` |
| 非开发环境缺 `ADMIN_INITIAL_PASSWORD` 启动 fail closed | `cmd/server/main.go` `resolveSeedHash`（`TestServerProcessRestart…` 不触碰；该路径属启动守卫，注释留痕） |
| 匿名 `401` / 缺权限 `403` | `T-API-08/09`（handler）+ browser E2E `shell.spec.ts` 401 断言 |
| 迁移失败事务回滚 | `TestMigrateFailClosed…` 与迁移 runner 不变量（S3 已覆盖） |

## 7. 清理

- 临时 DB 与其 `pre-v0003-*.sqlite` 快照随 `t.TempDir()` 自动清理；测试不写入默认 `./data/`。
- 进程级测试结束后 `Process.Kill()` + `Wait()`，不留孤儿进程。

## 8. 回归命令（S6 验收最低集合）

```bash
cd apps/api && go test ./...          # L1 + L2 + 全部 API 回归
cd apps/web && npm test               # web vitest 全量回归（含 T-UI-01～10）
cd apps/web && npm run test:e2e       # browser E2E（可选；需真实服务，验证登录/401/写路径）
```

## 9. 证据索引

- 新增：`apps/api/internal/handler/records_restart_test.go`（L1）、`apps/api/cmd/server/server_restart_test.go`（L2）。
- 既有：`store/records_test.go`（`TestRecordsPersistAcrossRestart`/`TestSeedRecords*`）、`store/restart_test.go`（`TestRestartPersistence`/`TestRestorePreV0002Snapshot`）、`store/migrate_test.go`、`handler/records_test.go`（T-API-08/09）、`web/e2e/shell.spec.ts`。
- 契约：I-007-001（毫秒 `updatedAt`、POST/PATCH/DELETE 语义）、I-007-002（迁移/seed/repository）。
