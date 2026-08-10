---
id: E-004-r3-c2-module-projection
doc: execution-entry
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: recorded
---

# E-004 · R3 C2 模块投影实施

## 已完成事实

- `kernel.Plan.HasModule` 成为运行时注册和投影的单一选择查询。
- API 组合根把同一 Plan 传入 handler；Settings/Activity 由各自模块包注册，
  Settings 依赖 always-on `core.operationlog`。
- Settings/Activity Schema JSON 移入模块自有 `schema` 包并通过 `go:embed` 提供；
  API Schema handler 按 Plan 过滤，旧 handler fixture 中的两个文件已移除。
- Manifest 响应标记 API 来源；Settings 成功 PATCH 返回
  `X-Schema-UI-Config-Changed: settings.branding`。
- Web 主入口通过通用 config-aware fetcher 发布 Host 事件，App/LoginPage 订阅
  branding namespace；开发静态 Manifest 命中会输出明确 warning。

## 边界

本条只记录 R3 有界切片；Users/Roles 全量模块迁移、中心 Schema 所有权的完整
清理和 R4 其他能力仍属于后续阶段。
