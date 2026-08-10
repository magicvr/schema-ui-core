---
id: D-002-r4-c4-criteria-scope
doc: decision-entry
goal: GOAL-010-r4-c4-schema-other-migration
source: user
date: 2026-08-05
status: accepted
---

# D-002 · C4.4 成功标准收窄与 ledger 门禁移交

## 决策

C4.4 成功标准收窄为 C4 可交付的 Manifest secrecy 扫描、Ready 失败反向清理、
PolicyID/Visibility/JSON 校验器与 Records historical-only 保持。**ledger
drift/unknown 运行时 fail-closed** 属 store/migration 数据路径改造，移交 C5 数据
门禁（fresh/upgrade/reconcile 深度验证），不在 C4 声明完成。

## 依据

- 冻结包 §4.2 的 ledger drift/unknown 需 store/migration 层运行时证据，与 C4 的
  settings/activity provider 迁移不同域。
- Grok A-008（C4 审计）F-IND-C4-002 指出：C4.4 字面含 ledger 项则不得无条件勾选。
- E-002 已诚实登记该 residual；本决策将 meta 条文与实施对齐。

## 影响

- C4.4 按收窄条文验收（secrecy/Ready 清理/校验器/Records）。
- C5 数据门禁新增：ledger drift/unknown 运行时 fail-closed（mvp/admin 双 Profile）。
- 其余 C4 范围不变。
