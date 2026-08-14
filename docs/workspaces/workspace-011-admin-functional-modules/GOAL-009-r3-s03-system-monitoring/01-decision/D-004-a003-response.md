---
id: D-004
goal: GOAL-009-r3-s03-system-monitoring
title: A-003 响应：1 required + 1 recommended 全 fixed
date: 2026-08-14
status: accepted
parent: GOAL-009-r3-s03-system-monitoring
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-004 · A-003 响应（关门）

## 结论

A-003（grok-build independent，security/data）verdict **conditional**、1 required + 1 recommended —— **全部 fixed**：

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| F-001 status 扁平对象非 list envelope → statCard fail-closed | required | **fixed** | handler 改单行 envelope {items:[row],total:1,page:1,pageSize:1}；row 增加 moduleCount（数值）供 valueField；schema valueField modules→moduleCount；Host 侧单测（resource.test.ts：真实 envelope 过 fetchResourceList 读全部 valueField + 扁平体 reject） |
| F-002 GetOperation 包裹哨兵 → 等值比较失效返回 500 | recommended | **fixed** | monitoringEntity.Get 改 errors.Is（与 activity 一致）；新增 404 OPERATION_NOT_FOUND 测试 |

## 验证

- 修复后回归：go test ./... 全绿、vitest 900/900、e2e mvp/admin 8/8。
