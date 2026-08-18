---
id: GOAL-028-w17-cron-preview-field-binding
title: 决策记录 · W17 · Cron 字段绑定与中文 describeCron
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# 决策记录 · GOAL-028 · W17

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 预览如何接到表单字段而不新增协议控件类型 | S1 / S2 | S1 | 对照 FormControls 白名单与 custom-component 注册表 | **verified** | — | D-001 |

## 决策索引

| 编号 | 标题 | 日期 | 状态 | 摘要 |
|------|------|------|------|------|
| [D-001](01-decision/D-001-w17-freeze.md) | W17 方案冻结：字段绑定 + 本地化 describeCron | 2026-08-18 | accepted | `afterComponent` 挂到 create/edit `cron` 字段；`describeCron` 按 Accept-Language 出中文/英文人话 |
