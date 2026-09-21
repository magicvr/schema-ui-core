---
id: D-001-r3-implementation-scope
doc: decision-entry
parent: GOAL-004-r3-version-and-diagnostics
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · R3 实施口径

- 门禁：`currentUser.permissions` 含 `monitoring.read`
- 版本权威：`GET /api/system-monitoring/status` 信封 `items[0].version`（与 `pkg/version` 同源）。失败则隐藏 chip，不打断 Shell
- 不展示 `commit` / `BuiltAt` / 模块列表
- 升级入口：`https://github.com/magicvr/schema-ui-core/blob/main/QUICKSTART.md`（R1 默认 GitHub blob）
- 诊断入口：有 `system-monitoring` 页面时链到 `/system-monitoring`；不新建页
- mvp/demo 默认无此权限 → 无 chip
