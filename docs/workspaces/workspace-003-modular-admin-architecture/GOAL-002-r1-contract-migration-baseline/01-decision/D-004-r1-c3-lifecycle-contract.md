---
id: D-004
title: 冻结 C3 模块公共契约与生命周期语义
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
---

# D-004 · C3 模块公共契约与生命周期语义

## 决定

1. Go 1.26 是当前 API 的兼容基线；Uber Fx 只作为组合根的实现候选，模块公共 API 必须保持框架无关，Fx 类型不得泄漏。
2. 模块描述至少包含稳定 `id`、版本、内核 API 兼容范围和显式 `DependsOn`；标准 Admin 模块的 `HTTP`、`Schema`、`Authorization`、`Navigation`、`Manifest`、`Persistence` 六项是必须贡献，`Configuration`、`Lifecycle`、`Observability` 仅按需。
3. Profile 展开为显式 `modules.enabled` 后，依赖图、贡献冲突和 protocol/capability 不兼容必须在启动前 fail closed；不得静默自动启用依赖或把未知页面降级为局部可用。
4. 组合根按“注册/图校验 → 拓扑启动 → 反向停止”管理生命周期；启动失败清理已启动资源。`/healthz` 只表示进程存活，`/readyz` 必须聚合模块图、迁移、system-data reconcile 与必需依赖。
5. R1 只冻结语义和错误分类边界；稳定 Go 类型、Fx 版本、错误 code、模块级 hooks、聚合 readiness 与失败清理实现全部由 R2 承接。

## 理由与未选方案

- 采用 Fx 作为组合根候选符合现行架构文档，但当前仓库未有依赖或实现；把 Fx 放入模块 API 会锁定实现细节，因此不选。
- 当前进程级 `http.Server` 生命周期只能作为现状基线，不能替代模块级 lifecycle 语义；不把现有 `/readyz` 的 SQLite ping 解释成终态 readiness。
- 不在 R1 直接修改 Go 代码，避免把契约决策、R2 实现和实现验收混为同一门禁。

## 约束

本决定服务 Root I-003 的信息收集，不将其状态改为 `verified`；R2 子目标必须引用本决定并为每个边界提供实现/测试证据。
