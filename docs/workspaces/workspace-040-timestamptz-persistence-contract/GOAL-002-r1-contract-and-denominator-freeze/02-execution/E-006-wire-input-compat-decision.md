---
id: E-006-wire-input-compat-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-006 · R1 wire 输入兼容事实

承接 Root D-005 / Root E-006：公共输出固定为 6 位微秒 RFC3339 UTC `Z`；入站与 fixture parser 兼容合法 RFC3339 0/3/6/9 位小数及 `+00:00` 等等价 offset，解析后统一 UTC。此条只记录用户裁决，formatter/parser/fixtures 尚未实施。
