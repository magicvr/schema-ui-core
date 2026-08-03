---
title: 决策 · Schema 驱动通用数据适配层
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.3.0
---

# 决策 · GOAL-010

## D-001 · 采用通用适配层改造路径关闭 F-002-001

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004）：Root A-002 F-002-001（required / high）走「通用适配层改造」关闭路径（`fixed`），不降级 VP-002 主张。
- **决定**：将表格/表单 transport、字段模型与 response mapping 提升为 Schema 驱动的通用适配层，提供明确的后端资源契约；records 迁移为该适配层的注册实例；新业务实体仅通过修改 Schema 即可接入，无需修改 Renderer 主路径。
- **理由**：VP-002 产品级成功标准 6 明确要求「通过修改 Schema 新增业务页面，而无需修改前端 Renderer 主路径」；当前实现（`schema-table.tsx` / `records.ts` 硬编码）被 A-002 证伪为不满足标准 1/4/6。改造是保住 VP 核心主张的正向路径。
- **未选方案**：
  - **降级为单一 records 示例声明**：需收窄 VP-002 成功标准 6 并重新获得用户裁决，弱化「Fork 后可直接接业务」的 VP 核心价值。
  - **维持现状宣称完成**：A-002 已以代码证据 fail 证伪，且与 VP-002 文本冲突。
  - **先做评估再定路径**：用户已直接裁决改造，无需重复评估轮。
- **信息门禁**：`I-010-001` / `I-010-002` 为实施前 required，登记于 `00-meta`；在各自最晚阶段前由证据关闭，不得以本决策代替契约冻结。
- **影响**：Root A-002 F-002-001 保持 `open`，直至本目标 S1～S5 完成并有 `/audit` 关闭证据后按 `fixed` 闭合。
- **后续**：先收集并冻结 `I-010-001`（资源契约）与 `I-010-002`（迁移策略），再按路线图 S1 → S5 推进。

## D-002 · 冻结通用资源契约（S1 · 关闭 I-010-001 / I-010-002）

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004）：用户指令「实施 GOAL-010 S1」= 授权方案冻结；契约关键取舍随本决策留痕，后续如需修订走契约版本化（v0.1.x），不静默改写。
- **决定**：采纳 [I-010-001 通用资源契约 v0.1.0](attachments/I-010-001-schema-resource-contract.md)：
  1. **资源标识**：`dataSource` 保持协议相对 URL（写端点由 action 显式声明，不引入资源名映射层）；缺省 fail-closed（不再回落 `/api/records`）。
  2. **字段模型/envelope**：统一 list envelope `{items,total,page,pageSize}` 跨资源冻结；`items` 为任意对象（解除五字段白名单）；行键 `rowKey`（默认 `id`）；`columns[].field` 仅为展示/排序目标。
  3. **后端注册形态**：Go 资源注册表（id/path/listable/sortFields/qSearch/entity 接口/create·patch 字段/权限键派生）+ 通用 handler 工厂；records 注册到 `/api/records`，权限键 `records.read/write`，**对外 HTTP 契约与 I-007-001 逐项一致（零 API 变更）**。
  4. **错误 envelope**：`{error,message}` 跨资源冻结，不新增字段；通用错误码全资源共享；NOT_FOUND = `{ID}_NOT_FOUND`（records 保持 `RECORD_NOT_FOUND`）。
  5. **迁移策略（I-010-002）**：后端 records handler 收敛为注册实例（零对外变更）；前端 `RecordItem`/`RecordList` 固定解析一次性迁移为通用解析（`schema-table` 不再依赖固定形状、删除 URL 回落）；现有 fixture/emulator/测试形状保持；不做新旧双轨并行。
- **理由**：与 D-001「通用适配层改造」路径一致；records 兼容路径保住既有验收证据（T-API-01～13、schema-crud.test.tsx）不被破坏，同时解除前端固定解析（A-002 证伪点）。
- **未选方案**：
  - **dataSource 改为资源名 + 前端注册表解析**：引入前端第二套资源映射，与后端注册表双真相；URL 形态现状零成本继承。
  - **后端暴露 `/api/resources/{resource}` 统一前缀并迁移 `/api/records`**：破坏既有路径/fixture/测试，迁移面大且无收益（路径本身就是资源标识）。
  - **items 仍强制五字段白名单**：等于维持 A-002 证伪点，不达目的。
- **信息门禁**：`I-010-001` → **verified**（本契约 v0.1.0）；`I-010-002` → **verified**（§6 迁移策略随本决策冻结，最晚阶段 S3 提前关闭）。S1 方案冻结门禁解除，S2 实施可放行。
- **影响**：Root A-002 F-002-001 仍 `open`（S1～S5 完成并审计后闭合）；不修改任何产品代码（S1 为文档冻结）。
- **后续**：S2 后端通用资源 CRUD（注册表 + records 实例化）→ S3 前端泛化 → S4 新实体 `catalog` 验证 → S5 关门。

## D-003 · 响应 A-001 F-001/F-002：契约升版 v0.2.0 + S3 正反测试（fixed）

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004）：用户指令「响应 F-001/F-002，走 fixed：补充契约语义和 S3 正反测试」；P-004 §3.1 裁决**不补 self 审计**（L0 下 `fixed` 不强制自审/独立复审）。
- **决定**：采纳 [I-010-001 契约 **v0.2.0**](attachments/I-010-001-schema-resource-contract.md) 修订并实施 S3——
  1. **F-001（§2 执行规则）**：`dataSource` 必须匹配 `^/(?!\/)[^\s\\?#]*$`（单斜杠同源绝对路径，禁 `//`/scheme/空白/反斜杠/`?`/`#`）；`records.ts` 新增 `isValidDataSource`，`fetchRecords` 在调用（认证）fetcher **前**校验；`schema-table.tsx` `schemaTableDataSource` 对缺失/非法返回 `null` → fail-closed（不请求、不渲染、可观察错误）；删除 `DEFAULT_RECORDS_URL` 回落。
  2. **F-002（§3 行键不变量）**：`rowKey` 为直接字段名（默认 `id`）；每行非空且唯一 string/finite-number 标量；无效响应（缺失/空/非标量/重复）停止渲染数据、禁止行 action 与选中，渲染可观察错误。`schema-table.tsx` `schemaTableRowKey`/`scalarRowKey`/`checkRowKeys` 实施；`RenderTableNode.props` 增加 `rowKey`。
  3. **泛化**：`records.ts` 去除 `RecordItem` 五字段白名单（`ResourceItem`/`ResourceList` + 统一 envelope 解析）；`readRecordApiError`/`buildRecordsQuery`/`create·update·deleteRecord` 保留（泛化 body）。
- **fixed 关闭证据**：契约 §2/§3（v0.2.0）+ 本决策 + S3 正反测试——`records.test.ts`（22 用例，含 `isValidDataSource` 正反例与 `fetchRecords` 非法 dataSource 不触 fetch）、`schema-table.test.tsx`（14 用例，含非 id `sku` 正例与缺失/重复/错误类型反例）+ 全量 `vitest run` **481/481** + `tsc -b`/`vite build` 干净。F-001/F-002 按 `fixed` 闭合（见 03-audit 响应节）。
- **未选方案**：
  - **契约仅补文档、S3 实施另轮**：F-001/F-002 必改要求 S3 正反测试作为关闭证据，仅文档不构成 `fixed`。
  - **同轮实施 S2 后端**：用户裁决本轮范围为「契约 + S3」（S2 下一轮）。
- **信息门禁**：`I-010-001` 维持 `verified`（v0.2.0 为响应修订，不改变冻结结论）；`I-010-002` 维持 `verified`。
- **影响**：**S3 检查点达成，GOAL-010 `1/5 → 2/5`**（串行偏差留痕：S3 为纯前端、不依赖 S2，因 F-001/F-002 关闭证据在 S3 而先于 S2 实施；S2 下一轮）；Root A-002 F-002-001 仍 `open`（待 S2～S5 完成 + S4 新实体验证 + 关门审计后闭合）。
- **后续**：S2 后端通用资源 CRUD；S3 关闭证据可请求窄 scope `/audit` finding-closure 复核后再推进 S2 实施门禁。

## D-004 · S4 拆出语义资源子目标并采用 users + roles

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**：确认新建目标承接 S4；将 `records` 替换为对绝大多数系统具有实际语义的设计，并把 S4 的“新增”改为第二种同样具有实际语义的资源。用户同意采用推荐的 `users + roles` 组合。
- **决定**：
  1. 新建 `GOAL-011-s4-semantic-admin-resources`，`parent: GOAL-010-a002-schema-adapter`；以 users 替换 records 默认代表实体，以 roles 作为新增的第二资源。
  2. 本目标 S4 改为父级验收门：GOAL-011 完成、records 按版本化策略从当前产品默认运行面退场、users/roles 均完成 Schema 驱动列表/CRUD 闭环且 Renderer 主路径无修改后，才可勾选。
  3. “只修改 Schema 接入”限定为后端资源已经显式注册后的前端页面接入；后端持久化、权限与领域规则不由 Schema 自动生成。
  4. D-002/I-010-001 的 records 零 API 变更保持为 S1～S3 的历史事实；S4 终态由 GOAL-011 的新契约和迁移策略承接，不静默改写既有迁移或历史治理记录。
- **理由**：records 已从演示基线扩张到迁移、种子、权限、菜单、操作日志、fixture 与回归主线；彻底替换并新增第二个真实资源跨越多个门禁域，不能诚实压缩为本目标一个实现检查点。users/roles 已属于现有真实认证与 RBAC 数据域，对 fork 项目有直接价值，并保持在 VP-002 的最小权限边界内。
- **未选方案**：
  - **仅修改 S4 文案并直接开工**：缺少领域安全、迁移和双资源验收的信息门禁。
  - **继续 records + catalog**：保留并增加无普遍语义的示例域，未解决 fork 代码污染。
  - **roles + menu_items**：风险较低，但账户管理价值与认证链闭环弱于 users + roles。
- **信息门禁**：GOAL-011 登记 `I-011-001`/`I-011-002`/`I-011-003` required，初始均 open；只允许 S1 收集，未放行产品实现。
- **影响**：GOAL-010 保持 `active / 3/5`，S4/S5 未勾选；Root A-002 F-002-001 继续 open，Root/VP-002 关门继续阻断。
- **后续**：由 GOAL-011 先冻结 users/roles 领域契约与 records 退场策略，再按其路线图实施；完成后回到本目标评估 S4。
