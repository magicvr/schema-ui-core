---
title: 执行记录 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-010-a002-schema-adapter
version: 0.3.0
---

# 执行记录 · GOAL-011

## 2026-08-03 · 目标立项

- 用户确认新建 `GOAL-011-s4-semantic-admin-resources`，选择 `users` 替换 records 默认代表实体、`roles` 作为第二个语义资源（D-001）。
- 从 canonical 目标模板建立五件套与 `attachments/`，设定 `parent: GOAL-010-a002-schema-adapter`，并写入五个等权顺序检查点。
- 登记三个 required 信息项：`I-011-001`（users/roles 领域与安全契约）、`I-011-002`（records 退场/迁移）、`I-011-003`（双资源集成验收）；初始均为 `open`。
- 同步修订 GOAL-010 S4 为本目标交付后的父级验收门，并更新当前工作区 `goal-tree.md` 的树与状态表。
- 未修改 API/Web/迁移/fixture 等产品代码；未移除 records；未实现 users/roles CRUD；未关闭 Root A-002 F-002-001。
- **计划（非事实）**：先收集 `I-011-001`/`I-011-002` 并形成版本化契约候选，经用户裁决冻结 S1 后再进入实现。

## 2026-08-03 · S1 契约冻结（D-002）

- 完成 `I-011-001`/`I-011-002` 信息收集：静态核对 auth/RBAC 表结构（users/roles/user_roles/permissions/menu/operation_log）、通用工厂 `resources.go`、身份投影 `account.User`、`userWithRoles` 双写、records 足迹（API/迁移/种子/权限/菜单/操作日志/fixture/前端/测试/数据文件）。
- **用户裁决（P-004）三项关键取舍（均采纳推荐）**：
  1. 通用工厂 + 最小契约扩展（`Resource.JSONFields` + `DomainError` 409 映射）——users/roles 均走工厂五路由；
  2. 操作日志纳入（migration `0005` 扩展 operation_log event CHECK，新增 users/roles 事件，保留 records/auth 历史值）；
  3. records 硬退场 DROP TABLE（migration `0006` 删表 + 清理权限/菜单行，既有库自动 `pre-v0006` 快照兜底）。
- 落盘两份冻结契约：`attachments/I-011-001-users-roles-contract.md` v0.1.0、`attachments/I-011-002-records-retirement.md` v0.1.0。
- `I-011-001`/`I-011-002` → **verified**（契约 + D-002）；`I-011-003` 保持 open（最晚 S4）。
- **S1 检查点达成，progress `0/5 → 1/5`**；未修改任何产品代码。
- 同步修订 `00-meta`（S1 勾选、信息表状态、版本 v0.2.0）、`01-decision`（D-002）、当前工作区 `goal-tree.md`（GOAL-011 `1/5`）。
- **计划（非事实）**：S2 后端 users/roles 资源闭环 → S3 records 退场 → S4 双资源 Schema 接入 → S5 回归审计关门；关键节点自审后调用 grok build 独立审计。

## 2026-08-03 · S1 交叉审计 + 合并响应（A-001 → A-002 → D-003）

- **S1 自审（A-001 · self · pass）**：契约完整、可实施、与既有事实一致；无 required；F-001～F-003 recommended（password 长度、fixture 文案、DomainError 优先级）。
- **grok build 独立审计（A-002 · independent · conditional）**：S1 主体成立；**F-001（required/med）** 工厂扩展未规定 actor 通道，SELF_OPERATION 不可在通用五路由内诚实实现；**F-002（required/med）** `migrate.go` 仅 first-pending 快照一次，0005+0006 同批时 `pre-v0006` 验收字面失败；F-003～F-006 recommended（roles 响应形状、`linkUserRole` 隐式建角色、父契约 409 双真相、承接 A-001 low）。
- **用户裁决（P-004 §3.2）**：A-001 pass 与 A-002 conditional 同 scope 冲突 → 裁决「**全部 fixed**」。
- **响应落盘（D-003）**：I-011-001 → **v0.2.0**（§7 actor 通道 + DomainError 优先级、§2.3 禁 ensureRole 隐式建角色、§3.0 roles 响应形状）；I-011-002 → **v0.2.0**（§2.3 每待应用数据变更迁移前快照、§5 验收对齐）；GOAL-010 **D-005** + I-010-001 **v0.2.2**（§5 账号域 409 注记）。
- `I-011-001`/`I-011-002` 维持 `verified`（v0.2.0 响应修订）；`I-011-003` 保持 open。S1 无开放 required，A-001/A-002 趋同；S2 实施门禁保持解除。
- 同步修订 `00-meta` v0.2.1、`01-decision` D-003、`03-audit` 响应节、`goal-tree.md` 注记。
- **计划（非事实）**：S2 按 I-011-001 v0.2.0 实施通用工厂扩展 + users/roles 后端闭环。

## 2026-08-03 · S2 后端 users/roles 资源闭环（I-011-001 v0.2.0）

- **通用工厂扩展**（`resources.go`）：`ResourceEntity.Create/Update/Delete` 增传 `account.User`（A-002 F-001 actor 通道）；`Resource.JSONFields`（users.roles 原始 JSON 透传，decode create/patch 支持）；`DomainError{Status,Code,Message}` + `writeEntityError` 统一映射（先 DomainError → 再 ErrNotFound → 最后 INTERNAL，A-002 F-006/A-001 F-003）；records 实体签名补齐并忽略 actor（零对外行为变化）。
- **store 领域方法**：`store/users.go`（ListUsers/GetUser/CreateUserManagement/UpdateUser/DeleteUser + `reconcileRoles` 双写集合一致性 + self/last-admin 保护 + refresh_tokens 级联清理 + **不隐式建角色** A-002 F-004）；`store/roles.go`（ListRoles/GetRole/CreateRole/UpdateRole/DeleteRole + system/in-use/invalid-key 保护）。
- **migration 0005** `operation_log_expand`：重建 operation_log event CHECK，新增 `users.*`/`roles.*` 事件，保留 `records.*`/`auth.*` 历史值；SQLite 表内重建 + 行迁移 + 索引；`transformID 0005:operation-log-expand:v1` 冻结 checksum。
- **种子**（`seed.go`）：新增 `users.read/write`、`roles.read/write` 权限、`menu_users`/`menu_roles` 菜单；admin 四权限 + 两菜单，editor/viewer 只读（users.read/roles.read）；records 种子保持（S3 退场）。
- **注册**：`health.go` 注册 `/api/users`、`/api/roles`；`account.StaticDevSession` 同步 users/roles 权限与菜单。
- **测试**：新增 `handler/users_test.go`（list/detail/CRUD、password_hash 隔离、重复 username 409、self 保护、401/403、password 仅写可登录、users.* 操作日志）、`handler/roles_test.go`（list/detail、创建/更新/删除、invalid-key/duplicate/system/in-use 保护、401/403、roles.* 操作日志）、`store/users_test.go`（last-admin、role 校验、无隐式角色、双写往返）、`store/roles_test.go`（create 校验、system/in-use 保护）；既有迁移账本/种子测试更新至 0005 + users/roles 种子。
- **证据**：`go test ./...` 全绿（126 顶层测试函数；`-v` RUN 计 151 含子测试）+ `go vet ./...` 干净；web `vitest run` 481/481 + `tsc -b` 干净（后端变更未破坏前端）。
- 修复实现中发现的问题：单连接嵌套查询死锁（ListUsers 先收行再 reconcile）、`updated_at` INTEGER 扫描 time.Time 类型错误（UpdateUser）。
- **S2 检查点达成，progress `1/5 → 2/5`**；records 仍注册（S3 退场）；Root A-002 F-002-001 仍 open。
- **计划（非事实）**：S3 records 退场（0006 DROP TABLE + 清理权限/菜单 + 每待应用版本前快照 F-002 落地 + API/种子/fixture/前端/测试退场）。

## 2026-08-03 · S2 交叉审计 + 合并响应（A-003 → A-004 → 响应节）

- **S2 自审（A-003 · self · pass）**：工厂扩展、领域不变量、错误码、401/403、0005 迁移、回归证据均对齐契约；无 required；F-001～F-003 recommended（per-pending 快照随 S3、重启/搜索/排序随 S4、password 长度非目标）。
- **grok build 独立审计（A-004 · independent · pass）**：与 A-003 同向，无 required；F-001（0005 行保留缺专用升级回归）、F-002（LAST_ADMIN 缺 HTTP 层断言）、F-003（StaticDevSession 投影缺回归断言）→ **fixed**（新增 `TestMigrate0005PreservesOperationLogRows`、`TestUsersLastAdminHTTP`、`TestAccountsMeDevSessionFallback` 增补）；F-004（承接 A-003 + 测试计数口径）→ **handled**（计数口径修正为 126 顶层 / 151 RUN）。
- 未触发 P-004 裁决（A-003/A-004 verdict 一致 pass，无冲突）。S2 无开放 required，可放行 S3。
- 同步修订 `03-audit` 响应节与索引、`goal-tree.md` 注记。
- **计划（非事实）**：S3 records 产品运行面退场。

## 2026-08-03 · S3 records 产品运行面退场（I-011-002 v0.2.0）

- **migration 0006 `records_retire`**：`DROP TABLE records` + 清理 records 权限/菜单行（`perm-records-*`、`menu-list-edit-lifecycle` 及 grants，先删 join 再删父行 FK 安全）；`transformID 0006:records-retire:v1` 冻结 checksum。**per-pending 快照**（A-002 F-002 / A-003 F-001）：`migrate()` 对每个待应用数据变更迁移（≥2）前各快照一次，0005+0006 同批时 `pre-v0005` 与 `pre-v0006` 均存在（`TestMigrateExistingV3ToV4` 断言 pre-v0006）。
- **后端退场**：`health.go` 移除 records 注册（改 users/roles）；删除 `handler/records.go`、`store/records.go`、`store/seed_records.go`；`seed.go` 移除 records 权限/菜单/grants（users/roles 保留）；`store.go` 移除 seedRecords 调用；`account.StaticDevSession` 移除 records 键；`jsonQuote`/`newOperationID`/`ErrRecordExists` 迁至共享位置。
- **fixture/manifest**：新增 `users.json`/`roles.json`（两语义资源 CRUD 页，I-011-002 §3.3）；删除 `list-edit-lifecycle.json`/`catalog.json`；`data-table.json`/`search-form-table.json` 改指 `/api/users`/`/api/roles` 并更新文案；manifest 移除 records/catalog 页 + `menu_list_edit_lifecycle`，新增 users/roles 页（`menu_users`/`menu_roles` visibleWhen）。
- **前端 records 专名退场**：`records.ts` 移除 `RecordItem`/`RecordList` 别名、`RecordsQuery`→`ResourceQuery`（`fetchRecords` 通用 transport 保留）；删除 `use-records.ts`；`App.tsx`/`main.tsx` `recordsFetcher`→`resourceFetcher`；`schema-table.tsx` footer "records"→"items"、empty/caption 文案去 records；`render.tsx` 成功 toast "Record *"→"Item *"。
- **测试退场/改指**：删 records handler/store/restart 测试；`cmd/server` 进程级重启测试改指 users（create/patch/delete + 重启持久化）；`resources_test`/`operations_test`/`seed_test`/`migrate_test`/`restart_test` 更新（6 迁移账本、无 records 表、users/roles 种子计数、`menu_users` 投影）；web `schema-crud.test.tsx`/`representative-pages*.test`/`navigation.test`/`app-manifest.test`/`upstream-fixtures.test`/e2e `schema-crud.spec`/`shell.spec` 改指 users/roles；QUICKSTART/smoke.sh/playwright.config 更新。
- **验收口径达成**：fresh install 无 `records` 表（`TestMigrateFreshDB`）；升级路径 pre-v0006 快照存在；产品代码 grep 无 `api/records`/`records.read`/`records.write`/`menu_list_edit_lifecycle`/`list-edit-lifecycle` 残留（仅 0006 清理语句与历史注释）；users/roles 操作日志事件生效。
- **证据**：`go test ./...` 全绿 + `go vet ./...` 干净；web `vitest run` 481/481 + `tsc -b` + `vite build` 干净；e2e `playwright test` 2/2 通过（users CRUD 真实 Go/SQLite 往返 + shell/auth 链）。
- **S3 检查点达成，progress `2/5 → 3/5`**；records 已从产品默认运行面退场；Root A-002 F-002-001 仍 open。
- **计划（非事实）**：S4 双语义实体 Schema 接入验证（I-011-003 冻结 + fresh/upgrade/restart/401-403 矩阵 + Renderer diff 证据）→ S5 回归审计关门。

## 2026-08-03 · A-007 必改修复与 I-011-003 冻结（D-004）

- **用户裁决**：不补同范围自审，A-007 F-001～F-003 全部直接 `fixed`；未选择 residual 或 overruled。
- **F-001 · 契约真实性**：将 A-007 审计时的候选矩阵 v0.1.0 修订为 v0.2.0，并在 D-004 实际落盘后冻结；`00-meta.md` 的 `I-011-003` 同步由 `open` 改为 `verified`。
- **F-002 · 双资源 Schema-only**：
  - `apps/web/src/renderer/schema-crud.test.tsx` 的 T-UI-10 增加 roles 真实 fixture 驱动的 create/update/delete 请求与 body/path 断言，并将 action-id 无硬编码检查扩展到 users + roles。
  - `apps/web/src/app/representative-pages.integration.test.tsx` 增加真实 manifest + roles fixture 的页面级渲染断言（toolbar、row actions、rows、recordView）。
  - Renderer/App 生产路径基线冻结为 `adfe15a17da770699d5e109f22402c41ece5eeea`；按 I-011-003 §3 的精确命令检查为零 diff。
- **F-003 · 后端完整边界**：
  - `TestUsersAuthGates` / `TestRolesAuthGates` 覆盖各自 list/detail/create/update/delete：匿名均为 401 `UNAUTHENTICATED`，viewer 读为 200、写为 403 `FORBIDDEN`。
  - `TestUsersOperationLogEvents` / `TestRolesOperationLogEvents` 增加 actor id/name、record id、create/update 非敏感 detail 与 delete nil detail 断言。
  - `TestServerProcessRestartPersistsUsers` 增加 roles create 响应时间戳格式、重启后 list/detail 身份及毫秒时间戳精确往返。
  - `TestMigrate0005PreservesOperationLogRows` 增加升级表上的 `roles.create`，关闭/重开后同时核对 users/roles 新事件及两条 legacy 事件。
- **验证事实**：
  - 定向 Web：`npm test -- --run src/renderer/schema-crud.test.tsx src/app/representative-pages.integration.test.tsx` → 2 files / **26 tests passed**。
  - 定向 API：`go test ./internal/handler ./internal/store` → 两包通过。
  - 完整 API：`go test ./...` → 全包通过，含 `cmd/server` 进程级重启测试。
  - 完整 Web：`npm test` → **23 files / 485 tests passed**；`npm run build` → `tsc -b` + Vite production build 通过。
  - Renderer/App 基线：I-011-003 §3 的 `git diff --exit-code adfe15a... -- <受限生产文件>` → exit 0、无输出。
- **治理投影**：A-007 F-001～F-003 均 `fixed`；I-011-003 v0.2.0 + D-004 → `verified`，只解除信息门禁。GOAL-011 保持 `active / 3/5`，S4/S5 未勾选，goal-tree 状态与进度不变。
- **计划（非事实）**：另行执行 S4 冻结矩阵，形成包含完整命令、revision 与结果的阶段验收收据；经阶段审视后再决定是否把 progress 推进至 `4/5`。

## 2026-08-03 · S4 阶段验收收据（I-011-003 v0.2.0 §5）

- **验收 revision**：`73bc93abba52db0440bc8c70eaf89969174a00cc`（GOAL-011 归档提交）；Renderer/App 基线 `adfe15a17da770699d5e109f22402c41ece5eeea`（I-011-003 §3）。
- **命令与结果**：
  | 命令 | 结果 |
  |------|------|
  | `git diff --exit-code adfe15a… -- apps/web/src/renderer/{render.ts,render.tsx,schema-table.tsx,modal.tsx,confirm.tsx,form-controls.ts,form-controls.tsx,row-action.ts,records.ts,reactions.ts,permissions.ts} apps/web/src/app/App.tsx` | **exit 0**，无 diff（Renderer/App 主路径自基线零修改） |
  | `go vet ./...`（apps/api） | 干净 |
  | `go test ./... -count=1`（apps/api） | 全包通过（cmd/server / account / auth / config / handler / store） |
  | `npx vitest run`（apps/web） | **23 files / 485 tests passed** |
  | `npx tsc -b`（apps/web） | 干净 |
  | `npx vite build`（apps/web） | 生产构建成功 |
  | `WEB_PORT=9999 npx playwright test`（apps/web/e2e） | **2/2 通过**（users CRUD 真实 Go/SQLite 往返 + shell/auth 链） |
- **§2 验收维度对照**（全部通过，证据指针见 I-011-003 §2）：
  - fresh fork：`TestMigrateFreshDB`（无 records 表 + users.create 事件可写）✅
  - 既有库升级：`TestMigrateExistingV3ToV4`（pre-v0005/pre-v0006 + records 表消失）、`TestMigrate0005PreservesOperationLogRows`（legacy 行保留 + users/roles 新事件重开持久化）✅
  - 重启持久化：`TestServerProcessRestartPersistsUsers`（users create/patch/delete + roles create；重启后双资源 list/detail 与毫秒时间戳往返）✅
  - 401/403 双资源：`TestUsersAuthGates`/`TestRolesAuthGates`（五路由匿名 401、viewer 读 200/写 403）✅
  - 操作日志：`TestUsersOperationLogEvents`/`TestRolesOperationLogEvents`（actor id/name、record id、非敏感 detail、delete nil detail）✅
  - Schema-only 页面：T-UI-10（users/roles action id 无 Renderer 硬编码）+ `representative-pages.integration.test.tsx`（真实 manifest + users/roles fixture 页面渲染）✅
- **§3 Renderer diff 边界**：S4 基线命令 exit 0（上表首行）；S4 未修改任何受限生产文件。
- **S4 阶段审视**：I-011-003 全部验收维度在验收 revision 通过；A-001～A-008 已响应，本 scope 无开放 required；S4 检查点达成。
- **治理投影**：**S4 勾选，GOAL-011 `3/5 → 4/5`**；同步修订 `00-meta`（S4 勾选、version 0.6.0）与 `goal-tree.md`（4/5 + 注记）；Root A-002 F-002-001 仍 open（S5 交接后闭合）。
- **grok build S4 独立审计调用记录（未完成）**：按用户指令，A-009 自审后调用 grok build 独立审计（scope: S4 双资源 Schema 接入验证）。**grok 服务端 5 次会话均被取消**（`stopReason: cancelled`，最短 2 轮即中断，未产出书面 A-010 意见）。S4 阶段审视暂以 A-009（self · pass）为主，叠加既有独立意见 A-007（conditional → fixed）/A-008（independent · pass，I-011-003 finding-closure）；grok 对 S4 的独立审计**留待 S5 关门时重试**，如仍不可用则作为用户可复核的待办记录（不以未落盘的独立意见作为放行依据）。
- **计划（非事实）**：S5 回归、审计与父级交接（全量回归 + 关门审计 + grok build 重试 + GOAL-010 S4 证据交接）。

## 2026-08-03 · S5 回归、审计与父级交接

- **全量回归（S5 事实，独立于 S4 收据重跑）**：
  | 命令 | 结果 |
  |------|------|
  | `git diff --check`（仓库根） | 干净（仅 LF→CRLF 提示） |
  | `go vet ./...`（apps/api） | 干净 |
  | `go test ./... -count=1`（apps/api，180s） | 全包通过（cmd/server / account / auth / config / handler / store） |
  | `npx vitest run`（apps/web） | **23 files / 485 tests passed** |
  | `npx tsc -b`（apps/web） | 干净 |
  | `npx vite build`（apps/web） | 生产构建成功 |
  | `WEB_PORT=9999 npx playwright test`（apps/web/e2e） | **2/2 通过** |
  | I-011-003 §3 Renderer 基线 `git diff --exit-code adfe15a… -- <受限生产文件>` | **exit 0** |
- **关门审计（self）**：A-010（close-out · self · pass）——五个检查点全部达成、无开放 required finding、无到期 required 信息项、GOAL-010 交接就绪；F-001（grok S4/关门取消）为可复核待办。
- **grok build 关门独立审计**：重试成功，产出 **A-011**（independent · conditional，文本态意见已代贴入 03-audit）——F-001（S5 执行事实未入 02-execution）、F-002（关门文档真相不一致）required → 本 S5 事实节 + meta/审计边界对齐后 `fixed`；F-003～F-005 recommended 随关闭或 handled。
- **父级交接**：GOAL-010 `02-execution` 已落「GOAL-011 S4 证据交接」节（users/roles 替换、Renderer 零修改、records 退场、双资源验收收据）。
- **治理投影**：S5 关门审计通过（A-010 + A-011 经 F-001/F-002 修复后趋同），**GOAL-011 置 `done`，progress `4/5 → 5/5`**；同步 `00-meta`（S5 勾选、status done、version 0.7.0）与 `goal-tree.md`（5/5 + done + 注记）；GOAL-010 据此可评估其 S4 勾选与 S5 关门；Root A-002 F-002-001 仍 open（GOAL-010 S5 关闭证据链）。

## 2026-08-04 · A-012 编排响应与关门门禁恢复（D-005）

- 核对 active Charter `schema-ui-core-admin-foundation@0.1.0`、VP-002 `active`、workspace-002 `delivery` / `primary_plan` / Root / canonical 绑定；`shared_materials_catalog: none`，本轮未使用共享资料，愿景与工作区链不构成响应阻断。
- 汇总同一 close-out scope：A-010（self · pass）、A-011（independent · conditional，经 required fixed 后趋同）与 A-012（independent · fail）。A-012 新增 F-001～F-005 五条 required，当前均无合法闭合记录；P-004 冲突保持未决。
- 本轮静态抽查确认 A-012 的主要代码/流水线落点仍存在：
  - F-001：users 写门仅 `users.write`，create/patch 可携带任意已存在 role key，未校验 actor 的可委派权限；
  - F-002：users fixture 密码字段仍为普通 `input`；通用字符串解码 trim 密码；改密事务未撤销 refresh token；
  - F-003：users Schema 未提供角色分配/改密，roles Schema 只管理 key/name；自定义 role 无 grant 管理路径，system 行仍渲染必失败动作；
  - F-004：活动 CI 仍调用 `/api/records` 并传 `SMOKE_RECORD_ID`/seed total 8；`smoke.sh` 注释默认 1、实现默认 8；通用 Web transport 与生产 imports 仍使用 records 专名；
  - F-005：`DeleteUser` 在事务外读取用户/角色，随后才开启最后管理员检查与删除事务。
- 新增 D-005 与 `I-011-004`；A-012 F-001～F-005 保持 `open`，F-006 保持 recommended。未修改 API/Web/CI/smoke 产品代码，未运行新的回归或复现命令，也未把 A-012 已有测试收据当作本轮修复证据。
- **治理投影**：GOAL-011 `done → active`；原 S1～S5 检查点仍全部完成，派生 `progress` 保持 `5/5`。同步 `00-meta`、`03-audit` 响应节与 `goal-tree.md`；GOAL-010、Root、VP-002 状态/进度不变。
- **计划（非事实）**：用户先裁决五项 finding 的闭合路径及 `I-011-004` 产品边界；推荐 F-001/F-002/F-004/F-005 走 fixed，F-003 选择明确边界后 fixed。实施完成后补 API/Web/CI/E2E 收据并请求限定范围独立复审。
