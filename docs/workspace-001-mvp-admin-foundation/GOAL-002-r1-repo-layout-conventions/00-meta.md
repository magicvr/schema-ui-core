---
id: GOAL-002-r1-repo-layout-conventions
title: R1 · 仓库布局与包管理约定
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# GOAL-002 · R1 · 仓库布局与包管理约定

## 概述

在本仓将 R1 的 monorepo 目录、包管理与本地运行约定落成**可引用文档与最小占位**，使 API/Web 脚手架（GOAL-003 / GOAL-004）有唯一布局真相，而不实现业务能力。

依据 Root `D-004`（I-STACK-001 / I-STACK-002）。

## 成功标准

- [ ] 根或 `docs/architecture/` 中有 monorepo 约定：`apps/web`、`apps/api`、与 `docs/`/`skills/` 边界说明
- [ ] 前端包管理写明：npm + `package-lock.json`（工作目录 `apps/web`）
- [ ] 后端包管理写明：Go modules（工作目录 `apps/api`，独立 `go.mod`）
- [ ] 本地运行入口（命令级）在根 README 或约定文档可找到；标明「骨架阶段，业务未实现」
- [ ] 未把订单/钱包/通知等业务域目录当作本仓 MVP 默认树

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 本目标范围由 Root D-004 锁定；无新增 required | — | — | — | — | — | I-STACK-001/002 已在 Root verified |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。与 GOAL-003 / GOAL-004 同阶段；本目标宜先完成文档约定，二者可并行 scaffold。
- 平行仓仅作结构参考，不整仓拷贝（见 Root D-004）。
