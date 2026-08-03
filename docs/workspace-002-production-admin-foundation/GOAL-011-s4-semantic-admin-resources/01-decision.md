---
title: 决策 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-010-a002-schema-adapter
version: 0.1.0
---

# 决策 · GOAL-011

## D-001 · 采用 users + roles，替换 records 默认语义并独立立项

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**：同意创建新目标承接 GOAL-010 S4；将 `records` 替换为对绝大多数系统具有实际语义的设计，并把 S4 的“新增”改为第二种同样具有实际语义的资源。用户确认采用编排器推荐的 `users + roles` 结构。
- **决定**：
  1. `users` 替换 `records` 作为默认代表实体；`roles` 作为新增的第二个语义资源。
  2. 创建 `GOAL-011-s4-semantic-admin-resources`，`parent: GOAL-010-a002-schema-adapter`，以五个顺序检查点承载契约、后端双资源、records 退场、Schema 集成与审计交接。
  3. GOAL-010 S4 改为本目标的父级验收门：本目标完成后，才判断 records 是否已从当前产品默认运行面退场、users/roles 是否完成双资源闭环，以及前端 Renderer 主路径是否保持不变。
  4. “Schema-only 接入”只约束前端页面接入：后端资源注册、持久化、权限与领域规则仍须显式实现和验证。
  5. users/roles 只交付最小可用 Admin 管理边界；密码哈希、refresh token 等敏感字段不得进入通用响应，完整 IAM、SSO、SCIM 与复杂策略编排保持非目标。
  6. GOAL-010 D-002/I-010-001 的 records 零 API 变更继续作为 S1～S3 的历史实施与回归事实；本目标的终态替换必须通过新契约和迁移版本演进，不静默改写既有决策、迁移或审计记录。
- **理由**：
  - `records` 已证明 transport/CRUD 功能，但其 `name/status/owner` 模型没有稳定业务语义；fork 项目需要主动移除 API、表、种子、权限、菜单、操作日志、fixture 与测试，构成产品基线污染。
  - users 与 roles 已属于本仓真实认证/RBAC 持久化域，绝大多数 Admin 系统均有实际价值，并直接服务 VP-002 的真实认证、最小权限与 fork 即用边界。
  - 替换涉及多个独立门禁域和至少两个可独立验收交付块，已超过 GOAL-010 单个 S4 检查点可以诚实承载的范围，按 P-001 建立子目标更可追踪。
- **未选方案**：
  - **仅在 GOAL-010 S4 内直接实施**：会把领域契约、数据迁移、双资源交付和退场清理压进一个检查点，门禁与证据混杂。
  - **继续使用 records，只新增 catalog**：仍要求 fork 项目清理两套无普遍语义的示例域，不能解决用户指出的代码污染。
  - **roles + menu_items**：实现风险较低，但缺少账户管理闭环，通用价值与 VP-002 真实认证链的结合弱于 users + roles。
- **信息门禁**：`I-011-001`、`I-011-002`、`I-011-003` 均为 required，初始 `open`；当前只放行 S1 信息收集，不放行 S2～S4 产品实施或验收。
- **影响**：本目标 `active / 0/5`；GOAL-010 保持 `active / 3/5`、S4 未勾选；Root A-002 F-002-001 继续 `open`，Root/VP-002 关门继续阻断。
- **后续**：先收集并提交 `I-011-001`/`I-011-002` 的方案裁决，冻结 S1 后再进入任何 users/roles 或 records 退场代码变更。

## D-002 · 冻结 S1：users/roles 领域契约与 records 退场策略

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004.4）：用户于契约冻结裁决点逐项确认三项关键取舍（均采纳编排器推荐）：
  1. **通用工厂 + 最小契约扩展**：users/roles 均走通用资源工厂五路由；扩展 `Resource.JSONFields`（任意 JSON 值字段透传，承载 users.roles）与 `DomainError{Status,Code,Message}`（实体返回领域错误，工厂逐字映射）；对 I-010-001 §5「不引入 409」做限定范围偏离（仅账号域 409，envelope 形状不变）。
  2. **操作日志纳入**：migration `0005` 重建 `operation_log` event CHECK，新增 `users.*`/`roles.*` 事件（保留 `records.*`/`auth.*` 历史合法值）；users/roles 写路径挂 `OnWrite`。
  3. **records 硬退场 DROP TABLE**：migration `0006` `DROP TABLE records` + 清理 records 权限/菜单行；既有库升级自动 `pre-v0006` 快照兜底，records 数据随表删除（可由快照恢复）。
- **决定**：采纳两份版本化契约冻结 S1——
  1. [I-011-001-users-roles-contract.md](attachments/I-011-001-users-roles-contract.md) **v0.1.0**：users/roles 资源契约（公开字段、敏感字段隔离、角色分配、self/最后管理员保护、system role 保护、grant 约束、权限键/菜单/操作日志、错误码）、最小 IAM 边界、通用工厂最小扩展。
  2. [I-011-002-records-retirement.md](attachments/I-011-002-records-retirement.md) **v0.1.0**：records 足迹盘点、fresh install 与 in-place upgrade 迁移矩阵（0005/0006）、硬退场数据处置、代码/种子/fixture/前端退场动作、S3 验收口径。
- **理由**：三裁决点均为「先例契约改写/数据处置/范围取舍」，须用户书面确认而非编排器静默推断；采纳推荐路径保持 S2「通用工厂之上」、审计链一致与 fork 基线干净。
- **未选方案**：users 自定义 handler（不改工厂契约但 S2 主张打折、双套门禁逻辑）；操作日志不纳入（省 0005 但账号变更无审计留痕）；软退场保留死表（fork 基线不干净、数据处置不明）。
- **信息门禁**：`I-011-001` → **verified**（契约 v0.1.0 + 本决策）；`I-011-002` → **verified**（契约 v0.1.0 + 本决策）。S1 方案冻结门禁解除，S2 实施与 S3 退场可放行。`I-011-003` 保持 `open`（最晚 S4 前）。
- **影响**：**S1 检查点达成，GOAL-011 `0/5 → 1/5`**；GOAL-010 保持 `active / 3/5`；Root A-002 F-002-001 仍 `open`，Root/VP-002 关门继续阻断。未修改任何产品代码（S1 为文档冻结）。
- **后续**：S2 后端 users/roles 资源闭环（通用工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403）→ S3 records 退场 → S4 双资源 Schema 接入 → S5 回归审计关门。

## D-003 · 响应 A-002：契约修订至 v0.2.0（fixed 闭合 F-001/F-002 + 采纳 F-003~F-006）

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004 §3.2）：A-001（self · pass）与 A-002（independent · conditional）同 scope verdict 冲突；用户裁决闭合路径「**全部 fixed**」——修订两契约 + 采纳 recommended，不补 residual/overruled。
- **决定**：
  1. **I-011-001 v0.1.0 → v0.2.0**：§7 补 `ResourceEntity` Create/Update/Delete 的 `account.User` actor 通道（F-001，SELF_OPERATION/LAST_ADMIN 可诚实实现）+ DomainError 检查先于 ErrNotFound/INTERNAL（F-006）；§2.3 禁 API 路径复用 `linkUserRole`→`ensureRole` 隐式建角色（F-004）；§3.0 冻结 roles 公开响应形状 `system:boolean` + 毫秒时间戳（F-003）。
  2. **I-011-002 v0.1.0 → v0.2.0**：§2.1/§2.3 快照语义改为「每个待应用数据变更迁移前快照」（至少 0006 前强制），0005+0006 同批时 `pre-v0006` 必然存在（F-002）；§5 验收句对齐。
  3. **GOAL-010 D-005 / I-010-001 v0.2.2**：父契约 §5 追加账号域 409 限定扩展注记（F-005，消除跨目标双真相）。
  4. A-001 F-001（password 长度）与 F-002（fixture 文案）维持 recommended，随 S2/S3 落实（F-006 承接）。
- **fixed 关闭证据**：两契约 v0.2.0 + 本决策 + GOAL-010 D-005/I-010-001 v0.2.2；A-001/A-002 差异经此趋同（见 03-audit 响应节）。
- **理由**：F-001/F-002 均为真实可核对缺口（工厂无 actor 通道则 self 保护不可诚实实现；快照机制与验收字面不一致），修文档成本低、无 residual 必要；recommended 采纳后 S2/S4 金标准更明确。
- **未选方案**：residual（缺口小、应修）；overruled（拒绝合理必改项无依据）；仅闭 required 延后 recommended（用户裁决全部采纳，提高 S2 可实施性）。
- **信息门禁**：`I-011-001`/`I-011-002` 维持 `verified`（v0.2.0 为响应修订，不改变冻结结论）；`I-011-003` 保持 open。
- **影响**：GOAL-011 保持 `active / 1/5`，S1 契约以 v0.2.0 为准；A-002 conditional 经 F-001/F-002 闭合后与 A-001 pass 趋同；S2 实施门禁保持解除；Root A-002 F-002-001 仍 open。
- **后续**：进入 S2 后端 users/roles 资源闭环（按 v0.2.0 契约落地工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403）。
