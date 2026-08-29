---
id: GOAL-004-r3-compose-cicd
title: R3 · compose/CI 实跑（主仓 compose 全服务 + consumer-regression 免凭据重构 + 信号级 drain harness linux 容器）
status: done
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-004 · R3 · compose/CI 实跑

## 概述

承接 Root R3 与 VP-024 判据 #3：主仓 compose/Dockerfile **真实实跑**（healthz/readyz + `stop_grace_period` + shutdown drain 日志）；golden-field `consumer-regression` workflow 按 v0.4.0 形态重构（**npmjs 公开免凭据** + pnpm/setup-go 就位 + serve 化冒烟）；R1 登记的**信号级 drain harness**（SIGTERM → exit 0/1 · `shutdown.timeout`）在 linux 容器补跑。核销 workspace-023 GOAL-005 F-001（compose 未实跑）与 A-002 F-007（workflow 无 pnpm setup · 跨仓 token）残留。

## 成功标准（可验证检查点）

- [x] C1：主仓 `docker compose up -d --build`（api + web 全服务）→ api healthy（readyz `{"status":"ok"}`）· web 反代 200；`docker compose stop`（grace 15s）→ 日志含 `shutdown.starting`/`shutdown.complete` 且 api **ExitCode 0**
- [x] C2：golden-field `consumer-regression.yml` 重构提交（免 GH token · setup-node pnpm cache · Go build + serve 后台 + healthz/readyz + 四探针 + SIGTERM 收尾断言 `shutdown.complete`）；**本地等价实跑全绿**（Windows 无 SIGTERM → 信号面由 C3 容器证据 + hosted 触发登记覆盖）
- [x] C3：信号级 drain harness（linux 容器 = CI 等价）：A = stop → `shutdown.starting`+`shutdown.complete` + **exit 0**；B = `HTTP_SHUTDOWN_TIMEOUT=1s` + 在途慢请求 → `shutdown.timeout` + **exit 1**
- [x] C4：I-024-002 核销（CI 环境 = 本地等价 + linux 容器实跑）；**hosted runner 实触发登记**为 R7 复核查项（不主张 hosted acceptance——workspace-002 先例）

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | workflow 重构 + compose 实跑准备（env 注入不落盘） | **已关门**（2026-08-29 · D-001 · workflow commit `c4d14ea`） |
| S2 | compose 实跑（api+web · healthz/readyz · stop drain 日志） | **已关门**（2026-08-29 · E-002） |
| S3 | drain harness A/B（linux 容器）+ workflow 本地等价实跑 | **已关门**（2026-08-29 · E-002） |
| S4 | 证据 + 自审 + 独立审计（grok）→ 关门 | **已关门**（2026-08-29 · A-002 grok independent `pass`（0 required）· F-001~F-003 fixed · GOAL-004 done） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-024-002 | required | CI 槽位环境（真实 runner / 用户环境等价 + 凭据） | 判据 #3 | R3 | workflow 实跑验证（本地等价 + linux 容器；hosted 触发登记） | **verified**（2026-08-29 · 本地等价 + linux 容器实跑 · hosted 触发登记 R7） | — | E-002 |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 = independent（grok build · R2 先例）；hosted runner 实触发登记为 R7 复核项（不主张 hosted acceptance）。