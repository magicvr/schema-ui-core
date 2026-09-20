---
id: E-004-public-wire-decision-recorded
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-004 · 公共时间输出合同用户裁决已落盘

用户选择统一改为 6 位微秒 RFC3339 UTC `Z` 输出。该裁决扩大 C2/C3/R3 的证据范围：需盘点现有固定 3 位毫秒/RFC3339 formatter、更新协议/测试夹具，并在 R3 验证微秒回读与 VP-020 展示时区往返。

证据：Root `01-decision/D-003-r1-public-wire-contract-user-decision.md`。
