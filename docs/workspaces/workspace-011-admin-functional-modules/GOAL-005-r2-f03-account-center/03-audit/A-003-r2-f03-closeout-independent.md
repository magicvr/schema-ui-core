---
id: A-003
goal: GOAL-005-r2-f03-account-center
title: S5 · 关门 independent 审计（grok-4.6 · security/authz · conditional）
date: 2026-08-14
source: independent
scope: S5 关门（提交 `5a50524` 声称范围；security/authz）
verdict: conditional
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-003 · S5 关门 independent 审计（grok build 代贴）

> provider：grok-4.6（Grok Build · 本地 CLI v1.0.3 · 只读：deny Write/Edit/Bash、--no-subagents、--disable-web-search）。原意见全文见会话记录；本文件为代贴摘要 + 完整 findings 台账。

## verdict：**conditional**

停用即时失效、权限 fail-closed、会话 IDOR、改密闭环、迁移账本、操作日志写入路径均成立；**last-admin 不变式不完整（F-001，required）**，不得无条件关门。

## findings（完整台账）

| id | severity | scope | 内容 | 编排器处置 |
|----|----------|-------|------|------------|
| F-001 | **required** | 启停守卫/last-admin/竞态 | `countAdminUsersExcluding` 不排除 `enabled=0`：委托 `users.disable` 可停用最后一名 **enabled** admin（停用 admin 仍被计入）；双管理员并发互停存在 TOCTOU；last-admin 测试与 self 测试同构，未独立打到分支 | **fixed**（D-004 / E-005）：enabled 计数 + 事务后不变式 + 独立用例 |
| F-002 | recommended | 启停错误契约 | 未知用户启停/解锁 → 500 INTERNAL 而非 D-002 承诺的 `USER_NOT_FOUND` 404 | **fixed**（D-004 / E-005） |
| F-003 | recommended | 改密 | 错误当前密码无限流/不计数（持有效 access 可在窗口内撞当前密码） | **fixed**（D-004 / E-005：复用 login limiter 模型，5 次/15 分钟 → 429） |
| F-004 | recommended | S3/S5 证据 | 冒烟未核 `account.session-revoke` 落盘；缺匿名 401 / 未知 404 / 73 字节密码 / 独立 last-admin 用例 | **fixed**（D-004 / E-005：补 6 个用例 + 冒烟核对） |
| F-005 | info | 会话模型 | 单条吊销不使对应 access 立即失效 | 与 D-002 `2 已文档化残余一致 |
| F-006 | info | 会话列表 | 过期未吊销显示 `status=active`（可幂等吊销） | 与 D-002 计算字段一致 |
| F-007 | info | 前端跨模块 | 键缺失时按钮渲染为 disabled 而非隐藏；detail 未展示 enabled/locked | 服务端 fail-closed 成立；视觉差异留痕（D-002 `5 语义为「隐藏」——实现为禁用，等效 fail-open 视觉） |

## 审计问题对照（8/8）

1. 停用即时失效：完整（Login 403 / Refresh 401 / Middleware 401 + tv+1 + refresh 吊销），未见绕过。
2. 权限 fail-closed：是（匿名 401、无键 403、PolicyAdmin 注册、模块关闭 404）。
3. self/last-admin：self 成立；last-admin 不完整 → F-001。
4. 会话 IDOR：否（双条件 + 外籍 404，列表无 hash）。
5. 迁移 0013/0014：未见破坏（checksum 冻结、CHECK 超集重建、编译目录不过滤 Profile）。
6. 改密闭环：成立（当前密码校验、8-72 字节、tv+1、refresh 吊销）；限流缺口 → F-003。
7. 前端跨模块：服务端 fail-closed 成立；视觉为禁用 → F-007。
8. 操作日志：写入路径完整；冒烟未核 session-revoke → F-004。

## 与既有意见

- A-001（self · S1 · pass）：方案级观察仍成立。
- A-002（self · S4 · pass）：在 last-admin 充分性与启停错误契约上**不一致**（独立审发现 F-001/F-002）；其余同向。
- 无 P-004 用户裁决冲突（F-001 采用 fixed，非 residual/overruled）。

## 建议

编排器：F-001 以 **fixed** 闭合后复跑全量测试再关门（D-004 / E-005 已响应）。
