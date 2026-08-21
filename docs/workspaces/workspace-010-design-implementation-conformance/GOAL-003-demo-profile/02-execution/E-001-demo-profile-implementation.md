---
id: E-001-demo-profile-implementation
doc: execution-entry
goal: GOAL-003-demo-profile
status: recorded
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# E-001 · W2 demo Profile 实施 + 回归 + go 留痕

## 事实（2026-08-11）

按已确认方案（`demo` / API + e2e demo）完成实施：

### kernel（`apps/api/internal/kernel/profile.go`）
- 新增 `ProfileDemo ProfileName = "demo"` 常量。
- `profileDefaults[ProfileDemo]` = **mvp 模块集 + `dev.examples`**（core 六项 + admin.users + admin.roles + dev.examples）。`ResolveProfile("demo", nil)` 成功，`Source = profile.default`。

### 测试分母（D-003 §6 口径扩展）
- `kernel_test.go`：`TestBuiltinProfilesResolveDeterministically` 增加 demo 解析断言；新增 `TestDemoProfileIsNonProduction`（mvp/admin 默认**不含** dev.examples；demo Source=profile.default）——S3 卫生。
- `composition_test.go`：新增 `TestDemoProfileManifest`——demo plan 含 dev.examples + admin.users；manifest `homePageRef=overview`、含 users/roles/overview/data-table/form-controls、无 settings/activity；`/api/schema/overview` 200；mvp/admin 仍无 dev.examples——S2/S3。
- config_test 无回归（mvp/admin 分母不变）。

### e2e（playwright）
- `playwright.config.ts`：`APP_PROFILE` 白名单放开 `demo`。
- `shell.spec.ts`：按 `isDemoProfile` 分支——demo 下 manifest 含 examples、`homePageRef=overview`、home→`/overview`、Overview/Data table 导航可见；mvp/admin 保持无 examples（W1 S5 卫生）。
- `schema-crud.spec.ts`：`signInAsAdmin` 按 profile 断言 home（demo→`/overview`，否则 `/users`）。
- `localization.spec.ts` 两测在 demo 下 skip（admin/mvp 边界证据，非 demo 目标）。

### 文档（S5）
- `apps/api/README.md`：`APP_PROFILE` 说明 + `demo`（非生产向演示 Profile = mvp 集 + dev.examples，home=overview）+ Profile 选择段。
- `README.md`（根）：启动示例、compose 示例、Profile 说明补 demo。
- `apps/web/README.md`：e2e 运行示例补 `APP_PROFILE=demo`。

### 回归证据
- `go test ./...`（apps/api）：23 包全绿。
- `npm test`（apps/web）：44 文件 / 746 测试通过。
- Playwright e2e：**mvp** 3 passed/1 skip、**admin** 3 passed/1 skip、**demo** 2 passed/2 skip（shell + schema-crud 在 demo 下通过，展示范例面 + home=overview + users CRUD）。

## VP-008 `go` 消费有效性 —— 判定（S6）

本波为**新增编译 Profile `demo`**（Profile 默认集变更，模块矩阵新增一个非生产候选）。判定：

- **mvp/admin 生产向默认集未变**（`TestDemoProfileIsNonProduction` 证据）；`demo` 为**非生产向**演示 Profile，不进入业务 VP 候选身份。
- 因此 **VP-008 `go` 消费保持有效，不触发暂挂**；W1 恢复的 `go` 证据（`4a2b8cd…`）对 mvp/admin 生产矩阵仍适用。
- 留痕：若未来业务 VP 以 `demo` 为候选，则按 VP-008 触发消费前 freshness review。

## 尚未发生

- W2 波次审计（roadmap 阶段 5：self + independent，grok-build；触及 Profile 矩阵）。
