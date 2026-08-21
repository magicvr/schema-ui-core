---
title: I-008-003 · 账号权限跨层集成 oracle（冻结）
status: frozen
doc_type: oracle
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-008-r6-integration-acceptance-vp-evidence
version: 0.1.0
related_info: I-008-003
---

# I-008-003 · 账号权限跨层集成 oracle（冻结）

> **已冻结（2026-08-01，I-008-003 verified / D-002 / A-002）**。本文件曾是「冻结候选」稿；正文中的候选表述为历史时点语义，权威冻结与执行结果以 [GOAL-008 03-audit](../03-audit.md)（A-002/A-003/A-005/A-006）、D-002/D-005 与 acceptance 证据为准。R6 已关门（GOAL-008 `done`）。

## 1. 范围与载体

- **身份载体**：R4 MVP 静态开发 session（`StaticDevSession`），经 `GET /api/accounts/me` 暴露，Web `loadAccountContext` 挂载为 `$context`（`user` / `features`）。无凭据库、无 token 签发。
- **覆盖层**：
  1. **Shell 导航权限**：`projectNavigation` 对 `permissions: { view: "$context.user.roles contains ..." }` 项做可见性投影（R3/R4）。
  2. **Renderer D-PERM**：`validatePermissions`（L2 fail-closed）→ `evaluatePermissionTargets`（AND 级联公式）→ `executeAction`（visible→permission→disabled→confirm）。（ADR-0023，GOAL-006 D-004 冻结）
  3. **HTTP 账号上下文**：`GET /api/accounts/me` 的 session JSON 形状。
- **依赖边界**：MVP 不依赖未声明业务领域模块；`dperm/cases.json` 17 例 fixture 以 SHA-256 pin 于 `permissions-inheritance.test.ts`。

## 2. 正向路径（allowed）

| # | 场景 | 入口 | 预期结果 | 证据 |
|----|------|------|----------|------|
| P-1 | dev session 暴露 `roles:[admin,editor]`、`features.beta:true` | `GET /api/accounts/me`（Web proxy → Go API） | HTTP 200；`user.id=dev-001`；`user.roles` 含 admin、editor；`features.beta=true` | `runtime-probes.log`（[3]）；E2E `shell.spec.ts` |
| P-2 | 拥有 admin 角色的导航项可见 | 真实 manifest 无 permission 门控项；用 `App.integration.test` 的同构断言（admin 注入 context 时 `Detail` 可见） | admin context → 门控链接可见 | `App.integration.test.tsx`（jsdom） |
| P-3 | D-PERM 有效权限：无 cascade 源 / 无目标权限时放行 | `evaluatePermissionTargets` + fixture cases（含 execution 时序断言 5 例） | 13 valid 求值 + 4 invalid 错误码符合 fixture | `dperm/cases.json`（17 例，SHA pin）；`permissions-inheritance.test.ts` |
| P-4 | 动作执行链通过（EXECUTED） | `executeAction`：visible+permission+confirm=true+confirmed=true | `outcome: EXECUTED`，events 含 `confirmShown`、`actionExecuted` | `permissions-inheritance.test.ts` |

## 3. 拒绝路径（denied / fail-closed）

| # | 场景 | 入口 | 预期结果 | 证据 |
|----|------|------|----------|------|
| D-1 | 无 admin 角色的导航项隐藏 | `App.integration.test`（viewer 注入 context 时 `Detail` 隐藏） | viewer context → 门控链接**不出现** | `App.integration.test.tsx` |
| D-2 | 缺失权限的 action 被阻断 | `executeAction`：effectivePermission=false | `outcome: BLOCKED`，`reason: PERMISSION_DENIED` | `permissions-inheritance.test.ts` |
| D-3 | 不可见目标被阻断 | `executeAction`：`visible=false` 或 target 不存在 | `outcome: BLOCKED`，`reason: NOT_VISIBLE` | `permissions-inheritance.test.ts` |
| D-4 | 结构校验 fail-closed | `validatePermissions`：版本过低 / 缺 capability / cascade 非法 / intent 挂载非法 | 返回对应 L2 错误码（`PROTOCOL_VERSION_TOO_LOW`、`CAPABILITY_REQUIRED`、`PERMISSION_CASCADE_*`、`PERMISSION_INTENT_*`） | `permissions-inheritance.test.ts`；`permissions.ts` |
| D-5 | 未知权限值拒绝 | `evaluatePermissionValue`：非 boolean / 非 string | 视为 `false`（deny） | `permissions.ts` `evaluatePermissionValue` |
| D-6 | 空 session provider fail-closed | `accountHandler.me` 的 `sessionProvider` 返回 `(nil, false)` | HTTP 401 `UNAUTHENTICATED` | `handler/account.go`（注入路径）；生产走静态 session 不触发 |

## 4. 排除 / 边界

- MVP 无凭据库与 token 生命周期：**不**在 R6 断言真实登录/登出、token 刷新或越权提权；`StaticDevSession` 是唯一运行时身份。
- 真实 manifest 当前**没有** permission 门控导航项（查看 `app-manifest.json` navigation 段）；P-2/D-1 的 shell 导航正向/拒绝以 `App.integration.test` 的注入 context 断言为 oracle，而不是浏览器内真实 manifest 的可见/隐藏（浏览器侧只能验证 session 已挂载）。
- 拒绝路径（D-2～D-6）以 D-PERM 单元/fixture 与 handler 注入测试为权威；浏览器 E2E 只覆盖正向 session 挂载（C-005），不伪造浏览器内拒绝。
- 上传域、多选批量、完整 component registry、跨业务模块权限均为既有排除（I-PROTO-001 v0.1.3）。

## 5. 阶段 2 执行对照

阶段 2 账号权限验收运行时，把 `dperm/cases.json` 的 executed/excluded 明细、`GET /api/accounts/me` 的 HTTP 结果与 E2E session 断言按证据 index 落盘；任何与上表不符的结果须显式列失败或排除，不得用总体 pass 掩盖。
