---
id: D-016-open-typecheck-guard-hardening
doc: decision
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-016 · 开设 GOAL-010 加固类型检查守卫并约定交叉审计与关门条件

## 决定

2026-09-18，用户在获知 `GOAL-008 A-002 F-001`（守卫按 `-p` 令牌判定检查型调用，`tsc --noEmit -p tsconfig.json` 空转仍会通过）后指示：

> 开一个小整改子目标（GOAL-010）加固守卫并补 `-p tsconfig.json` 变异用例。然后交叉审计（独立审计使用本地 grok build cli，模型 grok 4.6，思考强度 xhigh）确认没有问题之后，关门。

据此：

1. 开设 `GOAL-010-typecheck-guard-hardening`（Root 下的**非纲领整改子目标**，不改变 Root 六阶段分母与 `progress: 5/6`），承接并闭合 `GOAL-008 A-002 F-001`。
2. 范围限定为**守卫加固 + 合成变异用例**，不改正确口径、不改产品代码/脚本/CI（`D-001`）。
3. 审计模式为 **`independent`**（承接安全相关证据完整性问题，且用户明确要求交叉审计）：先 self 审计，再按项目级路径 `docs/architecture/independent-audit-execution.md` 调用本地 grok build；**本次思考强度按用户指令提高到 `xhigh`**（项目级决策默认为 `high`）。
4. 关门条件：self + independent 意见均已落盘并合并响应，且不存在未合法闭合的 required/必改 finding。

## 未选方案

- **不修守卫，仅保留 recommended**：用户明确要求修复，不采纳。
- **顺带修复其他 recommended**（`GOAL-008 A-001 F-002` CI 步骤正向断言、`GOAL-009 A-001 F-001/F-002` 暗色断言与页面覆盖）：会放大范围与审计成本，不采纳；保持各自台账 open。
- **把守卫加固并入 `GOAL-008`（回开该目标）**：`GOAL-008` 已以 `done · 4/4` 关门，回开需要改写其关门结论；按 `GOAL-009` 先例以平行整改子目标承接更清晰，不采纳回开。
- **只加 `-p tsconfig.json` 一条黑名单**：`D-001` 已否决（漏 `-p .`、`--project=` 与未来同类配置）。

## 边界与不变量

- 不改变 Root 六阶段分母、`progress: 5/6` 与 Root/VP 的 `active` 状态。
- 不关闭 `R5-I-004`（用户书面 Root/VP 关门确认）——本目标关门**不**等于 Root/VP 关门。
- 独立审计 provider 与强度按用户指令记录为 `grok build (grok 4.6 · reasoning xhigh)`；若该 provider 不可用或无产出，independent 门禁保持未满足，不得由 self 冒充。
