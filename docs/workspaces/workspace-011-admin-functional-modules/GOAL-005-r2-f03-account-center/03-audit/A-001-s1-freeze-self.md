---
id: A-001
goal: GOAL-005-r2-f03-account-center
title: S1 · 方案级 self 审视（D-002 冻结方案）
date: 2026-08-14
source: self
scope: S1 方案冻结
verdict: pass
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-001 · S1 方案级 self 审视

## 结论

**pass**（无 required finding；3 条观察均已在方案内处置）。

## findings

| id | 级别 | 内容 | 处置 |
|----|------|------|------|
| F-01 | info | 中间件对停用账号返回 401 UNAUTHENTICATED（与 superseded 同 envelope）属新增拒绝面：停用即时生效、不泄露状态 oracle | 方案已定（D-002 `3），S3 以集成测试验证 |
| F-02 | info | account 页放 navigation.user 区且对 viewer 可见：自服务页面人人可用，不暴露管理面 | 方案已定（D-002 `5），符合 PolicyAdminEditorViewer 语义 |
| F-03 | info | `users.enable` 同时覆盖 enable+unlock：键数最小化；若需分离审计可在 R3 增键（向后兼容，新增键仅影响授予） | 方案已定（D-002 `1），留痕 |

## 偏差

无。D-002 与 00-meta 边界一致（未改变登录/锁定协议语义；共享基架问题未引入）。