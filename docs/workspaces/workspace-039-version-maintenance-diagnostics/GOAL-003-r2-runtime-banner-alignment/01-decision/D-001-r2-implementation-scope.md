---
id: D-001-r2-implementation-scope
doc: decision-entry
parent: GOAL-003-r2-runtime-banner-alignment
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · R2 实施范围

无新 P-004。按 GOAL-002 D-001：

- 生产者 `maintenance`/`degraded`/`read-only` → Host `degraded`
- `/me.runtimeMode` 原样字符串；`fetchMe` 必须投影
- 横幅只读该字段；登录页不渲染 Shell 故无横幅
- 空 runtime 配置视为 `normal`
