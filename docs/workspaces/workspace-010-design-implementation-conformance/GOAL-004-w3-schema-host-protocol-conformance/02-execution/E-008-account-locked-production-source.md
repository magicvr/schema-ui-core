---
id: E-008
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: account-locked 生产源实现（用户 P-004 裁决：实现）— 锁策略 + 423 + Host 终态
status: recorded
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# E-008 · account-locked 生产源实现（用户 P-004 裁决：实现）

## 背景

A-008 F-1（account-locked 拟议 residual 无用户书面裁决）经用户 P-004 裁决：**实现生产源**
（不采取 accepted-residual）。D-002 §3 residual #1 随之关闭。

## 已完成事实

### 1. API 侧（锁策略 + 423 终态）

- **migration 12 `account_lock`**：`users` 表新增 `failed_login_count`（连续失败计数）与
  `locked_until`（unix 秒，0=未锁；到期自动解锁）；`TestCompiledMigrationCatalogOwnership`
  / `TestMigrateFreshDB` / `TestRestartPersistence` / operations 测试期望同步 11→12。
- **锁策略**（`internal/auth/auth.go`）：5 次连续密码失败 → 15 分钟锁窗（`LockThresholdFailures`/
  `LockWindow`）；锁期内登录直接 `ErrAccountLocked`（不做 bcrypt 工作）；锁开瞬间撤销该用户全部
  活动 refresh token（锁定账号不得继续轮换会话）；成功登录清零计数与锁窗。
- **仓储**（`authsession`）：`User` 增两字段（scan 全链路同步：accounts/userBy/ListUsers/
  UpdateUser）；`RecordLoginFailure`（返回是否开锁）、`ResetLoginFailures`、
  `RevokeAllRefreshTokensForUser`。
- **handler**：`POST /api/auth/login` 在 `ErrAccountLocked` 时返回 **423 `ACCOUNT_LOCKED`**
  （error catalog 双语注册：`error.accountLocked`）。
- **测试**：`TestAccountLockLifecycle`（5 失败→423、锁期内正确密码仍 423、到期自动恢复+清零）、
  `TestAccountLockRevokesSessions`（锁开撤销 refresh token）、`TestLoginRateLimit` 改用不存在
  用户（避免与账号锁交叉）。

### 2. Web 侧（锁定判别 + Host 终态）

- `auth-client.login`：423 → `AuthError("ACCOUNT_LOCKED", 423)`；测试 423 映射。
- `AuthContext`：`AuthStatus` 现由 `host/boot.ts` 的 `SessionAdapterState` 单一来源（ADR-0035
  D4 不变量收紧——映射不再分居两文件）；login 收到 423 时适配器进入 `locked`。
- `boot.ts`：`adapterAuthFor("locked")` → `{state:"locked"}`（归一化输入）；`lockedFailure()` →
  `HOST_ACCOUNT_LOCKED` 终态（`hostFailure.accountLocked`，home/support only、无 reauth、无
  retry——ADR-0036 D6）。
- `main.tsx`：`HostBootGate` 使用 `adapterAuthFor`；`AuthGate` 增加 `locked` 分支渲染
  `HostFailureScreen`（锁定终态优先于 anonymous 登录页）。
- i18n `hostFailure.accountLocked` 双语键已有（en/zh）。
- 测试：`boot.test.ts` locked 4 用例（含映射与 recovery 封闭断言）。

### 3. 门禁重跑

- apps/api：`go test ./internal/...` 全绿（auth/handler/authsession/store 含新用例）；
- apps/web：vitest 48 文件 **875 通过**；tsc 0 错误；
- Playwright：7 通过 + 1 既有 skip（真实双服务，bootstrap/host-failure/shell/schema-crud）。

## 阻塞 / 风险

无。剩余 O-001/O-002/O-003 为 recommended 观察项（vendored 行尾、app-navigation 2.7 pin），
不阻断 S6。

## 关联信息项 / findings

- A-008 F-1（account-locked P-004）→ **closed（已实现生产源）**。
- D-002 §3 residual #1（account-locked 生产源缺位）→ **关闭**。

## 下一步（计划）

claim 重生成（绑定本 commit）→ S6 台账终态（progress 6/6、status done、goal-tree/workspace/
索引同步）。
