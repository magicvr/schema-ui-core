---
id: A-021-r1-self-response-to-a020
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-020 / execution index and expression sub-item
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-021 · R1 self response to A-020

## 响应

- **F-I-016 → fixed/closed**：当前 child `02-execution.md` 索引已核对为 E-001～E-019 严格递增；磁盘文件与 ID 一致。A-020 读取到的是修正前状态，现行索引不再有 E-019 插入 E-017/E-018 之前的问题。
- **F-I-002 expression 子项 → fixed（不关闭总 finding）**：guardrails、column-contract、matrix 当前统一为 `date_trunc('microseconds', to_timestamp(seconds))` 与整数 `INTERVAL '1 millisecond'`，无 `/1000.0`；逐列 USING/rebuild/codec/test 仍开放。

## 仍开放 required

F-I-002、F-I-003、F-I-004、F-I-005、F-I-006。C2/C3 与 R2 仍阻断。
