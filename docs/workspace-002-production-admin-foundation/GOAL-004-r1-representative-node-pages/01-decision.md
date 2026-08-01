---
title: 决策 · R1 · 代表性 Node 页面与回归证据
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.2.0
---

# 决策 · GOAL-004

## D-001 · 代表性集合 = 列表 + 表单 + 组合/详情，均在白名单内

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：R1 代表性页面至少覆盖三类：列表向、表单向、组合/详情向；节点 type 与 form 控件不得超出 `I-PROTO-001 v0.1.3` §5 白名单。
- **理由**：对齐 VP-002 阶段 1「列表、详情、表单和组合」验收建议与 I-001 矩阵。
- **未选**：只做一个 hello-world section——不足以证明 Admin 主路径。

## D-002 · 完整主路径证明绑定 002+003；资产可先行

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：允许先写 Node JSON 与 `RenderPage` 单测；宣称「默认主路径上改 Schema 即出现页面」须在 GOAL-002/003 可用后用集成证据闭合。
- **理由**：拆分交付物，避免阻塞页面设计，同时防止用内嵌文档冒充产品主路径完成。

## D-003 · 页面资产 = 5 份迁移 Schema 文档（关闭 I-004-001）

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：本目标页面资产 = 5 份**迁移 Schema 文档**（改写既有 5 个手写示例语义，非新增独立资源），落于 Go `//go:embed fixtures/schema/*.json`，经默认主路径 `page.route → page.schemaUrl → GET /api/schema/<pageId> → loadPageDocument → RenderPage` 渲染：
  1. `data-table` — 列表向（table 节点 + columns + dataSource）
  2. `search-form-table` — 列表向 search+table 结构（form 搜索控件 + table）
  3. `form-controls` — 表单向（白名单控件）
  4. `form-with-reactions` — 表单向 + `$context` reactions
  5. `list-edit-lifecycle` — 组合/详情向（tabs + section + recordView + text + form）
- **依据**：GOAL-003 **D-003**（用户确认「迁移为 Schema」；`I-003-001` closed）；I-001 矩阵 §4 页面资产候选；本 meta `I-004-001` 结论列。
- **边界**：页面语义只在 `I-PROTO-001 v0.1.3` §5 白名单 + R1 能力（form/table/actionButton/`$context` reactions）内表达；search 表单→表格查询绑定、list-edit 完整 CRUD 生命周期属 R1 **Out**，不作为本目标承诺。
- **关联信息项**：`I-004-001` → **closed**（证据 = 本决策 + 5 fixtures）。
- **后续**：进入实施，落地 5 文档 + 默认表格数据注入 + 回归测试。

## D-004 · 列表数据复用 `/api/records`（关闭 I-004-002，联动 I-003-002）

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：R1 列表页数据继续使用既有 `GET /api/records` 演示 API（D-DATA）。默认表格数据注入读取 table 节点 `props.columns` / `props.dataSource`，经 records 列表包络（`items/total/page/pageSize`）渲染 `DataTable`，含 loading/error/empty 与列排序。
- **理由**：I-001 矩阵 D-DATA 允许 R1 reuse；避免在 R1 引入新数据契约或持久化。
- **边界**：静态 dev 数据 + 进程内变更，非生产持久化（R3/R4）；search 表单输入不绑定表格查询（R1 无 form→data 绑定原语）。
- **关联信息项**：`I-004-002` → **closed**（证据 = 本决策）；GOAL-003 `I-003-002`（列表页数据注入）随本实施联动关闭。
