---
id: GOAL-005-requestid-correlation
doc: decision
status: active
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| （继承）I-005 | required | request-id / correlation 如何写入 span（属性名、是否 baggage） | R4 方案冻结 | R4 关联前 | 本目标 D-001 | verified（继承） | — | `01-decision/D-001-requestid-correlation.md` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-22 | request-id ↔ span 关联合同（闭合 I-005） | accepted | `01-decision/D-001-requestid-correlation.md` |