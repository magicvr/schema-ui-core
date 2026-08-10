---
id: D-001
goal: GOAL-004-w3-security-audit-remediation
title: W3 修复范围与技术取舍
date: 2026-08-11
status: accepted
---

# D-001 · W3 修复范围与技术取舍

## 背景

2026-08-11 对 api/web 的安全审计（主路径 + 四路专项）确认多项 findings。workspace-009 Root 保持 active；本决策冻结 W3 波次范围。

## 范围内（按优先级）

### P0

1. **batch-delete 非原子**：注释称 whole-batch，实现逐条 `Delete` 各自提交。  
   **取舍**：在 `users`/`roles` repository 增加单事务 `DeleteMany`；handler 优先走 `BatchDeleter` 接口；失败整批回滚。不引入跨资源 2PC。
2. **recordSource.url 接受 `//host`**：`startsWith("/")` 旁路严格正则，预填 `authFetch` 可带 Bearer 跨源。  
   **取舍**：与 rowAction/upload 一致，仅允许 `isRelativeProtocolUrl`（单斜杠同源路径，可带 query）。

### P1

3. **nginx**：`client_max_body_size 8m`；`server_tokens off`；基线 `nosniff` / `X-Frame-Options DENY` / `Referrer-Policy` / 保守 CSP；`location /api/` 前缀。
4. **登录限流**：peer 为 loopback/private 时信任 `X-Real-IP`；失败键 = IP + username；**成功登录清桶**；map 容量上限驱逐最旧键。
5. **委托边界**：非 admin 不得对 admin 目标改密码；非 admin 不得从角色集中移除 `admin`（与「仅 admin 可授 admin」对称）。

### P2

6. logo 同源分支拒绝 `\`（防 `/\evil.com` 浏览器解析为外站）。
7. logout 等待/作废 in-flight refresh（generation 计数），避免 token 写回。
8. `srv.Serve` 非 `ErrServerClosed` 错误 → 进程退出（fail-closed，配合 compose restart）。

## 明确不做（本波）

- Cookie 化 refresh、TLS 终结、多实例限流后端、精细生产 CSP nonce。
- 改变 D3 schema 匿名可读 residual。

## 审计模式

security 高影响实施 → 完成后 **self** 审计（A-001）；独立审由编排器按 VP-009 provider 另开（不阻断本波落地）。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| batch 仅「先 Get 再删」两阶段无事务 | Get 不覆盖 LAST_ADMIN/SELF；仍有 TOCTOU |
| 限流直接信任 X-Forwarded-For | 客户端可伪造；仅 private peer 下 X-Real-IP |
| 密码重置拆独立 permission | 范围过大；本波用 admin 角色门槛收口 |
