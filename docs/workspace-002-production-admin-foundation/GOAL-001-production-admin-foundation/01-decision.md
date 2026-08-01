---
title: 决策 · 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.1.2
---

# 决策 · GOAL-001

## D-001 · 以新 delivery 工作区承接 VP-002

- **日期**：2026-08-01
- **决定**：建立 `workspace-002-production-admin-foundation` 与 Root `GOAL-001-production-admin-foundation`，将 VP-002 从 `planned` 激活为 `active`，并把本工作区设为其当前唯一 lead workspace；仓库 `primary_workspace` 保持 `workspace-001-mvp-admin-foundation`。
- **依据**：用户明确要求开启新工作区和 Root 承接 VP-002，并确认了工作区与 Root 命名。VP-001 与旧 Root 均已关闭/完成，新波次需要独立 canonical scope 和目标树。
- **边界**：不重开 VP-001 或旧 Root；不建立跨工作区 `parent`；只通过 Q2 路径引用已冻结历史基线。

### 未选方案

- **复用 workspace-001 并重开旧 Root**：会混合已关闭波次与新实施事实，破坏状态边界。
- **把新工作区设为 primary**：当前没有长期目的或仓库北极星换代，不应改写 Charter 的 `primary_workspace`。
- **把 VP-002 作为旧 Root 的子目标**：跨工作区 parent 被协议禁止，也无法提供独立交付树。

## D-002 · 采用五阶段串行纲领路线图

- **日期**：2026-08-01
- **决定**：Root 采用 `Renderer → 认证 → 持久化权限 → CRUD → 工程化与关门` 五个等权检查点，纲领阶段原则上串行，阶段内部可按依赖并行。
- **原因**：该顺序遵循 VP-002 的价值链，同时把原第三阶段拆成可独立验证的权限持久化、CRUD 和工程交付，避免一个阶段承载过多门禁。
- **执行约束**：只在当前阶段边界和 required 信息项就绪后创建具体子目标；`progress` 只按完成检查点数派生，不用于放行。

## D-003 · 先登记未知，再按阶段关闭

- **日期**：2026-08-01
- **决定**：将协议实施差量、认证机制、持久化模型、代表性 CRUD 实体及部署/fork 验收口径登记为阶段 required 信息项；操作日志保持 non-blocking。
- **原因**：这些未知不妨碍 Root 立项，但会改变对应阶段的方案或验收，必须在最晚门禁前以证据关闭或经用户书面接受有界 residual。

## D-004 · 采用 I-001 差量矩阵作为 R1 方案边界输入

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：
  1. 以附件 [I-001-implementation-gap-matrix.md](attachments/I-001-implementation-gap-matrix.md) 作为 `I-001` 的可核对答案，并将 `I-001` 标为 `verified`。
  2. **不改写** `I-PROTO-001 v0.1.3` 覆盖表；本波次仅继承其 include / include-partial / exclude。
  3. 将 R1 方案边界冻结为矩阵 §4：**In** = Schema 加载、默认 `RenderPage` 主路径、加载时结构校验、白名单节点/表单、`$context` reactions、代表性 Node 页、统一失败面、手写示例降级；**Reuse** = 现有 shell/fixture pin/演示 records API/静态 dev session（非生产身份）；**Out** = 真实认证、持久化 IAM、覆盖扩域、D-UPLOAD、`$deps` reactions、业务模块、fork/Docker 关门证据。
  4. 本回合**不**创建 R1 子目标、**不**修改产品代码；矩阵闭合只解除「R1 方案冻结前」的 `I-001` 信息门禁，**不**勾选 Root 路线图 R1、不宣称 Renderer 产品化完成。
- **理由**：用户确认 `/govern` 路径 A（仅文档关闭 I-001）。扫描显示库级/fixture 能力大量已有，但产品默认页面仍走 `EXAMPLE_PAGES`，`schemaUrl` 未驱动渲染——R1 主差量在主路径产品化，而非重新冻结协议。
- **关联信息项**：`I-001` → `verified`；证据路径见矩阵 frontmatter 与 §5。
- **后续**：可按矩阵候选拆分创建 R1 子目标并进入实施；I-002～I-005 仍分别阻断 R2～R5。

### 未选方案

- **跳过矩阵直接批量建 R1～R5 子目标**：违反 P-001/P-005；I-001 为 R1 方案冻结 required。
- **把现有手写示例与 fixture 通过直接标 R1 完成**：与 VP-002「Schema 驱动为默认主路径」成功标准冲突。
- **在本回合扩大 v0.1.3 或实现真实认证**：分别需要新覆盖决策与 R2 信息/方案，超出路径 A。

## D-005 · R1 拆为三个子目标（加载 / 主路径 / 代表性页）

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：在 I-001 已 verified 且 D-004 方案边界生效后，于本工作区创建三个 **R1 阶段内** 子目标（均可 `active`，阶段内可并行准备）：
  1. `GOAL-002-r1-schema-load-validate` — Schema 加载、结构校验、统一错误面  
  2. `GOAL-003-r1-default-render-path` — 默认 `RenderPage` 主路径与 `EXAMPLE_PAGES` 降级（硬依赖 002）  
  3. `GOAL-004-r1-representative-node-pages` — 代表性列表/表单/组合 Node 页与回归证据（完整主路径证明依赖 002+003）  
- **理由**：用户明确要求按 D-004 创建「加载 + 主路径 + 代表性 Node 页」；拆分对齐矩阵 §4 候选，避免单目标混杂门禁。
- **边界**：不创建 R2～R5 子目标；不勾选 Root `progress` R1；不改产品代码（本决策仅立项）。
- **后续**：优先实施 GOAL-002；GOAL-003 切换默认分支须 002 可测；GOAL-004 资产可先行。

### 未选方案

- **合并为一个 R1 大子目标**：验收与并行困难，与矩阵建议拆分不一致。
- **只建加载、不做主路径/页面目标**：无法关闭 VP-002 阶段 1「默认 Schema 驱动」主张。
