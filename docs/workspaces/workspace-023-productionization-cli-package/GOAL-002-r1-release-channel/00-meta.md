---
id: GOAL-002-r1-release-channel
title: R1 · 真实发布通道
status: done
parent: GOAL-001-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 4/4
---

# GOAL-002 · R1 · 真实发布通道

## 概述

承接 Root R1 与 VP-023 判据 #1：Go/npm 真实发布通道闭环——origin tag → 公共 Go proxy `go get` 消费；npm registry（GitHub Packages）上传 → `pnpm add @ver` 消费；golden-field 移除全部占位依赖（replace/file:）。

## 成功标准（阶段检查点）

- [x] **S1 · Go 通道实证**：tag 命名教训（`apps/api/v0.1.0` 而非 `v0.1.0`——子目录 module path 约定）+ 公共 proxy 真实拉取 + `go.sum` 哈希 + golden-field 运行全绿（E-002）；sumdb 收录时延 = 知识项，后续默认校验复核
- [x] **S2 · npm 通道实证**：GitHub Packages 双包发布（凭据 = 用户 .env `github_token`）+ golden-field registry 安装（lockfile tarball+integrity 证据）· 三探针全绿（E-002）
- [x] **S3 · 占位依赖清零**（go.sum/锁文件 → registry 语义；升级演练绑定 R2 发布，A-001 F-001）：golden-field 移除 file: tarball → registry 安装；一次 registry 升级（含 breaking 场景预演）零冲突
- [x] **S4 · 关门**：A-001 self `pass`（0 required）→ **判据 #1 满足声明**；GOAL-002 `done 4/4`

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-023-001 | required | Go 发布通道形态 | S1 | R1 | **verified**（origin tag `apps/api/v0.1.0` + 公共 proxy 实证 + sumdb 时延知识项） | E-002 |
| I-023-002 | required | npm registry 目标 + 发布凭据 | S2 | R1 | **verified**（Facebook Packages · scope `@magicvr` · 凭据 .env 注入 · 发布+安装实证） | E-002 |
| I-023-003 | required | PG 实例 | R4 实测 | R4 | open | — |

## 父目标

- `GOAL-001-productionization-cli-package`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。