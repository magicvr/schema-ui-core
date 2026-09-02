---
doc_type: goal-decision
id: D-001-r4-evidence-closure-plan
parent: GOAL-004-evidence-closure-and-closeout
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-001 · R4 证据闭环、边界核账与关门方案

## 触发

推进 Root `GOAL-001` 纲领阶段 R4，对 VP-029 的退出判据、红线与实证进行终审核验。

## 决定

1. **证据矩阵**：
   - 逐条对照 VP-029 方向级判据 #1～#7 进行证据映射与实证复核；
   - 确认各子目标（GOAL-002, GOAL-003）的审计意见全部闭合（open required = 0）。
2. **越界核账**：
   - 核验 `git diff origin/dev`，确认未改动 Charter、未改动 Profile 默认集、未引入支付网关或 Telegram 依赖、未重开 VP-011 账本。
3. **审计闭环路径**：
   - 执行 GOAL-004 自审（A-001）；
   - 调用本地 grok build（grok-4.6 · high）执行独立关门审计（A-002）；
   - 响应并在 open required = 0 确认后，执行 GOAL-004 与 Root GOAL-001 关门。
