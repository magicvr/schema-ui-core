---
id: D-004
goal: GOAL-005-r2-f03-account-center
title: A-003 响应 — F-001~F-004 全部 fixed（无 P-004 裁决冲突）
date: 2026-08-14
status: accepted
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-004 · A-003（grok independent · conditional）响应

> P-004 检查：F-001 为 required 且与 self（A-002）不一致。本响应选择 **fixed**（修复而非 residual/overruled）——安全不变式缺口，修复成本低、证据可核对；用户既定偏好（此前 A-002 的 required 均走 fixed）。无其它 required 冲突。

## 处置

| finding | 处置 | 修复 |
|---------|------|------|
| F-001 required | **fixed** | ① 预检查改为**只计 enabled=1 的 admin**（`countEnabledAdminUsersExcluding`）；② disable 事务内加**后置不变式**（更新后 enabled admin 计数为 0 → `ErrLastAdmin` 回滚）；③ 独立用例：先停用 admin2，再停用最后一名 enabled admin（user-admin）→ `ErrLastAdmin`，且用户保持 enabled。SQLite 单写者串行化使并发互停 fail-closed（busy 或后置检查） |
| F-002 recommended | **fixed** | users_state 对 `errResourceNotFound` 映射 404 `USER_NOT_FOUND`（enable/disable/unlock 三端点）+ 用例 |
| F-003 recommended | **fixed** | 改密端点错误当前密码复用 login limiter 模型（per client+user，5 次/15 分钟 → 429 `RATE_LIMITED`；成功 clear）+ 用例 |
| F-004 recommended | **fixed** | 补 6 用例：匿名 401（8 端点）、未知用户 404、73 字节密码 400、独立 last-admin、改密限流 429、`account.session-revoke` 操作日志落盘断言 |
| F-005/F-006/F-007 info | accepted（留痕） | 与 D-002 `2/`5 已文档化语义一致；F-007 视觉实现为 disabled（等效 fail-open），D-002 `5 表述「隐藏」以本留痕为准 |

## 验证

`go test ./... -count=1` 复跑全绿后关门（E-005 记录）。
