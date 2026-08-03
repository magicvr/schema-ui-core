---
title: 审计台账 · Schema 驱动通用数据适配层
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.3.0
---

# 审计台账 · GOAL-010

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-03 | S1 资源契约与 `I-010-001` 信息门禁 | conditional | **F-001 / F-002 → `fixed`**（契约 v0.2.0 + S3 正反测试，见 A-001 响应节）；不改变 status/progress |

> 本目标为 A-002 响应载体；Root 层面的 A-002 意见与响应记录见 [Root 03-audit](../GOAL-001-production-admin-foundation/03-audit.md)。

## A-001 · S1 资源契约与 I-010-001 独立审计（2026-08-03）

- **source**：independent
- **auditor**：Codex · `$audit`
- **类型 / scope**：design-plan + information-readiness；仅核对 GOAL-010 S1 的资源契约、`I-010-001` / `I-010-002` 作为实施前 required 信息门禁的充分性，以及它们与 Root A-002 F-002-001 的关闭路径是否一致。未复判 S2～S5 的实现、回归或 Root finding 关闭。
- **verdict**：**conditional**。通用适配层的方向、兼容策略和 S1 / 实施事实边界合理；但 `dataSource` 的认证请求边界与通用行键的输入不变量尚未冻结。`I-010-001` 因而不能在未补正前作为完整、无条件的 S2 实施输入。

### 范围与区间

- 工作区为 `workspace-002-production-admin-foundation`；`workspace.md` 的 Root、canonical root、`primary_plan: VP-002-production-admin-foundation` 与本目标 `parent` 一致，`shared_materials_catalog: none`。本意见未使用共享资料作为事实或关闭证据。
- 已核对 GOAL-010 五件套、D-002、[I-010-001 资源契约 v0.1.0](attachments/I-010-001-schema-resource-contract.md)、Root A-002 / D-014，以及契约直接指向的 Web 表格 transport、认证 fetch、DataTable 和 API records handler / 既有 records 契约。
- S1 的执行记录明确为决策和文档冻结，且明确未修改产品代码；本意见不将该事实误写为 S2～S5 已实施，也不将其作为 Root F-002-001 的关闭证据。

### 成果（有证据）

| 核对项 | 判断 | 证据 |
|--------|------|------|
| 关闭路径的方向 | **合理**：保留 VP-002 的“新增业务页面不改 Renderer 主路径”主张，以通用适配层而非降级 records 单例响应 Root F-002-001。 | `00-meta.md` S1～S5；`01-decision.md` D-001；Root A-002 F-002-001 / D-014 |
| API 兼容与资源边界 | **合理**：records 保持 `/api/records`、`records.read` / `records.write`、既有 envelope 和错误码；资源必须显式注册，未扩张为无约束的任意表 CRUD。 | 契约 §4～§6；I-007-001 v0.2.0 |
| 迁移边界 | **合理**：后端收敛为注册实例、前端一次性泛化、现有 records fixture / emulator 形状保持，且不以双轨掩盖迁移。 | 契约 §6；D-002 |
| S1 与实施的边界 | **清楚**：附件明确排除 S2～S5 成品，执行记录也明确未改产品代码。 | 契约 L15～L17；`02-execution.md` S1 节 |
| 作为完整实施门禁的精确度 | **不足**：存在下列两个 required 缺口；在响应前不能把 `I-010-001: verified` 扩写为对 S2～S4 的无条件放行。 | F-001、F-002 |

### Findings

#### F-001 · `dataSource` 的同源路径与认证边界未冻结（required / medium）

- **证据**：契约 §2 仅称 `table.props.dataSource` 为“协议相对 URL 字符串”，举例 `/api/records` / `/api/catalog`，但未定义可接受的语法、拒绝 `//host` 或绝对 URL 的行为，也未要求在调用认证 transport 前验证。相比之下，已固定的 Node `DataRef.url` 明确使用单斜杠相对路径约束（`docs/schemas/node.schema.json` 的 `^/(?!/)[^\\s\\\\]*$`）。当前 `main.tsx` 将 `authFetch` 作为表格 fetcher 注入，而 `authFetch` 在解析 endpoint 前即向任意输入附加 `Authorization: Bearer`；当前 `schema-table.tsx` 也只检查 `dataSource` 非空。
- **影响**：S3 若只按现有文字泛化，将把页面 Schema 控制的字符串直接交给认证 fetch。绝对 / protocol-relative URL、查询串拼接或其他非端点值的处理将由实现者临时决定，并可能扩大 Bearer 请求的目标范围；这不符合本目标要求的明确 transport 契约与 fail-closed 目标。
- **必改**：经 `/govern` 将契约修订为可执行的单一规则，例如仅允许单斜杠同源绝对路径、禁止 `//` / scheme / 空白 / 反斜杠，并明确 query / fragment 的处理；S3 必须在认证 fetch 前复核该规则，并补正反例测试。修订后应留下对应 D-002 响应和 finding-closure 证据。
- **影响门禁**：`I-010-001` 的完整实施输入；S2 及后续 S3/S4 不得据当前“verified”状态作无条件放行。

#### F-002 · `rowKey` 的响应不变量和失败行为未冻结（required / medium）

- **证据**：契约 §3 允许 `items` 为任意 JSON 对象，并仅把 `table.props.rowKey` 写为默认 `id` 的 string 配置；未规定每行在该字段上必须有何种可接受值、唯一性、缺失 / 非标量 / 重复时的 fail-closed 行为。`DataTable` 要求 `rowKey` 回调返回 string，当前 `SchemaTable` 则硬编码 `row.id` 作为 React key、选中态和行 action 的关联键；运行时 Node Schema 也未为 table props 提供该约束。
- **影响**：新资源可具有非 `id` 的主键、数值键或缺失 / 重复键。没有冻结输入与拒绝语义时，S3 的泛化实现会自行选择字符串化、静默回落或继续渲染，造成重复 React key、选中行 / edit target 错配或对错误行发起写操作，不能证明 S4 的“仅改 Schema 接入”是可靠的。
- **必改**：冻结 `rowKey` 的解析范围与数据要求（至少直接字段名、每行非空且唯一的键值、允许的 JSON 标量类型），并规定无效响应必须停止渲染数据和禁止行 action；S3 测试须覆盖非 `id` 的新实体键、无键、重复键和错误类型。
- **影响门禁**：`I-010-001` 的字段模型 / response mapping 部分；S2 及后续 S3/S4 不得据当前“verified”状态作无条件放行。

### 必改项汇总

- **F-001**、**F-002** 均为 required / medium，且尚未按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合。
- 它们不否定 S1 已完成“有版本化方案文档”的事实，也不等同于 S2～S5 已实施；但使 `I-010-001` 作为完整实施门禁的 `verified` 结论需要由 `/govern` 响应和补正证据后才能依赖。
- Root A-002 F-002-001 继续保持 `open`；本意见不关闭或重编号该 Root finding。

### 与既有意见的异同

- 与 Root A-002 的 F-002-001 一致：当前 records 专用实现不能支撑“只改 Schema 接入新业务实体”；本意见仅审查其已选修复路径的方案充分性，不声称已复核修复完成。
- 本目标此前无 self 或 independent 正式意见，故不存在同 scope verdict 冲突。

### 结论 + 建议给编排器 / 用户的下一步

- **conditional**：`I-010-001` 的架构取向合理，尤其是 URL 直接作为资源标识、显式后端注册和 records 零 API 变更；但它尚未冻结两个会直接影响认证 transport 与新实体 CRUD 安全性的边界。
- 建议 `/govern` 先按 P-004 处理 F-001 / F-002：建议路径为 `fixed`，将契约升版并记录精确语义、D-002 响应与 S3 正反验证要求；随后可请求一次窄 scope `/audit` finding-closure 复核，再推进受影响的实施门禁。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、派生 `progress`、方案正文或 `goal-tree.md`；响应、finding 闭合与阶段推进由 `/govern` 处理。

## A-001 响应 · F-001 / F-002 按 fixed 闭合（2026-08-03）

- **响应路径**（P-004 用户裁决）：`fixed`。用户指令「响应 F-001/F-002，走 fixed：补充契约语义和 S3 正反测试」；P-004 §3.1 裁决**不补 self 审计**（L0 下 `fixed` 不强制自审/独立复审）。
- **F-001 → fixed**（required / medium）：契约 [v0.2.0 §2](attachments/I-010-001-schema-resource-contract.md) 冻结 `dataSource` 单斜杠同源执行规则（`^/(?!\/)[^\s\\?#]*$`，禁 `//`/scheme/空白/反斜杠/`?`/`#`；query 由 `buildRecordsQuery` 追加、fragment 禁止）；`records.ts` 新增 `isValidDataSource`/`DATASOURCE_URL_PATTERN`，`fetchRecords` 在调用（认证）fetcher **之前**校验；`schema-table.tsx` `schemaTableDataSource` 对缺失/非法 `dataSource` 返回 `null` → fail-closed（不请求、不渲染、可观察错误）；`DEFAULT_RECORDS_URL` 回落删除。**反例测试**：`records.test.ts` `isValidDataSource` 正反例（`//host`/scheme/相对路径/空白/反斜杠/`?`/`#`/空串）+ `fetchRecords` 非法 dataSource 拒绝且不触 fetch；`schema-table.test.tsx` 缺 dataSource 不 fetch（`fetcher` 未调用）。
- **F-002 → fixed**（required / medium）：契约 [v0.2.0 §3](attachments/I-010-001-schema-resource-contract.md) 冻结 `rowKey` 为直接字段名（默认 `id`），每行非空且唯一 string/finite-number 标量；无效响应（缺失/空/非标量/重复）**停止渲染数据、禁止行 action 与选中态**、渲染可观察错误。`schema-table.tsx` 实施 `schemaTableRowKey`/`scalarRowKey`/`checkRowKeys`；`RenderTableNode.props` 增加 `rowKey`。**正反测试**：`schema-table.test.tsx` 非 id `sku` 正例 + 缺失/重复/错误类型反例（`tbody` 0 行 + `role=alert` 错误）。
- **泛化**：`records.ts` 去除 `RecordItem` 五字段白名单（`ResourceItem`/`ResourceList` + 统一 envelope 解析）；`readRecordApiError`/`buildRecordsQuery`/`create·update·deleteRecord` 保留（泛化 body）。`schema-crud.test.tsx`（T-UI-10 等）与集成测试保持全绿。
- **证据**：契约 v0.2.0（§2/§3/§10）；GOAL-010 [D-003](01-decision.md)；`apps/web/src/renderer/records.ts`、`schema-table.tsx`、`render.ts`（`rowKey` 字段）；`records.test.ts`（22）、`schema-table.test.tsx`（14）；全量 `vitest run` **481/481** + `tsc -b` + `vite build` 干净。
- **状态**：**F-001、F-002 均按 `fixed` 合法闭合**（P-003 三路径）；同 scope 无冲突意见；`I-010-001`/`I-010-002` 维持 `verified`。S3 检查点达成（GOAL-010 `1/5 → 2/5`，串行偏差留痕：S3 纯前端、不依赖 S2）。
- **可选加固**：S3 关闭证据可请求一次窄 scope `/audit` finding-closure 复核，再推进 S2 实施门禁；Root A-002 F-002-001 仍 `open`，不因此闭合而关闭。
