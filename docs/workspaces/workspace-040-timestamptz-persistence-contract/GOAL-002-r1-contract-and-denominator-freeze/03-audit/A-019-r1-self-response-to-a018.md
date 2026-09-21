---
id: A-019-r1-self-response-to-a018
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-018 / matrix expression and E-index correction
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-019 · R1 self response to A-018

## 响应

- `F-I-002`：column matrix PG expressions 已同步为 `date_trunc('microseconds', to_timestamp(seconds))` 与 `date_trunc('microseconds', epoch + milliseconds * INTERVAL '1 millisecond')`；不再保留 double/1000 或 typmod-only candidate。guardrails 与 column-contract 已同形；逐 owner USING/rebuild/codec/test 仍开放。
- `F-I-016`：child execution index 已修正为 E-001～E-019 单调，E-019 owner allocation 位于 E-018 predicate matrix 之后；关闭该 recommended。
- `F-I-003`～`F-I-006` 继续开放；D-012/D-013 方向已唯一但逐列/逐 callsite/restore/checksum evidence 未完成。

## 放行

C2/C3 仍不可冻结，R2 不启动；下一次 independent 复审应以当前 matrix/执行索引为准。
