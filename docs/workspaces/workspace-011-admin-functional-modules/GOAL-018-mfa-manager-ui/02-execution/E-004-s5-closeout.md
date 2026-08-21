---
id: E-004
goal: GOAL-018-mfa-manager-ui
title: S5 关门完成
date: 2026-08-15
status: recorded
parent: GOAL-018-mfa-manager-ui
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-004 · S5 关门（2026-08-15）

## 事实

- A-003（grok build · grok-4.6 · high · independent）S5 关门审计 **fail**：F-001（high，account.json custom 节点过不了运行时 D-VAL）、F-002（med，disable/rotate 未兑现 code|recoveryCode）。
- 响应（全 fixed）：F-001 → 运行时校验本地 overlay（runtime-schema-validate.ts 增加 component 属性；**不改动** I-PROTO-004 钉住的 node.schema.json / component-registry.json）+ load-page.test D-VAL 契约测试；F-002 → splitMFAInput（6 位→code，其余→recoveryCode）+ 契约测试（恢复码停用 + 6 位码轮换 + 请求体断言）；另修生产构建 tsc 类型错误 3 处、01-decision.md I-001 状态。
- A-004（grok reaudit）**pass**：A-003 required 合法闭合；recommended（真实 account.json D-VAL、rotate 断言）已随契约测试与无条件断言补强。
- 审计链 A-001（self 立项+方案）→ A-002（self S2-S4）→ A-003（independent S5 fail）→ A-004（reaudit pass）全部落盘。
- 波次级验证（2026-08-15）：e2e 双 profile 16/16、隔离 compose 容器冒烟 SM-001~007 **PASS**（exit 0）、go 全量 + web 976/976 + tsc 全绿。
- status: active → done；progress 3/5 → 5/5；goal-tree 同步。用户裁决依赖链兑现：GOAL-018 关门 → GOAL-017 回归关门。
