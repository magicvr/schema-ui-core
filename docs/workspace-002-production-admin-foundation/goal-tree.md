---
title: workspace-002-production-admin-foundation · 目标树
status: active
created: 2026-08-01
updated: 2026-08-02
parent: null
version: 0.27.0
workspace_id: workspace-002-production-admin-foundation
---

# 目标树

- 工作区：`workspace-002-production-admin-foundation`
- canonical root：`docs/workspace-002-production-admin-foundation/`
- Root Goal：`GOAL-001-production-admin-foundation`
- primary plan：`VP-002-production-admin-foundation`

```text
GOAL-001-production-admin-foundation [active] 生产级可用 Admin 基架 (4/5)
├── GOAL-002-r1-schema-load-validate [done] R1 · Schema 加载、校验与统一错误面 (4/4)
├── GOAL-003-r1-default-render-path [done] R1 · 默认 Renderer 主路径与示例降级 (4/4)
├── GOAL-004-r1-representative-node-pages [done] R1 · 代表性 Node 页面与回归证据 (5/5)
├── GOAL-005-r2-auth-session [done] R2 · 真实认证与请求级身份 (6/6)
├── GOAL-006-r3-persistent-rbac-menu [done] R3 · 持久化 RBAC、菜单投影与版本迁移 (6/6)
├── GOAL-007-r4-schema-crud [done] R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环 (6/6)
└── GOAL-008-r5-engineering-fork [active] R5 · 工程化、fork 体验与集成关门 (2/5)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| `GOAL-001-production-admin-foundation` | 生产级可用 Admin 基架 | `null` | `active` | `4/5` | 2026-08-02 |
| `GOAL-002-r1-schema-load-validate` | R1 · Schema 加载、校验与统一错误面 | `GOAL-001-production-admin-foundation` | `done` | `4/4` | 2026-08-01 |
| `GOAL-003-r1-default-render-path` | R1 · 默认 Renderer 主路径与示例降级 | `GOAL-001-production-admin-foundation` | `done` | `4/4` | 2026-08-02 |
| `GOAL-004-r1-representative-node-pages` | R1 · 代表性 Node 页面与回归证据 | `GOAL-001-production-admin-foundation` | `done` | `5/5` | 2026-08-02 |
| `GOAL-005-r2-auth-session` | R2 · 真实认证与请求级身份 | `GOAL-001-production-admin-foundation` | `done` | `6/6` | 2026-08-02 |
| `GOAL-006-r3-persistent-rbac-menu` | R3 · 持久化 RBAC、菜单投影与版本迁移 | `GOAL-001-production-admin-foundation` | `done` | `6/6` | 2026-08-02 |
| `GOAL-007-r4-schema-crud` | R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环 | `GOAL-001-production-admin-foundation` | `done` | `6/6` | 2026-08-02 |
| `GOAL-008-r5-engineering-fork` | R5 · 工程化、fork 体验与集成关门 | `GOAL-001-production-admin-foundation` | `active` | `2/5` | 2026-08-02 |

> Root `4/5` 由五个等权纲领检查点派生（R1、R2、R3、R4 已勾选；R5 待办）。子目标 progress 仅反映各自成功标准，不替代 Root 检查点。依赖：003 硬依赖 002；004 完整主路径证明依赖 002+003。`GOAL-002` 已于 2026-08-01 `done`（A-001 independent + A-002 self 关门审计；无开放 required）。`GOAL-003` 已置 `done`（2026-08-02：A-001/A-002/A-003/A-004 全 pass，无开放 required；`I-003-001/002` closed；F-001 recommended open → R4 follow-up）。`GOAL-004` 已置 `done`（2026-08-02：A-001 self + A-002 independent 关门审计全 pass，无开放 required；`I-004-001/002` closed；F-001 → fixed，F-002 recommended open → R4 follow-up）。R1 三个子目标（002/003/004）全部 `done`，Root R1 检查点已勾选（I-001 覆盖矩阵 verified + Renderer 默认主路径 + 425/425 回归与 fail-closed 证据）。`GOAL-005` 已置 `done`（2026-08-02：A-001 independent close-out `conditional`，F-001 → **fixed**（Linux CI run #30711903555 browser E2E `1 passed` + 匿名 401 断言 + 403 由 records_test 承担）；A-002 self close-out `pass`；`I-005-001/002/003/004/005` verified，无开放 required；D-002～D-007）；**Root R2 检查点已勾选**（Root progress 1/5 → 2/5）。Root D-009 已关闭 `I-003`（`verified`）并冻结方案 B、`features` 菜单投影、两步迁移、读写权限与恢复证据口径；`GOAL-006` 已立项为 `active / 6/6`，其 `I-006-001/002` 已由 D-002/D-003 与附件验证关闭，当前无开放 required 信息门禁。S1 已实现（D-004：migration runner + `schema_migrations` + `0001/0002` 迁移链 + pre-v0002 恢复快照）；S2 已实现（D-005：阶段 B 终态规范化权威读 + 双写 + 集合核对，A-002 F-004 required → fixed、A-003 独立复核 pass）；S3 已实现（D-006：`seedRBAC` 增量幂等种子接线到 Open，A-004 self 阶段审计 pass）；S4 已实现（D-007：records 读写经身份携带的 permission key 门禁，匿名 401 / 缺权限 403）；S5 已实现（D-008：`/me.features` 从持久化菜单 grants 投影 + 真实 manifest `visibleWhen`）；**S6 已实现（2026-08-02：恢复/重启/回归证据——`TestRestartPersistence`、`TestRestorePreV0002Snapshot`、E2E 重启冒烟、全仓 API/Web 回归；A-001 F-001/F-002 与 A-004 F-101 已闭合为 fixed）**。**A-005 F-005 已闭合为 fixed（2026-08-02：API README 端点表与鉴权边界同步 R3 权限键语义，API 回归通过；意见台账开放 0）**。**`GOAL-006` 已置 `done`（2026-08-02：A-006 independent 复核 F-005 关闭 `pass`；A-001～A-006 全部 responded，无开放 required；`I-006-001/002` verified；close-out 用户裁决不补 self 自审）。Root R3 检查点已勾选（Root `2/5 → 3/5`）**。六项检查点已全勾选（`6/6`）。

> Root D-010 已关闭 `I-004` 的方向/立项目门禁；D-011 已创建 `GOAL-007-r4-schema-crud`。2026-08-02：GOAL-007 D-002/D-003 已关闭 `I-007-001`/`I-007-002`（verified），勾选 S1/S2，进度 `0/6 → 2/6`；证据见目标附件 API 契约与 SQLite 迁移计划。**A-001 F-001 已闭合为 fixed（2026-08-02：D-004 统一 `updatedAt` 毫秒精度与「严格晚于」断言，附件更新 v0.2.0，A-002 self 复核 pass），S3 持久化 API 实施已放行但尚未实施**；`I-007-003`/`I-007-004` 仍为 open required（阻断 Schema 写交互与 S6 验收）；Root R4 仍未勾选，Root 保持 `3/5`。**S3 已实施（2026-08-02：0003 records_persist（updated_at Unix 毫秒）+ 快照通用化 pre-v0003 + records repository + seedRecords 空表种子 + handler 重写走 SQLite，新增 POST 201/INVALID_CREATE_*；staticRecords 与进程切片移除；T-API-08～13 与 T-DB-01～09（含毫秒钳制/往返）全绿，API `go test ./...` 与 web vitest 443/443 通过；API README 同步 R4）。成功标准 S3 勾选，进度 `2/6 → 3/6`**。

> **A-003 已响应、I-007-003 已冻结（2026-08-02）**：A-003（independent · pass）响应节已写入——A-002 R-001（毫秒钳制/往返）→ 已落实；A-003 R-001（Root 台账同步）→ fixed（Root 00-meta 纲领 R4 已同步 `3/6`、S3 已实施，Root 保持 `3/5`）；A-003 R-002（PATCH trim 一致性）→ fixed（`UpdateRecord` 入库 trim + store/handler 回归）。**`I-007-003` → verified（D-005 冻结 `list-edit-lifecycle` 代表页 + table actions/toolbar + form `submitAction` + 字段/交互/权限矩阵 + T-UI-01～10），首个 Schema 写交互代码变更已放行**。`I-007-004` 仍 open required（S6）；Root R4 未勾选。S3 阶段 self 审计按用户裁决留待 S4/S5 放行或关门前。

> **A-004/A-005 已响应、I-007-003 升级 v0.2.0（2026-08-02）**：A-004（independent · pass）复核 S3 + A-003 R-001/R-002 关闭证据成立；A-005（independent · conditional）三处 required 均 `fixed`——I-007-003 **v0.2.0** + D-005 补记闭合 F-002（唯一结构：table + modal create-form + modal edit-form + 行 delete）、F-003（`records.write` → `permissions.edit/delete` 表达式 + cascade，禁止 intent 单独）、F-004（§9 最小 actions/capabilities/`$row`/search 归属），R-001/R-002 handled。**`I-007-003` 保持 verified（v0.2.0），S4/S5 实施放行维持，可无歧义开工**。当前 scope 无开放 required；`I-007-004` 仍 open（S6）；Root R4 未勾选。后续意见从 A-006 起。

> **A-006 已响应、I-007-003 升级 v0.2.1（2026-08-02）**：A-006（independent · conditional）复核 v0.2.0 修订——F-002/F-003 维持 fixed；F-004 残余由 F-005/F-006 承接并随本轮闭合——I-007-003 **v0.2.1** + D-005 v0.2.1 补记闭合 **F-005**（§9.1a 冻结 form submit 的 `{id}` 槽绑定：从打开 modal 时捕获的行上下文解析，`formAction` 有界扩展落入 §9.5 白名单并补测试）与 **F-006**（§9 对齐 `action.schema`/registry：顶层 5 action（3×RequestAction + 2×ModalAction）、`onSuccess.behavior`、挂载字段 `actionRef`、`confirm` 移到 rowAction、delete `requestMapping.path.id` 留 rowAction）；R-001（§9.5 允许新增 modal/confirm 文件）与 R-002（D-005 旧表述取代）handled。**`I-007-003` 保持 verified（v0.2.1），S4 fixture/actions 接线可无歧义开工**。当前 scope 无开放 required；`I-007-004` 仍 open（S6）；Root R4 未勾选。后续意见从 A-007 起。

> **A-007 已响应、I-007-003 升级 v0.2.2（2026-08-02）**：A-007（independent · conditional）复核 v0.2.1——F-005/F-006 主体维持 fixed；F-006 残余 **F-007**（§9.2 `confirm` 写成对象、registry 要求 **string**）随本轮闭合——I-007-003 **v0.2.2** + D-005 v0.2.2 补记将 §9.2 改为 **`confirm: "Delete this record?"`**（string，一行修补）；R-001（§9.5 补入 `request-construction.ts`/`row-action.ts`）与 R-002（§9.1 注明 `reload` 隐含关 modal）handled。**`I-007-003` 保持 verified（v0.2.2），§9 字面形状已对齐 `action.schema`/registry，S4 代表页 fixture 可按 v0.2.2 编写**。当前 scope 无开放 required；`I-007-004` 仍 open（S6）；Root R4 未勾选。后续意见从 A-008 起。

> **S4/S5 已实施（2026-08-02）**：`list-edit-lifecycle` 演进为 v0.2.2 结构（table + toolbar/rowActions + 双 modal form + 行 delete 确认 + recordView 选中行 + 顶层 5 actions）；`search-form-table` 纳入 search form-to-query；渲染层一次性补齐（`render.tsx` SchemaCrudProvider、`schema-table.tsx` toolbar/actions/选中、`data-table.tsx` 行选中、新 `modal.tsx`/`confirm.tsx`、`records.ts` `createRecord` + envelope）；D-006 固定实现落点（`{id}` 槽绑定落 `render.tsx`、T-UI-10 判据、fetcher 首值保持、只读=禁用）。T-UI-01～10 全绿；web `vitest run` **458/458**、`tsc -b` 干净、`vite build` 成功、`go test ./...`（apps/api）全绿。**S4/S5 勾选，GOAL-007 进度 `3/6 → 5/6`**；`I-007-003` 保持 verified（v0.2.2）；`I-007-004` 仍 open（S6）；Root R4 未勾选（Root 保持 `3/5`）。无自审（沿用用户裁决，留待放行或关门前）。后续意见从 A-008 起。

> **A-008/A-009 已响应、S6 已实施（2026-08-02）**：A-008（independent · pass）S4/S5 execution-facts 复核；用户按 P-004 §3.1 裁决「先补 self 审计」→ **A-009（self · pass）** 补齐 S4/S5 scope 的 self 覆盖并统一响应 A-008。**`I-007-004` → verified**（D-007 + [I-007-004-restart-e2e-protocol.md](GOAL-007-r4-schema-crud/attachments/I-007-004-restart-e2e-protocol.md)：L1 HTTP 层 + **L2 进程级**重启协议）。**S6 已实施**——L1 `handler/records_restart_test.go` `TestRecordsSurviveRestart`（同文件 store 关闭→重开，全 HTTP CRUD→list/detail）+ **L2 `cmd/server/server_restart_test.go` `TestServerProcessRestartPersistsRecords`（真实子进程终止→同 `DB_PATH` 重启）**；`go vet` 干净、`go test ./... -count=1` 全绿（L1 0.13s、L2 4.38s）、web vitest 458/458。**S6 勾选，GOAL-007 进度 `5/6 → 6/6`，六项成功标准全部达成**；`I-007-001/002/003/004` 全 verified。**本目标仍 `active`（未 `done`）**：关门需先做关门审计（self 或 `/audit`）+ 用户裁决 + Root R4 勾选。Root R4 未勾选（Root 保持 `3/5`）。后续意见从 A-010 起。

> **A-010 已响应、F-008 已闭合为 fixed（2026-08-02）**：A-010（independent · conditional · close-out）F-008 required（L2 未断言 PATCH `updatedAt` 跨进程毫秒精确保持）→ **fixed**——L2 `cmd/server/server_restart_test.go` 保留 POST/PATCH 响应的 `updatedAt`，Phase 2 对新建记录与 `rec-1` 分别 GET detail 断言毫秒精确一致（I-007-004 §3.6/§4）；`go vet ./cmd/server/` 干净、focused L2 PASS（4.32s）、`go test ./...`（apps/api）全绿。**关闭证据已请求 finding-closure 复核（`/audit`）**。R-003/R-004 保持 recommended 非阻断。GOAL-007 仍 `active / 6/6`，未 `done`；Root R4 未勾选（Root 保持 `3/5`）。后续意见从 A-011 起。

> **GOAL-007 已关门、Root R4 已勾选（2026-08-02）**：A-011/A-012（independent · finding-closure）确认 F-008 `fixed`；用户按 P-004 §3.1 裁决补 close-out self 审计 → **A-013（self · close-out · pass）**（成功标准 S1～S6 全 `6/6`、`I-007-001/002/003/004` verified、无开放 required；`go test ./...` 全绿 + web vitest 458/458）。**GOAL-007 已置 `done`**；**Root R4 已勾选（Root `3/5 → 4/5`）**。R5 与 VP-002 保持 open。

> **A-010 R-003/R-004 已闭合（2026-08-02 关门后补充）**：R-003 → fixed（`apps/api/README.md` 端点表阶段标注统一 R4）；R-004 → fixed（`e2e/schema-crud.spec.ts` 真实浏览器 list-edit-lifecycle CRUD + `login()` 后 `fetchMe` 修复 features 投影；Playwright `WEB_PORT`/临时 DB；本机 `WEB_PORT=9999` E2E **2 passed**）。GOAL-007 保持 `done / 6/6`；Root 保持 `4/5`；A-010 scope 无开放 required/recommended。

> **A-014 已响应、R5 已立项（2026-08-02）**：A-014（independent · finding-closure · pass）复核确认 A-010 R-003/R-004 `fixed` 关闭成立（README 端点表 R4 + 真实浏览器 CRUD E2E + login features；vitest 9/9 + Playwright 2 passed）；GOAL-007 保持 `done / 6/6`，Root `4/5`。**Root R5 信息门禁已关**：`I-005` → verified、`I-006` → closed（Root D-012/D-013：部署基线 A + 建议 15 分钟口径/复现方法 + 操作日志方案甲）；**`GOAL-008-r5-engineering-fork` 已立项（active / 0/5）**，登记 `I-008-001/002/003` 实施前 required。Root R5 检查点仍未勾选（待 GOAL-008 完成证据）。

> **GOAL-008 A-001 已响应（2026-08-02）**：A-001（independent · conditional）F-001 → **fixed**——GOAL-008 D-002 + D-001 修订 + S2 对齐：**Docker Compose 为 R5 必须交付和验收的第二启动路径**（S2 核心检查点、计入进度分母，非 S6 式可选加分项）；fork 用户可选本地双进程或 Compose；完整生产拓扑/CI-CD 仍非目标。R-001 → handled（I-005 附件 v0.2.1 时态清理）；R-002 → handled（`I-008-001/002` 信息表补最低收集清单）。`GOAL-008` 保持 `active / 0/5`；`I-008-001/002/003` 仍 open；Root `4/5`。

> **GOAL-008 A-002 已响应（2026-08-02）**：A-002（independent · finding-closure · pass）确认 F-001 `fixed` 关闭成立、R-001/R-002 handled；R-003 → handled（GOAL-008 概述 / Root 进度说明「R5 已立项待实施」/ I-005 附件 v0.2.2 §2 三处投影清理）。`GOAL-008` 保持 `active / 0/5`；`I-008-001/002/003` 仍 open；Root `4/5`。

> **GOAL-008 self 审计 + I-008-001 冻结（2026-08-02）**：A-003（self · pass）补齐立项与方案边界同 scope 自审（P-004 §3.1 闭环）。**`I-008-001` → verified（D-003 + [契约 v1.0.0](GOAL-008-r5-engineering-fork/attachments/I-008-001-engineering-contract.md)）**——S1/S2 方案冻结门禁解除。`I-008-002`/`I-008-003` 仍 open（阻断 S3/S4、S6 若实施）；`GOAL-008` 保持 `active / 0/5`；Root `4/5`。

> **GOAL-008 S1/S2 已实施（2026-08-02）**：`I-008-001` verified 后实施 S1（env 清单 + health/启动验证 + dev/prod 区分文档）与 S2（`apps/api/Dockerfile` + `apps/web/Dockerfile` + 根 `compose.yaml` + `nginx.conf` `/api` 反代 + CI `container-smoke`）；契约 **C-001～C-007 本机验证通过**（`docker compose up` → healthz/登录/`/me`/SPA fallback/重启与 down-up DB 持久化）。`GOAL-008` 进度 `0/5 → 2/5`；`I-008-002`/`I-008-003` 仍 open；Root R5 未勾选，Root `4/5`。
