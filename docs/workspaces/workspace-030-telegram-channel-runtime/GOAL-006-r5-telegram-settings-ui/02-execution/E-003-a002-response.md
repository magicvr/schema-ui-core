---
doc_type: goal-execution
id: E-003-a002-response
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-003 · A-002 意见响应（R-001～R-004 fixed）

## 事实

用户指令（2026-09-03）：响应 GOAL-006 A-002（independent pass），R-001～R-004 顺手修，不阻断关门。

1. **R-001（nav 绑定 settings.read）**：`channel.telegram` DependsOn 增加 `admin.settings`（`modules/channel/telegram/provider.go` + `kernel/profile.go` BuiltinModules）；`menu_telegram` `Permission: "settings.read"`（同 ContributionSet 内由 admin.settings 声明）。测试：module integration 改经 registry 解析 plan + `settingsPermissionStub` 声明 settings.read，断言 `set.Navigation[0].Permission == "settings.read"`；composition 相关 plan/ModulesEnabled 对齐（补 `core.operationlog`、`admin.settings`）。
2. **R-002（导航顺序）**：`kernel/provider.go` `DefaultNavigationOrder` 在 `menu_mail_outbox` 后插入 `menu_telegram`；`navigation_order_test.go` 快照同步。
3. **R-003（组合根 schema 探测）**：新增 `TestTelegramSettingsSchema_MountAndDisable`（`composition_telegram_test.go`）：启用模块经 `newMux` → `GET /api/schema/telegram-settings` 200（dev-session auth）；禁用（mvp）→ 404。
4. **R-004（UI 清除密钥）**：`telegram-admin-tab.tsx` 两步「清除已保存密钥」（确认 → PATCH 空串），状态翻转为未配置；`en-US.json`/`zh-CN.json` 新增 `schema.telegram.clear.*` + `feedback.cleared`；`telegram-admin-tab.test.tsx` 新增清除用例。

## 验证

- `go build ./...`、`go vet`、`gofmt` 干净；`go test ./...`（apps/api）全部 ok。
- 前端：`telegram-admin-tab.test.tsx` 3/3、`schema-keys.structural.test.ts` 4/4 ok。

## 评估

R-001～R-004 全部 fixed；GOAL-006 维持 done。详见 [A-003](../03-audit/A-003-a002-response.md)。
