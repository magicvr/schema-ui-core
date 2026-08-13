---
id: A-002
goal: GOAL-005-r2-f03-account-center
title: S4 · self 审计（实现/验证与 D-002 一致性 + 安全面自查）
date: 2026-08-14
source: self
scope: S2/S3 实现与验证
verdict: pass
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-002 · S4 · self 审计

## 结论

**pass**（无 required finding；4 条观察均已在实现内处置或留痕）。

## findings

| id | 级别 | 内容 | 处置 |
|----|------|------|------|
| F-01 | info | 单条会话吊销不吊销对应 access token（短窗自过期）——D-002 `2 已文档化残余 | 留痕（D-002 `2） |
| F-02 | info | 停用账号的 access token 在中间件被 401 拒绝（而非 403）——避免状态 oracle，与 superseded 同 envelope | 实现即此语义（auth.go），测试覆盖 |
| F-03 | info | `users.enable` 同时覆盖 enable+unlock；服务端同一键、前端两个动作 | D-002 `1 留痕；R3 可拆键 |
| F-04 | info | users 页启停动作依赖 admin.account 路由（跨模块）；模块未启用时权限键不存在 → 按钮隐藏（fail-open 视觉 + 服务端 fail-closed） | 权限表达式门控实现；D-002 `5 留痕 |

## 偏差

无。实现与 D-002 冻结一致；安全敏感路径（停用即时失效、越权 403、外籍会话 404、改密吊销）全部有集成测试证据（E-004）。
