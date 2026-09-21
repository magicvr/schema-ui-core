---
id: A-017-r1-self-response-to-a016
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-016 / precision-voucher-monotonic-owner evidence
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-017 · R1 self response to A-016

## 已响应

- F-I-002：冻结包的 PG legacy expressions 已统一为 `date_trunc('microseconds', to_timestamp(...))`；milliseconds 改为整数 `INTERVAL '1 millisecond'`，不再使用 double/1000；新 Go writes 使用 truncate-to-microsecond。逐 owner DDL/codec 仍开放。
- F-I-003：D-012/D-013 已唯一冻结 voucher `0→NULL/negative fail closed` 与 users/roles `max(truncatedNow, old+1µs)`；matrix 与 predicate draft 已同步。90 列 old→new→r/w 仍需逐 owner 证据。
- F-I-004：CreateRecoveryPoint draft、PG/SQLite provider、restore-to-new-db 后置条件已落盘；Port/package/tool/rollback 仍开放。
- F-I-005：v73 owner allocation、`core.persistence` schema ledger owner、72-entry test/timeNames 清单已落盘；allocation 与 test rewrite 仍 proposed。
- F-I-006：predicate matrix 已补 monotonic D-013 与 voucher D-012；exact SQL/rebuild order/callsite tests 仍开放。
- F-I-015：E-ID collision 已 fixed，E-001～E-019 当前单调。

## 放行

A-016 的五条 required 仍开放；C2/C3 不冻结，R2 不启动。下一步是逐 owner 写出 accepted SQL/codec/Port/test evidence，再进行 grok independent 复审。
