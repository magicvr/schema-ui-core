---
id: GOAL-008-w7-yaml-config
title: W7 · YAML 主配置体系（config.yaml + env 仅敏感信息）
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 1/5
---

# GOAL-008 · W7 · YAML 主配置体系

## 概述

本子目标是 VP-010 / workspace-010 的**第七波**（共享基架整改，用户 2026-08-14 裁决）：当前配置全部走环境变量（`config.Load` 纯 env），不符合业界主流分层——**config.yaml 承载全部配置内容（非敏感直接写值，敏感字段写 `${VAR}` 占位符引用环境变量），真实敏感值由 .env（开发，gitignore）/ 进程环境变量（生产）提供**。此模式与项目现有 `compose.yaml` 的 `${AUTH_JWT_SECRET:?...}` fail-closed 插值先例一致，属把同一模式推广到应用自身配置。

## 当前边界

- 范围：`apps/api/internal/config`（YAML 加载 + `${VAR}` 插值 + fail-closed）、新增仓库默认 `configs/config.yaml`、`.env.example`、compose/smoke/e2e/部署文档同步、env 兼容迁移路径。
- **不**改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin / 运行时行为语义（配置载体变化，值语义不变）；go 判定在 S4 确认不暂挂。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：YAML 字段清单（14+3 项）、优先级与插值、.env、迁移兼容（D-002/A-001，2026-08-14）
- [ ] **S2 · 实现**：config.yaml 加载器 + 默认文件 + 迁移现有配置 + .env.example + compose 同步
- [ ] **S3 · 验证**：单元/集成（插值、fail-closed、env 覆盖）+ 双路径实测（纯 YAML / YAML+env）+ 全量回归
- [ ] **S4 · go 影响判定 + 自审**（配置载体变化 → go 判定）
- [ ] **S5 · 关门**：独立审计（grok，data 门禁）+ required 闭合 + goal-tree 同步

progress: 1/5 由五个等权检查点派生。

## 审计策略

独立审计沿用 grok build（用户书面偏好）；配置体系涉及部署/生产门禁，S5 独立审计。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | YAML 覆盖哪些配置项（现有 env 全集迁移？保留哪些 env 直读） | S1 方案 | 现 config 清单对照 | **closed**（D-002 §3：14 项 + upload 3 项） |
| I-002 | required | 敏感项清单与插值规则（未定义时 fail-closed 语义） | S1 方案 | compose :? 先例对照 | **closed**（D-002 §2） |
| I-003 | required | 现有部署只设 env 的迁移路径（YAML 缺失时行为） | S1 方案 | 向后兼容对照 | **closed**（D-002 §1：env 覆盖保留，零迁移） |
| I-004 | required | go 影响判定（Profile/模块矩阵语义不变） | S4 | VP-008 接口对照 | open |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
