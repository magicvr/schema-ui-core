---
doc_type: goal-decision
id: D-001-r3-goal-establishment
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
status: done
version: 0.2.0
---

# D-001 · R3 目标建立与阶段路线

## 已发生决策

- 在 R2 `GOAL-003-r2-connection-settings` 已由 A-018 Grok independent `pass`、A-019 response 关闭后，建立 R3 子目标 `GOAL-004-r3-session-operator-console`，父目标为 `GOAL-001-telegram-operator-console`。
- R3 采用四个串行检查点：C1 合同/信息/用户裁决；C2 入站文本与会话落盘；C3 API/权限/运行时接线；C4 UI、发言权反馈与独立审计。
- 当前只冻结 VP-033 已有边界与门禁，不冻结会话主键、发言权、权限或发送状态方案；这些列为 C1 的信息需求，未获用户裁决前不得实施依赖它们的代码。
- 现阶段不创建更细粒度子目标；待 C1 方案、写集和并行价值明确后，再判断 C2/C3 是否需要拆分子目标承载治理上下文。

## 不变边界

R3 只处理 Telegram 实际投递的文本、未绑定人工台和现有 sender；不引入历史回灌、FSM、群发、频道、多 bot、多实例 polling、独立进程或 SSE/WebSocket。

## 待用户裁决

待处理 `I-033-009`、`I-033-010`、`I-033-019`、`I-033-021`、`I-033-022`；`I-033-020` 的 Update 幂等边界需在同一 C1 合同中明确。

## 候选方案材料（非决策）

`attachments/r3-c1-option-analysis.md` 已记录 I-033-009/010/019～022 的证据、互斥候选、AI 推荐和主要取舍。该附件只承载 C1 信息收集，不代表用户已选择任何方案；方案依赖的代码、迁移和权限写入继续保持门禁。
