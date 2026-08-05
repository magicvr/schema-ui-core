---
id: E-005-r4-c3-cutover-schema-manifest
doc: execution-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-005 · R4-C3 Schema 与 Manifest 特例迁移（C3.3 续作）

## 已发生事实

- **Schema 内容迁入模块**：`modules/{users,roles}/schema/` 承载 users/roles 页面文档
  （原 `handler/fixtures/schema/{users,roles}.json` 已删除）；`handler/schema.go`
  `staticSchemaDocuments` 聚合模块 schema 包。
- **Manifest 内容迁入模块**：`modules/{users,roles}/manifest/` 承载模块 manifest
  fragment（协议版本/能力/app 身份/页面/导航）；provider 的 FragmentContribution.JSON
  读取模块 fragment；嵌入基线 `app-manifest.json` 移除 users/roles 页面与 admin 导航。
- **manifest 聚合**：新增 `ForModulesWithFragments`，基线投影 settings/activity +
  合并 provider fragment；composition 传 `set.Fragments`；`adminModules` 收敛为
  settings/activity。
- **Schema owner map 残余**：`schemaDocumentsForPlan` 的 owner map 仍按 plan 门禁
  users/roles/settings/activity 文档暴露；内容已模块所有，map 作为 plan 投影辅助，
  与 settings/activity（C4）共享，登记为 C3.3 文档化残余。

## 验证

- API `go test ./...`（12 包）与 `go vet` 通过；Web `vitest run`（24 文件 495 测试）
  通过（fixture 目录映射更新为模块 schema 包）。
- 生产 mux：users/roles 路由经 `RegisterContributions`；schema 内容经模块包聚合；
  manifest 经 provider fragment 合并；无中心 users/roles 注册/内容/投影。
