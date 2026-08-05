---
id: E-004-r6-c61-dead-adapter-removal
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-004 · R6 C6.1 死适配器与双轨删除

## 已发生事实

- **测试环境改走真路径**：`testhelpers_test.go` 不再用 `MountProviderRoutes`/
  `RegisterSettings`/`RegisterActivity`；改为内联挂载各模块 Provider 使用的同一资源
  工厂路由（`SettingsRoutes`/`ResourceRoutes`/`resourceRoutes`），与生产
  `RegisterContributions` 行为一致（同 handler）。RegisterContributions 契约由
  kernel/composition 测试覆盖。
- **删除死适配器**：`handler/health.go` `MountProviderRoutes`、`handler/settings.go`
  `RegisterSettings`+`settingsHandler`、`handler/operations.go` `RegisterActivity`+
  `registerOperations` 全部删除；`operations.go` 未用 imports 清理。
- 生产仍仅 provider finalize 一条装配路径；handler 仅 core auth/accounts/health。

## 验证

API `go test ./...`（14 包）+ `go vet` + Web `vitest run`（495）通过。

## C6.1 状态

测试真路径 + handler 级死适配器删除完成（C6.1 主体）；`schemasHandler`/`registerResource`
仍用于测试路径（resources_test 等），业务 surface 已全走 provider 工厂。
