---
id: A-001
goal: GOAL-006-w5-recordview-declared-fields
source: self
date: 2026-08-14
scope: W5 波次关门（recordView declared-fields 契约 + dev/文档卫生 + HEAD 回归）
verdict: conditional
parent: GOAL-006-w5-recordview-declared-fields
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-001 · W5 关门自审（self）

## 结论

**conditional**：本波自身范围（recordView declared-fields + fail-open + dev/文档卫生）完成且证据充分；发现一条**波次外的跨门禁 finding（F-1）**——容器 smoke 复现性破损（W3 引入、post-go），影响 VP-008 `go` 消费重验证，不阻断本目标关门，但阻断 `go` 的干净消费，移交 freshness review 决策。

## 成果

- 波次代码全部入库并绿：V-001～V-006 于 HEAD `c8ae108` 实测通过（E-004）；
- recordView 声明字段契约 + fail-open 兜底语义已落盘（E-001/E-002，I-001 verified）；
- dev/文档卫生归档（E-003）。

## Findings

| finding | 级别 | 主张 | 处置 |
|---------|------|------|------|
| F-1 | major（跨门禁） | 容器 smoke 复现性破损（F-1a claim git 强制 + F-1b nginx upstream 作用域 + F-1c SM-007 陈旧断言，均为 post-go 引入） | ✅ **fixed（2026-08-14 · GOAL-007 W6）**：E-002 实证 V-007 exit 8 + V-008 exit 0 完整绿 |

## 偏差

无（本波范围完整；F-1 属波次外既有引入，W3 关门漏检）。
