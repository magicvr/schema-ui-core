---
id: E-001
title: 立项：R3 证据 harness（合同 §8）
date: 2026-08-27
status: done
---

# E-001 · GOAL-004 立项（2026-08-27）

## 事实

1. Root 纲领 R3 → `GOAL-004-r3-evidence-closeout` 立项（五件套；检查点 C1 harness 方案 / C2 实施与证据 / C3 关门审计与 Root 关门；progress 0/3）。
2. 承接：GOAL-003 A-001 F-001（进程级 harness）；合同 §8（A/B/C harness + 双方言 + compose）；I-004 口径（日志断言，指标不进分母）。
3. goal-tree 同步（R3 挂树；Root 纲领 2/3 → 3/3 待 R3 关门）。
4. 关门审计计划：A-001 self + A-002 grok build independent（grok-4.6 · high，本地 build）。

## 验证 / 后续

- C1：harness 方案落盘（D-001）——A clean drain（存量请求 + 运行中 Job → SIGTERM → exit 0）；B timeout（1s 预算 + 拖住连接 → exit 1 + shutdown.timeout）；C 重启 reclaim（attempt+1）；SQLite + PG 矩阵；compose stop 核对。