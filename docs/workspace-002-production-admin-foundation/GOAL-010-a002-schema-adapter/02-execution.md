---
title: 执行记录 · Schema 驱动通用数据适配层
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.2.0
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
