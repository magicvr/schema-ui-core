---
id: GOAL-002-r1-tx-port-and-config
title: R1 · 内核 Tx 端口与配置键名冻结
status: done
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.2.0
progress: 2/2
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 冻结内核持久化端口（Tx ≠ *sql.Tx）与 db 配置键名；v1.1 补 upsert/时间/WasFresh/path 数据根。不实现驱动、不对写台账。
---

# GOAL-002 · R1 · 内核 Tx 端口与配置键名冻结

## 概述

Root 纲领 **R1**：把 VP-013 的持久化端口与配置面写成可实施合同。本目标**只冻结**，不改 `apps/api` 运行时。

权威正文：[attachments/r1-tx-port-and-config-freeze.md](attachments/r1-tx-port-and-config-freeze.md) **v1.1.0**（D-001 + D-002）。

## 范围

- `kernel.Store` / `kernel.Tx` 方法面、事务语义、占位符、禁止项。
- 配置键：`db.dialect` / `db.path` / `db.dsn` 与 env。
- 缺省仍为 SQLite 文件路径。
- v1.1（A-002 响应）：方言 SQL 硬规则（upsert / 时间）、`WasFresh` 方言中立语义、postgres 下 `db.path` 作数据目录根、`ErrNoRows` / `Dialect()` / `OpenOptions` 可扩展边界。

## 非目标

- 不选 PostgreSQL 驱动（I-002 → R2）。
- 不实现 postgres 方言、不对写迁移、不收口模块仓库签名（R2–R4）。
- 不改 Compose 默认、不引入 ORM。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 扫描现行 `db.path` / `WithTx(*sql.Tx)` 泄漏面 | ✅ E-001 |
| S1 | 冻结端口 + 配置键；self 审视 | ✅ D-001 / A-001；v1.1 补丁 D-002 / A-003（不另计检查点） |

## 成功标准

1. Tx 端口合同落盘：公共类型不含 `*sql.Tx` / 驱动类型；`Run` 为一事务。
2. 配置键名与缺省/fail-closed 规则落盘；与 VP-013 配置面一致。
3. 未修改应用代码。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行打开路径与 WithTx 形状是否足以冻结端口 | S1 冻结 | S1 前 | 读 store/config/composition | **verified** | 2026-08-20 | E-001：打开 + `WithTx`。A-002：关闭过窄（未含 upsert/时间/path 文件根）。**不另立 I-00N**：缺口写入合同 v1.1（D-002 / A-003 `fixed`） |
| I-002 | required | 本目标审计模式 | S1 关门 | S1 前 | 按内核契约、无代码 | **verified** | 2026-08-20 | D-001：当时 `self`。A-002 已补冻结层 independent；响应见 A-003。不回溯改本条 |

## 父目标

- `GOAL-001-store-dialects`

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
