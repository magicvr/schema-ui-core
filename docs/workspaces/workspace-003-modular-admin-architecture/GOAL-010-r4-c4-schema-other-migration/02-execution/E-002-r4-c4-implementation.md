---
id: E-002-r4-c4-implementation
doc: execution-entry
goal: GOAL-010-r4-c4-schema-other-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R4-C4 settings/activity 迁移与门禁实施

## 已发生事实

- **C4.1 settings provider 化**：`modules/settings/provider.go`（HTTP 4 路由含公开
  branding + 鉴权 settings list/detail/patch、Schema `settings`、权限
  settings.read/write、导航 menu_settings、Manifest fragment）；`handler.SettingsRoutes`
  构建器；completion 消费 provider。
- **C4.2 activity provider 化**：`modules/activity/provider.go`（只读 operations
  HTTP + Schema `activity` + 权限 operations.read + 导航 menu_activity + Manifest
  fragment）；`handler.OperationsResource` 导出；POST → 405（只读）。
- **Manifest 全 fragment 化**：基线 app-manifest.json 移除 settings/activity 页与
  导航（核心专用）；`ForModulesWithFragments` 简化为 core + provider fragments；
  `adminModules` 完全删除。
- **C4.3 Schema owner map 贡献驱动**：`modules/{users,roles,settings,activity}/schema`
  暴露 ModuleID/PageIDs；`schemaDocumentsForPlan` 由模块 contributor 派生 owner，
  plan 门禁保留（解决 F-IND-C33-001）。
- **C4.4 门禁**：Manifest secrecy（`rejectFragmentSecrets`，含 secret 键 fragment
  fail closed）；PolicyID/Visibility/JSON 校验器（kernel）；Ready 失败反向 Stop 由
  composition `registerLifecycle` 处理（已有）；ledger drift/unknown 运行时
  fail-closed 登记 C5 residual（store/migration 路径，需迁移层改造）。

## 验证

- API `go test ./...`（含新增 settings/activity provider 测试、manifest secrecy 测试）
  通过；`go vet ./...` 干净；Web `vitest run`（495）通过。
- composition 消费 4 个 admin provider（users/roles/settings/activity）；中心
  RegisterSettings/RegisterActivity 分支不再用于生产。

## 边界

- ledger drift/unknown 运行时 fail-closed、运行时双 Profile 深度矩阵登记 C5。
- handler 测试环境仍走中心 RegisterSettings/RegisterActivity（与生产 provider
  路径同工厂等价，C33-002 已文档化）。
