---
id: GOAL-004-r3-recovery-evidence
title: R3 轮换后恢复证据（SQLite + PG）
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
serves_summary: 承接 Root R3：在既有 SQLite VACUUM INTO 与 PG pg_dump/pg_restore 路径上核对「轮换后从备份启动 + 鉴权合同」。不重做 dump。
---

# GOAL-004 · R3 轮换后恢复证据（SQLite + PG）

## 概述

承接 [workspace-016] GOAL-001 纲领阶段 **R3**（依赖 R2 ✓）。备份/恢复机制本身由 VP-013 交付（SQLite `VACUUM INTO` 快照 + PG `pg_dump -F c`→`pg_restore` 合同），本目标只补**轮换语义组合证据**：密钥不在库中，但必须证明「旧密钥时代的数据备份 × 新密钥配置」组合下应用完整启动且鉴权合同成立。

## 边界

- 不实现应用内 Backup API、不重做 dump 工具、不改持久化端口。
- 不做 PITR / `pg_basebackup` / 备份代理 / SQLite→PG 搬运器。

## 检查点（progress 来源）

| # | 检查点 | 状态 |
|---|--------|------|
| 1 | D-001 决策关闭 I-004（最小恢复剧本：备份点相对轮换点、两方言命令、鉴权断言） | done（D-001；Root I-004 verified） |
| 2 | 可重复自动化：composition 级恢复循环测试（SQLite 全量必跑；PG 按 pgtest 门控 + dump/restore 工具门控） | done（E-001） |
| 3 | 双方言实跑证据落盘（E 记录含命令与输出摘要） | done（E-001：SQLite PASS + PG PASS 实跑记录） |
| 4 | self 审计 + goal-tree 同步 | done（A-001 pass · 0 required） |

`progress` = 已完成检查点 / 4 = **4/4**。

## 信息就绪

I-004 由本目标 D-001 关闭。无其他新增未知。

## 父目标

- GOAL-001-key-rotation-and-backup（[Q2](../GOAL-001-key-rotation-and-backup/00-meta.md)）

## 台账布局

三台账目录平铺追加；索引文件保留 frontmatter 与条目索引。
