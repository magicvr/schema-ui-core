---
id: E-002-r3-implementation
doc: execution-entry
parent: GOAL-004-r3-version-and-diagnostics
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · T-4 实施事实

- `version-chip.tsx`：`monitoring.read` 门禁；`GET /api/system-monitoring/status` 取 `items[0].version`；失败隐藏
- 升级链接：`QUICKSTART_UPGRADE_URL` GitHub blob `main/QUICKSTART.md`
- 诊断：manifest 含 `system-monitoring` 时 `onNavigate("/system-monitoring")`
- 不展示 commit
- 测试：无权限不 fetch；展示 version 不含 commit；403 隐藏
- vitest version-chip + catalog + App.integration + auth-gate + ui-bilingual 全绿
