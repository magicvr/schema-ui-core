---
id: D-001
doc: decision-entry
goal: GOAL-006-channel-provider-contract
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · 开设 R5 子目标

## 背景

Root D-006 已否决组合层关门并扩展纲领 R5～R8。P-001：按阶段创建子目标。本回合只开 R5。

## 决定

1. 创建 `GOAL-006-channel-provider-contract`，parent = `GOAL-001-outbound-mail`，status `active`。
2. 本目标工作 = 冻结渠道合同并关闭 Root I-011；**本条不**冻结 I-011。
3. 审计模式：scaffold **none**；合同冻结落盘后 **self**。
4. 本回合不改 `apps/api` / `apps/web`，不创建 R6～R8 子目标。

## 未选方案

- 一开子目标就写 Resend 客户端：I-010/I-011 未关闭。
- 同时创建 R6～R8：违反 P-001。
