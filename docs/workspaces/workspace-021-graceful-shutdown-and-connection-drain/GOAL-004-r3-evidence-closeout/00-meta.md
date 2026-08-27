---
id: GOAL-004-r3-evidence-closeout
title: R3 证据与关门（合同 §8 harness · 双方言 · Root 关门）
status: active
parent: GOAL-001-graceful-shutdown-and-connection-drain
created: 2026-08-27
updated: 2026-08-27
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-021-graceful-shutdown-and-connection-drain
primary_plan: VP-021-graceful-shutdown-and-connection-drain
serves_summary: 承载 VP-021 R3 阶段：合同 v0.1.0 §8 的证据 harness（clean drain / timeout / 重启 reclaim）+ SQLite·PG 双方言 + compose stop 核对；Root 关门审计与 VP-021 退出判据 1～5 对照。
---

# GOAL-004 · R3 证据与关门

## 概述

按合同 v0.1.0 §8 实施进程级证据 harness（GOAL-003 A-001 F-001 承接；I-004 验收口径 = 结构化日志断言，指标不进分母）：单进程 + 信号 → 排空 → 退出码可核对；双方言（SQLite / PG）一致性；迁移台账 checksum 不变；随后 Root 关门审计（self + grok build independent）与 VP-021 退出判据对照。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **harness 方案**：A（clean drain：存量请求 + 运行中 Job → SIGTERM → exit 0）、B（timeout：1s 预算 + 拖住连接 → exit 1 + `shutdown.timeout` 日志）、C（重启 reclaim：停机中断 Job → 重启 attempt+1）；双方言矩阵；compose stop 核对方式 | 进行中 |
| C2 | **实施与证据**：harness 落地（进程级子进程或等价 test harness）；SQLite + PG 各跑；退出码 / 日志 / Job reclaim 断言；checksum 回归锁 | 待定 |
| C3 | **关门审计与 Root 关门**：A-001 self + A-002 grok build independent（grok-4.6 · high）双审；VP-021 退出判据 1～5 对照；Root done；决策层 VP-021 关门提案 | 待定 |

`progress` = 已关门检查点数 / 3。当前 **0/3**。

## 成功标准（与 VP-021 退出判据 1～5 及合同 §8 镜像）

1. 停机顺序 / 超时 / 退出码合同在单进程 + Compose 路径可核对（harness A/B）。
2. 运行中 Job 停机语义有明确行为证据（harness C：中断 → 重启 reclaim）。
3. 双方言（SQLite / PG）Store 排空语义一致可核对；checksum 台账不变。
4. 未进 A3 余项；未改 Charter；未改 Profile 默认集作为成功条件。
5. 开放 required finding = 0（自审 + 独立审）。

## 信息就绪与未知项

I-001～004 已 verified（R1 裁决；I-004 验收口径 = 结构化日志断言）。无新增开放 required；若 harness 暴露实现/合同缺口，按「发现后回流」暂停并记录。