---
title: I-011-003 · 双语义资源 Schema 接入验收矩阵
status: active
doc_type: contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-011-s4-semantic-admin-resources
version: 0.2.0
related_info: I-011-003
related_decision: D-004
supersedes: null
---

# I-011-003 · 双资源验收矩阵（冻结）

> **性质**：回答「双资源验收如何证明前端 Schema-only 接入以及 fresh fork / 升级 / 重启 / 401/403 的完整边界」。本 v0.2.0 已由 GOAL-011 **D-004** 冻结，`I-011-003` 已置为 `verified`，解除 S4 集成验收的信息门禁。
> **依据**：S2（I-011-001 v0.2.0 实施事实）、S3（I-011-002 v0.2.0 records 退场）、A-007 F-001～F-003 的 fixed 修复与 2026-08-03 本机预检证据。
> **不是**：S5 全量回归与关门审计（属 S5）；Root A-002 F-002-001 关闭证据（属 GOAL-010 S5）。

## 1. 目的与范围

- 目标：证明 users 与 roles 两语义资源的代表性列表/CRUD 页面**仅通过 Schema（fixture + manifest）接入**，前端 Renderer 主路径无修改；并覆盖 fresh fork、既有数据库升级、重启持久化与权限失败路径的完整边界。
- 范围：验收矩阵（§2）、Renderer diff 边界（§3）、信息门禁关闭与 S4 完成定义（§4/§5）、证据索引（§7）。本契约冻结不自动完成 S4。

## 2. 验收矩阵

| 维度 | 验收口径 | 证据（已实施） |
|------|----------|----------------|
| **fresh fork（新库）** | 全新 DB 经 0001～0006 迁移后：**无 `records` 表**；users/roles/permissions/menu 就绪；`operation_log` event CHECK 含 users/roles 事件；users/roles 五路由可用 | `store/migrate_test.go TestMigrateFreshDB`（无 records 表 + users.create 事件可写）；handler/users_test.go + roles_test.go（fresh env 全 CRUD） |
| **既有库升级（in-place）** | 0004 态库升级到 0006：records 表 DROP、records 权限/菜单行清理、`pre-v0005`+`pre-v0006` 快照存在（per-pending 快照）、operation_log 既有行保留；升级表接受 users/roles 事件且关闭/重开后仍保留 | `store/operations_test.go TestMigrateExistingV3ToV4`（pre-v0005/pre-v0006 断言 + records 表消失）、`TestMigrate0005PreservesOperationLogRows`（两条 legacy + users.create + roles.create 重开持久化） |
| **重启持久化（双资源）** | users 与 roles 的写操作在**进程级重启**后存活；users、roles 均有 list/detail 身份断言，更新/创建响应的毫秒时间戳跨进程精确往返 | `cmd/server/server_restart_test.go TestServerProcessRestartPersistsUsers`（users create/patch/delete；roles create；重启后双资源 list/detail 与时间戳断言） |
| **401/403 双资源** | users 与 roles 各自 list/detail/create/update/delete 五路由：匿名均为 401 `UNAUTHENTICATED`；viewer 的 list/detail 为 200，create/update/delete 均为 403 `FORBIDDEN` | `handler/users_test.go TestUsersAuthGates`、`handler/roles_test.go TestRolesAuthGates` |
| **操作日志（双资源）** | users/roles 的 create/update/delete 事件记录准确 event、actor id/name、record id；create/update detail 仅含 username/key，delete detail 为 nil，不含密码 | `handler/users_test.go TestUsersOperationLogEvents`、`handler/roles_test.go TestRolesOperationLogEvents` |
| **领域保护（双资源）** | users self/last-admin/username 唯一；roles system/in-use/key 格式 | handler + store 测试套件（S2 A-003/A-004 已核） |
| **Schema-only 页面（双资源）** | users/roles 均由真实 fixture + manifest 进入通用 Renderer；roles create/update/delete 的 method/path/body 由 fixture action 驱动；两资源 action id 不在 Renderer/App 生产源码硬编码 | `renderer/schema-crud.test.tsx` T-UI-10（roles CRUD + 双资源 action id 反证）；`app/representative-pages.integration.test.tsx`（真实 manifest + roles fixture 页面） |

## 3. Renderer diff 边界（Schema-only 接入）

- **约束**：users/roles 页面接入只允许修改 **Schema/fixture（`fixtures/schema/*.json`）与 manifest（`app-manifest.json`）**；**不得**修改 Renderer 主路径（`src/renderer/render.tsx`、`schema-table.tsx`、`records.ts`、`form-controls.*`、`reactions.ts`、`permissions.ts` 与 `src/app/App.tsx` 的渲染逻辑）。
- **S3/S4 边界**：S3 已完成的 Renderer 泛化/专名清理（`recordsFetcher`→`resourceFetcher`、footer/toast 文案、`RecordsQuery`→`ResourceQuery`、删 `use-records.ts`）是 **records 退场动作**，不是 S4 的 Schema 接入。S4 基线固定为提交 **`adfe15a17da770699d5e109f22402c41ece5eeea`**。
- **受限生产文件**：`apps/web/src/renderer/{render.ts,render.tsx,schema-table.tsx,modal.tsx,confirm.tsx,form-controls.ts,form-controls.tsx,row-action.ts,records.ts,reactions.ts,permissions.ts}` 与 `apps/web/src/app/App.tsx`。测试文件不在受限列表内，但测试变化不能替代生产路径零 diff 证明。
- **可重复命令**：

```powershell
git diff --exit-code adfe15a17da770699d5e109f22402c41ece5eeea -- apps/web/src/renderer/render.ts apps/web/src/renderer/render.tsx apps/web/src/renderer/schema-table.tsx apps/web/src/renderer/modal.tsx apps/web/src/renderer/confirm.tsx apps/web/src/renderer/form-controls.ts apps/web/src/renderer/form-controls.tsx apps/web/src/renderer/row-action.ts apps/web/src/renderer/records.ts apps/web/src/renderer/reactions.ts apps/web/src/renderer/permissions.ts apps/web/src/app/App.tsx
```

- **通过判定**：命令 exit 0 且无 diff；T-UI-10 同时断言 users/roles action id 只由 fixture 声明，真实 manifest + fixture 页面测试覆盖 users 与 roles。

## 4. 信息门禁关闭结果（D-004）

- [x] D-004 已采纳本契约 v0.2.0，候选/冻结状态与 canonical 决策一致（A-007 F-001 fixed）。
- [x] §2 的 roles 页面级证据与后端双资源完整断言已补齐（A-007 F-002/F-003 fixed）。
- [x] §3 已固定 baseline revision、受限生产文件与可执行 diff 命令；2026-08-03 冻结预检 exit 0。
- [x] `I-011-003` → `verified`，S4 信息门禁解除。

## 5. 完成定义（S4 勾选口径）

- [ ] 以本契约 v0.2.0 形成单独 S4 验收收据，记录 revision、完整命令、时间与结果；冻结预检不自动替代阶段验收。
- [ ] §2 全部验收维度在 S4 验收 revision 上通过；§3 Renderer/App 基线命令 exit 0。
- [ ] S4 阶段审视已汇总，且无开放 required finding。
- [ ] S4 检查点达成后，才将 progress `3/5 → 4/5` 并同步 `goal-tree.md`。
- [ ] GOAL-010 S4 父级验收门证据链就绪（S5 交接）。

## 6. 非目标

- S5 全量回归与关门审计；Root A-002 F-002-001 关闭证据；扩大 `I-PROTO-001` 覆盖；新增业务资源（非 users/roles）。

## 7. 证据索引

- `apps/api/internal/store/migrate_test.go`、`operations_test.go`、`seed_test.go`、`restart_test.go`
- `apps/api/internal/handler/{users_test.go, roles_test.go, account_test.go, operations_test.go}`
- `apps/api/cmd/server/server_restart_test.go`
- `apps/web/src/renderer/schema-crud.test.tsx`（T-UI-10）、`apps/web/src/app/representative-pages.integration.test.tsx`、`apps/web/src/renderer/representative-pages.test.tsx`
- `apps/api/internal/handler/fixtures/schema/{users.json, roles.json}`、`apps/web/public/.well-known/schema-ui/app-manifest.json`
- 冻结预检：`go test ./...`；Web `npm test`（23 files / 485 tests）；`npm run build`；§3 Renderer/App diff 命令（均于 2026-08-03 通过）

## 8. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 0.1.0 | 2026-08-03 | 候选；A-007 审计发现状态真实性、双资源 Schema-only 与后端断言三项 required 缺口，未冻结 |
| 0.2.0 | 2026-08-03 | D-004 冻结；固定 Renderer baseline/命令，补齐 roles 页面级与后端双资源断言，A-007 F-001～F-003 fixed，`I-011-003` verified |
