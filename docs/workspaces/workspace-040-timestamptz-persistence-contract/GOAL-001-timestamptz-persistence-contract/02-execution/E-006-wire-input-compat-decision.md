---
id: E-006-wire-input-compat-decision
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-006 · wire 输入兼容用户裁决

用户选择公共时间输出严格为 6 位微秒 RFC3339 UTC `Z`，但入站/fixture/parser 兼容合法 RFC3339 的 0/3/6/9 位小数与 `+00:00` 等等价 offset；解析后统一 UTC，写入由 storage codec 规范化。C2 必须把 formatter/parser/fixtures 影响列入方案与测试矩阵。

证据：Root `01-decision/D-005-r1-wire-input-compat.md`。
