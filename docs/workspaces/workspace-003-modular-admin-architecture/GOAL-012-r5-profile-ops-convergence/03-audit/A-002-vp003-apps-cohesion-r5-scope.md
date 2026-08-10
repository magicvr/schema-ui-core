---
id: A-002-vp003-apps-cohesion-r5-scope
doc: audit-entry
goal: GOAL-012-r5-profile-ops-convergence
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: R5 对 VP-003 代码内聚审计（Root A-010）的继承范围 · C5.1 residual / R5-I001 可见性
audit_type: ad-hoc
verdict: conditional
status: recorded
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# A-002 · R5 scope · VP-003 apps 内聚审计继承（2026-08-05）

- **source**：independent
- **auditor**：Grok Build / grok-4.5
- **类型**：ad-hoc（R5 门禁相关子集；全文见 Root）
- **scope**：GOAL-012 C5.1 / R5-I001 residual 与数据生命周期可见性；继承 Root [A-010](../../GOAL-001-modular-admin-architecture/03-audit/A-010-vp003-apps-cohesion-alignment.md)
- **verdict**：conditional

## 声明

本意见只写 GOAL-012 审计台账，**不**改 status/progress/goal-tree/方案/代码。  
完整证据与 F-001～F-009 正文在 Root **A-010**；本条只固定 **R5 必须响应的子集**，避免推进 C5.1 时只读本目标台账而漏债。

## 结论摘要

Root A-010 verdict **conditional**：`apps/api` 在 R4 表面 Provider 化之后，**store 上帝对象、生产迁移未走 `CollectPersistence`、seed 非贡献驱动、Schema 非 ContributionSet 发布** 仍与 VP-003 退出 #2/#3/#4 冲突。  
GOAL-012 既有 residual（E-002）覆盖 Schema/适配器/PolicyID/双 Profile/Configuration，**未覆盖 store·Persistence 内聚** → 构成 R5 信息门禁缺口（下表 F-R5-IND-001）。

## 本目标开放 Findings

### F-R5-IND-001 · R5 residual / R5-I001 未纳入 store·Persistence 内聚债

- **严重度**：high
- **建议**：required
- **状态**：open
- **关联**：Root A-010 **F-008**（主）；并牵引 F-001 / F-002 / F-005 的阶段登记
- **描述**：在响应 R5-I001 / 勾选 C5.1 前，必须把下列债写入 residual 清单或 required 信息项（实现可后移 R6，但**不得不可见**）：
  1. `internal/store` 领域/迁移/seed 所有权迁出模型（A-010 F-001）
  2. `CollectPersistence` 生产接线与历史 0001–0008 descriptor 归属（A-010 F-002）
  3. seed/RBAC reconcile 以 Authorization（或 system-data）贡献为源（A-010 F-005）
- **证据**：`00-meta.md` C5.1 / R5-I001；`02-execution/E-002-r5-readyz-and-residuals.md`；Root A-010
- **建议修复**：`/govern` 更新 R5-I001 证据列与 residual 表；或 D 记录明确「模型在 R5、迁出在 R6」的最晚阶段与复审触发；未登记前不得将 C5.1 中「R4 residual 闭合」写成已覆盖 Persistence 终态。

### F-R5-IND-002 · Schema 非 ContributionSet 驱动（继承）

- **严重度**：med
- **建议**：required
- **状态**：open
- **关联**：Root A-010 **F-003**；E-002 residual「Schema 完全贡献驱动」
- **描述**：与既有 residual 同债；本条确认其在 R5 C5.1 仍为 **required** 闭合项（fixed 或 accepted-residual + 复审触发），不得在 C5.1 叙事中静默降级为可忽略。
- **证据**：`apps/api/internal/handler/schema.go`；Root A-010 F-003

### F-R5-IND-003 · 中心 Settings/Activity 适配器删除（继承，recommended）

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **关联**：Root A-010 **F-004/F-007**；E-002 pending / R6
- **描述**：保持 E-002 登记；R5 不强制完成删除，但 C5.1/R6 交接清单须继续列出。

## 继承对照（Root A-010 → R5）

| Root finding | R5 动作 |
|--------------|---------|
| F-001 store 上帝对象 | 经 F-R5-IND-001 登记；实现可切片 |
| F-002 CollectPersistence 未接线 | 经 F-R5-IND-001 登记；与 C5.2 数据生命周期叙述交叉时不得宣称已生产化 |
| F-003 Schema 贡献驱动 | F-R5-IND-002 · C5.1 |
| F-004 / F-007 适配器双轨 | F-R5-IND-003 · residual → R6 |
| F-005 seed 贡献驱动 | 经 F-R5-IND-001 登记 |
| F-006 catalog 双源 | recommended；非 C5 阻断 |
| F-008 台账缺口 | **本条 F-R5-IND-001** |
| F-009 Web | 无 R5 必作 |

## 必改项（本目标）

1. **F-R5-IND-001**（required）：扩展 residual / R5-I001  
2. **F-R5-IND-002**（required）：Schema 贡献驱动闭合策略  

## 结论与下一步

**verdict: conditional** — R5 可继续实施，但 **C5.1 / R5-I001 在 F-R5-IND-001/002 开放时不得无条件勾选完成**。  
请用 **`/govern`** 响应本条与 Root A-010（优先台账登记，再排实现）。

### 声明

本意见不修改 status/progress；不关闭 Root finding；不替代 Root A-010 全文。
