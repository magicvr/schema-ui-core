---
title: 执行记录 · Schema 驱动通用数据适配层
status: active
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.5.0
---

# 执行记录 · GOAL-010

## 2026-08-03 · 立项

- 用户按 P-004 裁决：F-002-001 走「通用适配层改造」`fixed` 路径（Root D-014），本目标承接；F-002-002/003 归 `GOAL-009-a002-auth-form-fixes`。
- 建立五件套与 `attachments/`；登记实施前 required 信息项 `I-010-001`（资源契约）与 `I-010-002`（迁移策略）；goal-tree 已同步。
- 未修改任何产品代码；Root A-002 F-002-001 保持 `open`。
- **计划（非事实）**：收集并冻结 `I-010-001` → S1 方案冻结 → 冻结 `I-010-002` → S2 后端 → S3 前端 → S4 新实体 → S5 回归/审计。

## 2026-08-03 · S1 已实施（资源契约与方案冻结）

- 用户指令「实施 GOAL-010 S1」授权冻结；**D-002** 决策 + 附件 [I-010-001-schema-resource-contract.md](attachments/I-010-001-schema-resource-contract.md) **v0.1.0** 落盘。
- 契约核心（冻结）：`dataSource` 保持协议相对 URL（写端点由 action 显式声明；缺省 fail-closed 不再回落 `/api/records`）；统一 list envelope `{items,total,page,pageSize}` 跨资源、`items` 任意对象（解除 `RecordItem` 五字段白名单）；行键 `rowKey`（默认 `id`）；Go 资源注册表（id/path/listable/sortFields/qSearch/entity 接口/create·patch 字段/权限键派生）+ 通用 handler 工厂；records 注册 `/api/records` 保持 `records.read/write` 与全部错误码（**零对外 API 变更**）；错误 envelope `{error,message}` 不新增字段、NOT_FOUND = `{ID}_NOT_FOUND`（records 保持 `RECORD_NOT_FOUND`）。
- 迁移策略（I-010-002）：后端收敛为注册实例；前端 `RecordItem`/`RecordList` 一次性泛化、删除 URL 回落；不做新旧双轨；现有 fixture/emulator/测试形状保持。
- **`I-010-001` → verified、`I-010-002` → verified**（D-002 冻结；I-010-002 提前于最晚阶段 S3 关闭）；S1 方案冻结门禁解除。
- 未修改任何产品代码（S1 为文档冻结）。
- **计划（非事实）**：S2 后端通用资源 CRUD（注册表 + records 实例化，保持 T-API-01～13 全绿）。

## 2026-08-03 · A-001 F-001/F-002 响应 + S3 已实施（契约 v0.2.0 + 前端泛化）

- 用户指令「响应 F-001/F-002，走 fixed：补充契约语义和 S3 正反测试」；P-004 §3.1 裁决不补 self 审计。**D-003** 决策；[I-010-001 契约 **v0.2.0**](attachments/I-010-001-schema-resource-contract.md) 修订（§2 dataSource 执行规则、§3 rowKey 不变量）。
- **S3 · 前端通用适配层**实施：
  - `records.ts`：去除 `RecordItem` 五字段白名单 → 泛化 `ResourceItem`/`ResourceList` + 统一 envelope 解析；新增 `isValidDataSource`/`DATASOURCE_URL_PATTERN`（F-001）；`fetchRecords` 调用（认证）fetcher 前校验 dataSource；`DEFAULT_RECORDS_URL` 删除。
  - `schema-table.tsx`：`schemaTableDataSource` 对缺失/非法 dataSource 返回 `null`（fail-closed，不请求不渲染）；新增 `schemaTableRowKey`/`scalarRowKey`/`checkRowKeys`（F-002：非空唯一标量，无效响应停止渲染 + 禁行 action）；React key/选中态改用 rowKey 字段；`render.ts` `RenderTableNode.props` 增加 `rowKey`。
  - `use-records.ts`：移除 `DEFAULT_RECORDS_URL` 依赖。
- **测试（正反）**：`records.test.ts`（22，`isValidDataSource` 正反例 + 非法 dataSource 不触 fetch + 任意对象行）；`schema-table.test.tsx`（14，非 id `sku` 正例 + 缺失/重复/错误类型反例 + 缺 dataSource 不 fetch）；`representative-pages.integration.test.tsx` 错误文案 `records fetch failed` → `resource fetch failed`。
- **证据**：全量 `vitest run` **481/481** + `tsc -b` + `vite build` 干净。
- **F-001/F-002 → fixed 闭合**（见 [03-audit A-001 响应节](03-audit.md#a-001-响应--f-001--f-002-按-fixed-闭合2026-08-03)）；`I-010-001`/`I-010-002` 维持 `verified`。
- **S3 勾选，GOAL-010 `1/5 → 2/5`**（串行偏差留痕：S3 为纯前端、不依赖 S2，因 F-001/F-002 关闭证据在 S3 而先于 S2）。
- **计划（非事实）**：S2 后端通用资源 CRUD（注册表 + records 实例化，保持 T-API-01～13 全绿）→ S4 新实体 `catalog` 验证 → S5 回归/关门。S3 关闭证据可先请求窄 scope `/audit` finding-closure 复核。

## 2026-08-03 · S2 已实施（后端通用资源 CRUD）

- **`resources.go`**（通用资源注册表 + handler 工厂，I-010-001 §4）：`Resource` 描述符（id/path/listable/sortFields/qSearch/entity/create·patch 字段/权限键默认派生 `{id}.read`/`{id}.write`/NotFoundCode/NewID/OnWrite）+ `ResourceEntity` 接口（List/Get/Create/Update/Delete，行 = JSON map）+ `registerResource` 挂 list/create/detail/update/delete 五路由，统一 `requirePermission`（401/403）、4 KiB body 上限、`{error,message}` 写错误、`INTERNAL` 兜底、NOT_FOUND = `{ID}_NOT_FOUND`（records 显式 `RECORD_NOT_FOUND`）。
- **`records.go` 收敛为注册实例**：`recordsResource(st)` + `recordsEntity` 适配器——`recordToMap` 保持固定毫秒 `updatedAt`、`recordsOnWrite` 保持 `records.create/update/delete` 操作日志、`newRecordID` 保持 `rec-<hex>`；删除手写 `recordHandler`/`recordPatch`/`decodeCreateBody`/`validatePatch` 等。
- **`health.go`** Register 改为 `registerResource(mux, a, recordsResource(st))`。
- **零对外 API 变更**：全部现有 records 测试（T-API-01～13、权限 401/403、updatedAt 毫秒/单调递增、重启持久化、操作日志）保持全绿。
- **新增 genericity 测试 `resources_test.go`**：内存 `memEntity` 上注册合成 `catalog` 资源（不同 path/字段/sort 白名单/id 格式，**无手写 handler**）走通 list/create/detail/update/delete + 共享门禁（INVALID_SORT_FIELD / INVALID_CREATE_FIELD / INVALID_PATCH_FIELD / INVALID_CREATE_BODY / 权限 403）+ 默认权限键派生（无显式键 → `widget.read` 未授权 403）。
- **证据**：`go test ./...` 全绿（handler 7.2s 含新测试）+ `go vet` 干净；`gofmt` 仅 `internal/config/config_test.go` 为既有格式问题（非本轮改动）。
- **S2 勾选，GOAL-010 `2/5 → 3/5`**（串行恢复：S1/S2/S3 全勾选）。Root A-002 F-002-001 仍 `open`。
- **计划（非事实）**：S4 新实体验证（`catalog` fixture 仅改 Schema 接入，需 DB 迁移/种子 grants 注入/匿名 401·缺权限 403）→ S5 回归、审计与关闭。

## 2026-08-03 · S4 语义资源子目标已立项

- 用户确认 D-004：S4 不再以 `catalog` 等无普遍语义 fixture 作为终态验证；采用 `users + roles`，并将完整替换/新增工作拆入子目标 `GOAL-011-s4-semantic-admin-resources`。
- 已创建 GOAL-011 五件套与 `attachments/`，设定 `parent: GOAL-010-a002-schema-adapter`、五个等权检查点和三个开放 required 信息项 `I-011-001`/`I-011-002`/`I-011-003`。
- `I-010-001` 附件升至 v0.2.1 交接附注：S1～S3 技术契约不变；`catalog` 只保留为 genericity 历史示例，S4 产品终态由 GOAL-011 承接。
- 本目标 S4 文案已改为父级验收门，依赖 GOAL-011 完成证据；`progress` 保持 `3/5`，S4/S5 均未勾选。
- 未修改产品代码；records 仍在当前运行面；users/roles CRUD 尚未实施；Root A-002 F-002-001 继续 open。
- **计划（非事实）**：下一拍在 GOAL-011 收集并冻结 `I-011-001`/`I-011-002`，未冻结前不实施受影响范围。

## 2026-08-03 · GOAL-011 S4 证据交接（父级验收门就绪）

- GOAL-011 已完成 S1～S4（progress 4/5）：users/roles 后端闭环（S2）、records 产品运行面退场（S3，0006）、双语义实体 Schema 接入验证（S4，I-011-003 v0.2.0 冻结并 verified）。
- **本目标 S4 父级验收门证据链就绪**：
  - users 替换 records 默认代表实体、roles 作为第二语义资源：`apps/api/internal/handler/{users.go, roles.go}` + fixture `users.json`/`roles.json` + manifest users/roles 页。
  - 仅通过前端 Schema 接入列表/CRUD 页面、Renderer 主路径无修改：I-011-003 §3 基线 `adfe15a` 零 diff（`git diff --exit-code` exit 0）+ T-UI-10。
  - records 已从产品默认运行面按版本化策略退场：migration `0006 records_retire`（DROP TABLE + 清理权限/菜单行）+ per-pending 快照；产品代码 grep 无 `api/records` 等残留。
  - fresh fork / 既有库升级 / 进程级重启 / 401-403 双资源验收：I-011-003 v0.2.0 全维度 + S4 验收收据（见 GOAL-011 02-execution）。
- 完整证据见 [GOAL-011 五件套与附件](../GOAL-011-s4-semantic-admin-resources/)；GOAL-011 S5（回归审计与关门）完成后可勾选本目标 S4。
- **计划（非事实）**：GOAL-011 完成后评估本目标 S4 勾选与 S5 关门；Root A-002 F-002-001 关闭证据链仍待 GOAL-010 S5。

## 2026-08-04 · GOAL-011 A-012 响应交接（父级验收门保持未放行）

- GOAL-011 已记录原 S1～S5 检查点 `5/5`，但 A-012（independent · fail）新增 F-001～F-005 required/open，F-006 recommended/open；`/govern` 已将其状态恢复为 `active`，并登记 `I-011-004` 等待 F-003 产品边界裁决。
- 本交接只同步当前子目标状态，不把子目标 `5/5` 当作父级 S4 通过证据；GOAL-010 维持 `active / 3/5`，S4/S5 不勾选，Root A-002 F-002-001 不变。此前 S1～S4 交接段保留为当时事实。
- **计划（非事实）**：待用户裁决、整改实施、回归收据和限定范围 finding-closure 复审完成后，再按 D-004 重新评估 GOAL-010 S4 验收门。

## 2026-08-04 · GOAL-011 A-013 closure 交接（子目标关门，父级状态不变）

- GOAL-011 已按 D-006 / I-011-004 v0.1.0 完成 A-012 F-001～F-005 整改；候选提交 `fb5cd06` 与完整本地回归/Compose 收据已落盘，A-013（independent · finding-closure · pass）逐项确认五项 `fixed`、无新增 required，D-007 据此恢复子目标 `done / 5/5`。
- A-012 F-006 保持 `recommended / open / non-blocking`；GitHub-hosted Actions 未运行，子目标没有把本地 Linux Compose 收据写成远端 CI acceptance。
- 本交接只同步子目标当前关闭事实。GOAL-010 仍为 **`active / 3/5`**，S4/S5 均不勾选；Root A-002 F-002-001、Root 与 VP-002 状态不变。是否采纳子目标证据勾选本目标 S4，须由 GOAL-010 后续 `/govern` 独立评估。
