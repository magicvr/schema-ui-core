---
id: E-005
goal: GOAL-005-r2-f03-account-center
title: A-003 修复 + S5 关门（required 全闭合 + independent 复核通过）
date: 2026-08-14
status: recorded
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-005 · S5 关门

## A-003 修复事实（2026-08-14）

| 修复 | 落地 |
|------|------|
| F-001 | `countEnabledAdminUsersExcluding` / `countEnabledAdminUsers`（join users.enabled=1）+ disable 后置不变式 + 独立用例 `TestDisableLastEnabledAdminRejected` |
| F-002 | 未知用户 → 404 `USER_NOT_FOUND`（`users_state.go` 三端点）+ `TestUserStateUnknownUser404` |
| F-003 | 改密错误当前密码限流（5/15min → 429）+ `TestPasswordChangeRateLimited` |
| F-004 | 匿名 401（`TestAccountEndpointsAnonymous401`）、73 字节（`TestAccountPasswordChangeTooLongNew`）、session-revoke 落盘（`TestSessionRevokeLogsOperation`） |

## 关门验证

| 项 | 结果 |
|----|------|
| `go test ./... -count=1` | ✅ 全绿（含新增 6 用例） |
| `npm test` / `npm run build` | ✅ 889/889 / 构建通过 |
| A-003 required 闭合 | ✅ F-001 fixed（D-004） |
| A-001/A-002/A-003 汇总 | 无开放 required；info 均留痕 |
| goal-tree / workspace.md | 同步（5/5 done） |

## 关门结论

GOAL-005-r2-f03-account-center **done（5/5）**：方案冻结 → 实现 → 验证 → go 判定（无影响不暂挂）→ independent 关门审计（conditional → required 修复后全绿）。
