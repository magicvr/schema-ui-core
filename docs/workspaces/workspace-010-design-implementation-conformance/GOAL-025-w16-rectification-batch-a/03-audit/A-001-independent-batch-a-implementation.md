---
id: A-001
goal: GOAL-025-w16-rectification-batch-a
title: 独立审计 · W16 批 A 实施（F01/F07/F08）
source: independent
date: 2026-08-17
verdict: conditional
scope: S2 实施（F01/F07/F08），不含 S3 全量回归/S4 关门
---

# A-001 · 独立审计 · W16 批 A 实施

- **source**: independent
- **auditor**: grok-build / grok-4.6 · reasoning high
- **type / scope**: execution-facts · GOAL-025 S2 implementation of W16-F01 / W16-F07 / W16-F08
- **verdict**: conditional

## 范围与区间

- covered: GOAL-025 五件套、D-001/E-001/E-002、F01/F07/F08 代码与前端实现
- excluded: 全量 Go/Web 回归、e2e、批 B/C

## 成果与证据

| 主张 | 证据 |
|------|------|
| F01 存储/门禁/重签实现 | `authsession/migration/migration.go`、`auth/auth.go`、`handler/account_self.go` |
| F07 端点与重签 | `handler/account_self.go` `revokeOthers` |
| F08 UI | `LoginPage.tsx`、`mfa-manager.tsx` |
| 定向测试通过 | `TestForcedPasswordChangeGateAndReissue`、`TestRevokeOthersReissuesTokensAndRevokesOtherSessions`、web 22/22 |

## Findings

### F-001 · Revoke-others does not refresh the session list
- level: **required**
- severity: med
- status: **open**（A-002 响应后 fixed）
- evidence: `account-session-toolbar.tsx` 仅刷新 `/me`，未触发 `reloadList`；会话表保持旧 active 行。

### F-002 · Forced password change accepts the same password
- level: **required**
- severity: med
- status: **open**（A-002 响应后 fixed）
- evidence: `changePassword()` 未拒绝 `newPassword == currentPassword`，强制改密可用原密码清标记。

### F-003 · Existing DBs do not backfill seed admin
- level: recommended
- severity: med
- status: open
- evidence: migration 0038 仅 ADD COLUMN DEFAULT 0；老库种子 admin 不自动置 1。

### F-004 · Related e2e still assumes seed admin can use business APIs
- level: recommended
- severity: med
- status: open（A-002 响应后 fixed）
- evidence: `shell.spec.ts` / `schema-crud.spec.ts` 直接用 admin/admin 进入业务面。

### F-005 · Revoke-others can race a concurrent 401 refresh
- level: recommended
- severity: med
- status: open
- evidence: bump 后旧 refresh 可能被并行 refreshAccess 清会话。

### F-006 · Forced change sets the “please sign in again” notice
- level: recommended
- severity: low
- status: open
- evidence: `auth-client.ts` 对所有 password 成功写 `password.changedNotice`。

### F-007 · Rotated MFA recovery codes have no download
- level: recommended
- severity: low
- status: open
- evidence: 仅 enroll 阶段提供下载，rotate 后无下载按钮。

## Required items

1. F-001：成功后 refetch 个人中心会话列表。
2. F-002：拒绝 newPassword == currentPassword。

## 结论

F01/F07/F08 主体实现完整，定向测试通过；因两条 required findings 未闭合，verdict 为 conditional，S2 尚不能视为完全接受，S4 关闭前需响应。
