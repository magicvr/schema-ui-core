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
