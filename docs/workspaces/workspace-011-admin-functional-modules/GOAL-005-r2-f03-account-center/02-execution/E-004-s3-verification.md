---
id: E-004
goal: GOAL-005-r2-f03-account-center
title: S3 · 验证（单测/集成 + 全量回归 + 本地冒烟）
date: 2026-08-14
status: recorded
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-004 · S3 · 验证

## 事实（2026-08-14）

| 项 | 结果 | 说明 |
|----|------|------|
| `go test ./... -count=1`（apps/api） | ✅ 全绿 | 含 store（迁移账本 0014）、handler（新增 `account_self_test.go`：改密后旧 refresh 401 / 停用后登录 403 / 解锁 423→200 / 越权 403 / 外籍会话 404 / 中间件停用拒绝）、kernel、composition（permissions 8/11、navigation 3/5） |
| `npm test`（apps/web） | ✅ 889/889 | i18n 结构测试扩展到 account.json + account fragment（13 页分母） |
| `npm run build` | ✅ | vite build 通过 |
| 本地 API 冒烟（admin profile，真实 SQLite） | ✅ | 登录→profile→sessions→改密（旧 token 401、新密码登录）→建用户→disable（登录 403 ACCOUNT_DISABLED）→enable（恢复）→5 次失败锁 423→unlock→恢复→editor 越权 disable 403→会话吊销 204/revoked→operationlog 事件（users.enable/disable/unlock、account.password-change）全落盘→schema/account 200、users actions 含 enable/disable/unlock |
| 容器 smoke（V-007/V-008） | 留待波次末尾统一回归 | 本目标未改容器装配；与 GOAL-003/004/006 一起在 R2 波次收尾时执行 |

## S3 门禁结论

S3 检查点完成：单元/集成（会话吊销后旧 token 失效 ✅、停用后登录拒绝 ✅）+ 全量回归 ✅。
