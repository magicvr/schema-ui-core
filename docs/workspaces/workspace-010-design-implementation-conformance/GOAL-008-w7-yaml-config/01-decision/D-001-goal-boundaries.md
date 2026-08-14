---
id: D-001
goal: GOAL-008-w7-yaml-config
title: 立项边界：YAML 主配置体系（W7）
date: 2026-08-14
status: accepted
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（W7 YAML 配置体系）

## 1. 背景与归属

- 用户 2026-08-14 裁决：主流配置分层 = **config.yaml 承载全部配置（敏感字段 `${VAR}` 占位符）+ .env/env 存真实敏感值**；当前纯 env 配置不符合主流 → 归 workspace-10（VP-010 共享基架整改，workspace-11 已声明基架问题回流 VP-010）。
- 与项目现状一致：compose.yaml 已用 `${AUTH_JWT_SECRET:?...}` fail-closed 插值 + .env（gitignore）。

## 2. 目标边界

- 交付：应用层 config.yaml 加载（非敏感直接值 + 敏感 `${VAR}` 引用）；默认配置文件（仓库内）；.env.example；未定义变量 fail-closed；compose/smoke/e2e/文档同步；现有纯 env 部署兼容迁移路径。
- 排除：不引入 vault/secret-store 集成（env 注入即敏感来源）；不改协议 pin / Profile 默认集 / 模块矩阵语义。

## 3. 关键设计约束（S1 冻结细化）

- 优先级建议：config.yaml 为默认源；敏感项仅 env（YAML 引用）；env 可覆盖非敏感项（12-factor 允许）。
- 插值：Go os.ExpandEnv 语义；未定义且无默认 → fail-closed（compose :? 先例）。
- 迁移：现有 env 名保持可读（YAML 缺失时回退纯 env 读取，或文档化一键迁移）。

## 4. 信息就绪

I-001/I-002/I-003（S1）、I-004（S4）见 00-meta；均 open，S1 方案冻结关闭前不得实施。
