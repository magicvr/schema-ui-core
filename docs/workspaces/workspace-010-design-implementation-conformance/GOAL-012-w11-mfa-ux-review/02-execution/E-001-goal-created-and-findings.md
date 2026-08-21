---
id: E-001-goal-created-and-findings
doc: execution-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-001 · 目标建立与问题落盘（含 MFA 代码定位）

## 事实

1. **目标建立（2026-08-15）**：用户指令在 workspace-010 新增子目标，落盘 UX 审视改进项 + 3 项 MFA 缺陷。五件套与 ledger 目录已建；goal-tree / workspace 波次表已同步。
2. **M 类缺陷只读代码定位**（证据，未改动任何代码）：

### M-01 无二维码（前端缺失）
- API 已返回 otpauthURL：apps/api/internal/modules/mfa/service.go `Enroll()` 返回 `otpauthURL("Schema UI Core", name, secret)`。
- 前端只渲染文本：apps/web/src/components/mfa-manager.tsx（data-mfa-enroll 区）仅两个只读 Input（secretBase32 / otpauthURL）+ 恢复码 textarea，无二维码渲染。
- 结论：纯前端改动（由 otpauthURL 生成二维码），无需改 API。

### M-02 / M-03 输错码 → 强制登出（401 语义冲突）
- 服务端：apps/api/internal/handler/mfa.go `writeMFAError` 将 `ErrMFAInvalid` 映射为 **http.StatusUnauthorized（401）** + 码 `MFA_INVALID`（`/api/mfa/confirm`、`/api/mfa/disable`、`/api/mfa/recovery/rotate` 共用）。
- 前端：apps/web/src/account/auth-client.ts `authFetch`——非 auth 端点收到 **任意 401** → 尝试 refreshAccess() → refresh 失败或重试仍 401 → `clearTokens()` + `onAuthLost?.()` → AuthContext 切回登录页（强制登出）。
- 触发链（M-03）：mfa-manager.tsx `confirm()` → postJSON("/api/mfa/confirm") 输错码 → 服务端 401 → authFetch 登出。pending 状态仍在服务端，但用户已被登出，无法重填。
- 触发链（M-02 失败路径）：mfa-manager.tsx `disable()` → postJSON("/api/mfa/disable") 输错码 → 401 → 登出；MFA 仍 active（用户报告与之一致）。
- 触发链（M-02 成功路径的另一面）：服务端 `Disable` 成功 → `BumpTokenVersionAndRevokeAll` 吊销全部会话（含当前）→ 前端随后 `refresh()`（GET /api/mfa/status）401 → refresh token 已吊销 → 登出，**无任何成功提示**。
- 例外（不受影响）：登录二步验证 `/api/auth/mfa/verify` 由 `mfaVerify()` 直接处理 401 并抛 AuthError（不经过 authFetch，不清 token），故登录场景无此问题。

### 修复方向（待 I-001 用户裁决）
- 自服务端点 `ErrMFAInvalid` 改映射 **400**（业务校验失败语义），authFetch 不再误判为会话过期；
- 前端 mfaErrorKey 已有 MFA_INVALID → error.mfaInvalid 文案（仅需错误能到达页面）；
- 解绑成功路径：前端需在 `refresh()` 401 登出前给出明确提示（如「MFA 已解绑，当前会话已失效，请重新登录」）。

## 基线

- 本波开始前基线：go test 全绿（上波 W10 记录 2026-08-15 全量绿）；web vitest 991/991（W10 关门时记录）。本 E-001 无代码改动。
