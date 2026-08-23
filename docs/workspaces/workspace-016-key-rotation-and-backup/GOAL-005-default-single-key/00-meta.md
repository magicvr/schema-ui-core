---
id: GOAL-005-default-single-key
title: R4 默认单密钥仍可用（证据整合）
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
serves_summary: 承接 Root R4：证明未配置 previous 时本地/Compose 默认仍能开发与快测，轮换不是启动硬依赖。以既有证据映射 + compose 配置解析实证为主。
---

# GOAL-005 · R4 默认单密钥仍可用（证据整合）

## 概述

承接 [workspace-016] GOAL-001 纲领阶段 **R4**（依赖 R2 ✓）。VP 方向级判据 2：**未配置 previous 时，本地/Compose 默认仍能开发与快测；轮换不是 mvp/dev 启动硬依赖**。R1～R3 的实现与测试已覆盖绝大部分面，本目标把判据映射到可核对证据并补齐唯一缺口（compose 配置解析实证），不新造机制。

## 边界

- 不改任何产品代码与配置默认值；只做证据核对与记录。
- 不强制 Compose 常驻双密钥、不引入备份 sidecar。

## 检查点（progress 来源）

| # | 检查点 | 状态 |
|---|--------|------|
| 1 | D-001：判据 2 → 证据映射表（config / composition / compose / dev 启动四层） | done（D-001） |
| 2 | E-001：证据实跑核对（含 `docker compose config` 解析实证：未设 PREVIOUS 时配置合法且值为空） | done（E-001，6/6 成立） |
| 3 | self 审计 + goal-tree 同步 | done（A-001 pass · 0 required） |

`progress` = 已完成检查点 / 3 = **3/3**。

## 信息就绪

无新增信息项；本阶段不依赖开放 required 信息门禁。

## 父目标

- GOAL-001-key-rotation-and-backup（[Q2](../GOAL-001-key-rotation-and-backup/00-meta.md)）

## 台账布局

三台账目录平铺追加；索引文件保留 frontmatter 与条目索引。
