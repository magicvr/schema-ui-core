---
id: GOAL-004-s2-module-contract-access-drill
doc: execution
status: active
parent: GOAL-001-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# 执行记录 · GOAL-004

## 执行条目索引

| E-ID | 日期 | 动作 | 结果 | 证据 / 文件 |
|------|------|------|------|-------------|
| E-001 | 2026-08-10 | 非领域化接入演练（probe 模块） | pass | API 侧 `apps/api/internal/composition/s2_access_drill_test.go`；Web 侧 `apps/web/src/app/s2-access-drill.render.test.tsx` |
| E-002 | 2026-08-10 | M1–M6 逐模块核验闭环 + 依赖/Profile/迁移正反测试复核 | pass | S1 模块检查表 + kernel/composition 既有正反测试 |
| E-003 | 2026-08-10 | I-READINESS-002 闭合 | pass | 本子目标证据汇总 |

## 非领域化接入演练（S2-3）

**API 侧（Go composition）**：test-only probe 模块 `admin.probe`（generic 资源 `probe-items`，映射框架共性能力 list/detail/create/update/delete）实现 `kernel.Provider` 六项贡献。`newMuxWithExtraProviders` 装配 seam（production `newMux` 委托、extra=nil）验证：

- probe route `GET /api/probe-items` 可访问（200）
- probe schema `GET /api/schema/probe-items` 发布（200）
- probe 页面进入聚合 Manifest
- probe 不进入 mvp/admin 默认 Profile（`TestS2AccessDrill_ProbeAbsentFromDefaultProfiles`）
- 无 handler 中央业务注册、无 Web Renderer/Shell 改动

**Web 侧（Renderer/Shell）**：内存 Manifest 注入 probe 页面 `probe-items`（宿主从未硬编码该 pageId）+ probe schema 文档，App 经通用 schema-driven 路径在 `/probe-items` 渲染（h1 + body + 侧边导航），证明新增模块页面无需 Renderer/Shell 中央业务注册。

**probe 生命周期**：`admin.probe` 仅存在于 `s2_access_drill_test.go` 测试文件（test-only 候选集），不进入编译默认候选集、不进入 `mvp`/`admin` 默认启用集、不进入生产 Manifest；无独立 migration（归属 `core.auth-session`，与 `admin.users` 同 C3.3 模式）。S5 前维持 test-only 资产，或按 VP-008 清理。

**M1–M6 闭合**：S1 已核验 4 standard-admin 模块 M1–M6 全 pass；本演练证明新增模块按同一契约接入。依赖缺失（custom 无 modules → `PROFILE_MODULES_REQUIRED`）、Descriptor mismatch（→ `MODULE_API_MISMATCH`）、贡献冲突（→ `MODULE_CONTRIBUTION_CONFLICT`）、Profile 未知（→ `PROFILE_UNKNOWN`）正反测试由 kernel/composition 既有测试覆盖（`TestResolvePlanUsesConfiguredProfileAndRejectsMissingDependencies`、`TestDualProfileRegisterValidationFailClosed`、`TestRegisterContributionsCrossModuleConflict` 等）。

## 基线影响

composition.go 增加 `newMuxWithExtraProviders` seam（production `newMux` 行为不变，仅委托）。候选基线从 `852ee7e` 前进至本提交；`go build ./...`、`go test ./...`（全包）、`npm test`（41/729）实测通过，无回归。

## 记录规则

只写已发生事实；正反测试、probe 模块、命令结果绑定本提交。probe 生命周期与 test-only 归属已记录。
