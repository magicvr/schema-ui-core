---
title: S1 缺口盘点 · module-architecture / 代码路径 / QUICKSTART §5
status: active
created: 2026-08-06
updated: 2026-08-06
parent: null
version: 0.1.0
---

# S1 缺口盘点

对照 [module-architecture.md](../../../../architecture/module-architecture.md)、仓库真实模块路径与根 [QUICKSTART.md](../../../../../QUICKSTART.md) §5。

## 已具备（素材源）

| 区域 | 观察 | 路径 |
|------|------|------|
| 架构终态 | 核心六项、组合根、Profile、迁移、Manifest、横切边界已写清 | `docs/architecture/module-architecture.md` §1–§8 |
| 内核契约 | `Module` / `Provider` / `RegisterContributions` / Profile 默认集 | `apps/api/internal/kernel/` |
| 组合根 | 静态 import 一方模块；按 plan 装配 Provider | `apps/api/internal/composition/composition.go` |
| 一方标准模块 | users/roles/settings/activity 等 | `apps/api/internal/modules/{users,roles,settings,activity}/` |
| 横切 | operationlog | `apps/api/internal/modules/operationlog/` |
| 全局迁移收集 | 全候选 PersistenceProviders | `apps/api/internal/modules/compiled/persistence.go` |
| QUICKSTART §5 | 「加页面」级步骤：schema 目录 + Manifest/Navigation；禁 public 静态 Manifest | 根 `QUICKSTART.md` §5 |

## 缺口（S1 时）

| ID | 缺口 | 影响 |
|----|------|------|
| G1 | 无独立「接模块」MUST 清单（id/版本/六项/组合根/Profile/迁移/验证） | 合作者/AI 只能拼 architecture + QUICKSTART，易漏组合根/Profile/全局迁移 |
| G2 | 无对等 DO NOT 表（Renderer 中央注册、静态 Manifest、平行认证、按需误读、热插拔） | 过度施工风险 |
| G3 | 无可执行 Core vs 模块归属判定树与正反例入口 | 横切 vs 标准 Admin 易混淆 |
| G4 | overview / QUICKSTART 未链到操作 playbook（仅 architecture 终态 + 加页面） | 可发现性不足（VP-004 exit #4） |
| G5 | QUICKSTART §5 未升级到「完整一方模块」层级 | 读者可能以为只加 JSON 页面即可 |

## 非缺口

- 不需重开 VP-003 迁移或改 runtime 作为主交付
- 不需改 principles / workspace-protocol
- 不需默认脚手架

## 权威路径建议（供 D-002 冻结）

**新建** `docs/architecture/module-contribution-playbook.md` 为 must/must-not/归属法单一权威入口；由 `module-architecture.md` §9 链出；overview + QUICKSTART 接线。不把长 playbook 整段并入 architecture 决策正文，避免与 VP-003 终态边界文混读。
