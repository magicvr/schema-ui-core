---
id: D-005-c2-c3-guardrails-proposed
doc: decision-entry
status: proposed
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-005 · C2/C3 guardrails 草案

本条把当前可执行方案草案集中到一处，供 self + grok independent 审视；不把草案当成用户已接受的最终实现细节。

- codec boundary：Store/Tx 适配 + Repository 使用 `time.Time`；共享 internal temporal codec；禁止 driver type 泄漏。
- conversion：seconds/milliseconds 按列转换；PG `USING` / SQLite table rebuild；fixed-6 UTC TEXT；sentinel/NULL 映射逐列冻结。
- catalog：v1–v72 immutable；模块 owner 从 v73 append-only；测试冻结表 append。
- predicates：jobs/recycle/digital-offer/login-failures/task-runs/nullable fields 的 CHECK/index/WHERE/ORDER 必须清单化并重建/改写。
- backup：D-006 统一 Backup SPI/Service，SQLite/PG native provider，metadata/verification/restore-to-new-db，rollback 优先；完整 backup product 排除。
- wire：D-003/D-005 6 位输出 + RFC3339 兼容输入；formatter/parser/fixture 清单与 VP-020 R3 接口。

最终冻结前仍需处理 A-006 required F-I-002～F-I-006、F-I-010，并接受下一次 grok independent 复审。
