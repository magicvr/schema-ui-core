---
id: D-003-wire-input-compat
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-003 · C2 wire 输入兼容承接

承接 Root D-005：公共输出严格 6 位微秒 RFC3339 UTC `Z`；入站/fixture/parser 兼容合法 RFC3339 0/3/6/9 位小数与 `+00:00` 等等价 offset，解析后统一 UTC。C2 必须列出受影响 formatter、parser、fixture、Web/Go 测试与非 DB 文件系统 ModTime 输出边界。
