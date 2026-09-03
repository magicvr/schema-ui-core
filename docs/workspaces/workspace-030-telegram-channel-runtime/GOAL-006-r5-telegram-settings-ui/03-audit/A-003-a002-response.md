---
doc_type: goal-audit
id: A-003-a002-response
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
source: self
scope: GOAL-006 A-002 独立复审意见响应（R-001～R-004 顺手修）
audit_type: finding-closure
verdict: pass
open_required: 0
---

# A-003 · GOAL-006 A-002 独立复审意见响应（R-001～R-004）

## 1. 响应背景

用户指令（2026-09-03）：响应 GOAL-006 [A-002](A-002-independent-ui-tab-audit.md)（independent `pass` · required=0），R-001～R-004 recommended **顺手修一下，不阻断关门**。

## 2. Recommended 处置台账

| ID | 级别 | 处置 | 证据 |
|----|------|------|------|
| **R-001**（导航未绑 `settings.read`） | recommended | **fixed** | `channel.telegram` DependsOn 增加 `admin.settings`（`provider.go` + `kernel/profile.go` BuiltinModules），`menu_telegram` 的 `Permission: "settings.read"` 由同 ContributionSet 内 admin.settings 声明（权限键全局唯一、无新键红线保持）。测试：module integration（registry 解析 plan + settingsPermissionStub 声明 settings.read 断言 nav.Permission）+ composition 测试全部对齐。 |
| **R-002**（`menu_telegram` 不在 DefaultNavigationOrder） | recommended | **fixed** | `kernel/provider.go` DefaultNavigationOrder 在 `menu_mail_outbox` 后插入 `menu_telegram`；`navigation_order_test.go` 快照同步。 |
| **R-003**（无组合根 schema 200 断言） | recommended | **fixed** | 新增 `TestTelegramSettingsSchema_MountAndDisable`（composition_telegram_test.go）：经 `newMux` 组合根，启用模块 `GET /api/schema/telegram-settings` → 200（dev-session auth），禁用（mvp）→ 404。 |
| **R-004**（UI 不能清空已保存密钥） | recommended | **fixed** | `telegram-admin-tab.tsx` 新增两步「清除已保存密钥」动作（确认 → PATCH `{"bot_token":"","webhook_secret":""}`），清空后 status 翻转为未配置；i18n 双语 keys；`telegram-admin-tab.test.tsx` 新增清除用例（3/3）。 |

## 3. 验证证据

- `go build ./...` 通过；`go vet` 干净；`gofmt` 干净。
- `go test ./...`（apps/api）全部 ok。
- 前端：`telegram-admin-tab.test.tsx` 3/3、`schema-keys.structural.test.ts` 4/4 ok；改动文件 `tsc` 无错误（`form-controls.tsx` 类型错误为既有问题，非本次引入）。

## 4. 结论

A-002 independent `pass` 接受；R-001～R-004 全部 **fixed** 闭合。开放 required = 0；无新 recommended。GOAL-006 维持 `done`。
