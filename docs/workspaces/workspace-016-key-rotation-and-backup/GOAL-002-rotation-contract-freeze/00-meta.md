---
id: GOAL-002-rotation-contract-freeze
title: R1 轮换合同冻结与配置面落地
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
serves_summary: 承接 Root R1：以 D-002 冻结轮换合同后，交付 previous 密钥的配置面（YAML/env 解析、ValidateProd 规则、单测）。不改 Authenticator 验签（属 R2）。
---

# GOAL-002 · R1 轮换合同冻结与配置面落地

## 概述

承接 [workspace-016] GOAL-001 的纲领阶段 **R1**。合同本体已在 Root `D-002`（accepted）冻结：

- current = `auth.jwt_secret` / `AUTH_JWT_SECRET`（不变）；previous = `auth.jwt_secret_previous` / `AUTH_JWT_SECRET_PREVIOUS`（新增）。
- 缺省单密钥；previous 已配置时非开发环境沿用与 current 相同的熵规则；同值守卫；secret 不入库不进日志；重启生效。

本目标把该合同落成**配置面代码 + 单测 + 文档注记**，并做 self 审计后关门。

## 边界

- **不改** `internal/auth` 验签逻辑、不改 composition 装配签名（R2/GOAL-003）。
- 不改 Compose 默认依赖；不做热加载；不做 Admin 密钥页。
- 安全 finding → VP-009；符合性 gap → VP-010。

## 检查点（progress 来源）

| # | 检查点 | 状态 |
|---|--------|------|
| 1 | Config 字段 + YAML/env 双通道解析 + 同值守卫与熵规则校验 | done（E-001） |
| 2 | 单元测试覆盖（缺省单密钥不变式 / previous 合规通过 / 弱 previous 与短 previous 拒绝 / 同值拒绝 / dev 不受影响） | done（E-001/E-002，8/8 + 9/9 + 全套件 exit 0） |
| 3 | 配置样例文档注记（config 样例或 README 键名表） | done（两 YAML 样例 + compose 透传 + README 两表） |
| 4 | self 审计（04，source=self）+ goal-tree 同步 | done（A-001 pass · 0 required） |

`progress` = 已完成检查点 / 4 = **4/4**。

## 信息就绪

无新增信息项；I-001/I-002 已在 Root 层 verified（D-002），本目标实施不依赖其他开放 required 信息项。

## 父目标

- GOAL-001-key-rotation-and-backup（[Q2](../GOAL-001-key-rotation-and-backup/00-meta.md)）

## 台账布局

三台账目录 `01-decision/`、`02-execution/`、`03-audit/` 平铺追加；索引文件保留 frontmatter 与条目索引。
