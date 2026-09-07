---
doc_type: goal-execution
id: E-011-operator-config-and-test-admin
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-011 · Operator 配置启用 channel.telegram + TEST_ADMIN 测试账户机制

## 背景

用户反馈：界面看不到 Telegram bot 配置页。诊断结论：`channel.telegram` 为 VP-030 红线**不进任何默认集**，需显式 `app.modules` 启用；且用户本地 postgres 的 `admin` 密码已改，bootstrap 只对零用户库生效，无法用 `ADMIN_INITIAL_PASSWORD` 验证。用户裁决：约定测试专用账户/密码放 `apps/api/configs/.env`（gitignored）。

## 事实

1. **Operator 配置**（`apps/api/configs/config.yaml`）：`profile: mvp` → `profile: custom` + `app.modules.list`（admin 全量 23 模块 + `channel.telegram`）。
2. **测试账户机制**（可选，零侵入）：
   - `config.go`：新增 `TestAdminUsername`/`TestAdminPassword`（env `TEST_ADMIN_USERNAME`/`TEST_ADMIN_PASSWORD`；仅 PASSWORD 非空时启用）。
   - `systemdata/bootstrap.go`：新增 `EnsureTestAdmin`——创建（不存在）或重置密码 + 清 `must_change_password`（存在），授予 admin/editor 角色，不动既有 `admin` 用户；空密码 = no-op。
   - `composition.go` `openStore`：bootstrap 后按需 upsert 测试账户。
   - `.env.example` 登记两个变量。
3. **dev.cmd 默认直读 config.yaml**（`dev.cmd`）：默认 profile `admin` → `config`——`CONFIG_FILE` 指向 `apps/api/configs/config.yaml`（完整 operator 配置：profile/modules + db/mail/telegram/upload），仅端口经 `HTTP_ADDR` env 覆盖；显式 `--profile/--modules` 仍走 overlay 覆盖。避免旧行为（overlay 只含 app+http，丢弃 config.yaml 其它段）。
4. **测试**：`reconcile_test.go` 新增 `TestEnsureTestAdmin`（创建/重置/不动 admin/空密码 no-op）；`config` 新增 `TestCustomModulesResolveWithTelegram`（内联模块列表解析含 channel.telegram + admin.settings 依赖，不依赖 operator 文件）。

## 验证（端到端）

- `go test ./...`（apps/api）全部 ok；`gofmt`/`go vet` 干净。
- 用户 `.env` 配置 `TEST_ADMIN_USERNAME=testadmin` + `TEST_ADMIN_PASSWORD` 后重启 API（postgres）：
  - `readyz` 200；
  - `testadmin` 登录成功；
  - `/api/accounts/me` features：`menu_telegram = True`（与 `menu_mail`/`menu_mail_outbox`/`menu_settings` 并列）；
  - `GET /api/schema/telegram-settings` → 200。
- 前端「Telegram channel」侧栏菜单将按 `features.menu_telegram` 显示。
- `dev.cmd start --no-browser`（默认 config 模式）实测：API readyz 200；testadmin 登录成功；`menu_telegram=True`；manifest 侧栏含 `telegram-settings`（visibleWhen `menu_telegram`）；schema 200；`dev.cmd stop` 正常停服。

## 评估

判据 #5 Admin UI 入口已可被用户实际访问（`dev.cmd start` 直读 config.yaml 即启用）；测试账户机制为本地/CI 验证提供稳定凭据，不触碰既有 admin。
