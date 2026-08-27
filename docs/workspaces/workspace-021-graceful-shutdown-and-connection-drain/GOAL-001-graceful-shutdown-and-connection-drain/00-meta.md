---
id: GOAL-001-graceful-shutdown-and-connection-drain
title: 优雅停机 / 连接排空合同
status: done
parent: null
created: 2026-08-27
updated: 2026-08-27
version: 0.3.0
progress: 3/3
plan_refs:
  - VP-021-graceful-shutdown-and-connection-drain
primary_plan: VP-021-graceful-shutdown-and-connection-drain
serves_summary: 承载 VP-021（架构 RT-D02 · 优雅停机 / 连接排空合同）实现：停机顺序合同、HTTP drain、运行中 Job 停机语义、双方言 Store 排空；单进程 + Compose 基线。消费 VP-012 Job 六态、VP-013 双方言 Store、VP-015 可观测面。不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。
---

# GOAL-001 · 优雅停机 / 连接排空合同

## 概述

本 Root 承载 [VP-021-graceful-shutdown-and-connection-drain](../../../vision/plans/VP-021-graceful-shutdown-and-connection-drain.md)（**`active`** v0.2.0 · 2026-08-27 激活，VRev-046 self `pass`）的实现：在既有进程生命周期代码路径（`cmd` / 进程入口）之上，把「进程生命周期有、但无明确 drain 合同」的后端收成可核对的**优雅停机 / 连接排空合同**——停机顺序、`http.Server` Shutdown 语义、运行中 Job（VP-012 六态）停机行为、双方言（SQLite / PG）Store 连接关闭顺序与迁移窗口重叠语义。

**激活门禁已全部满足**（2026-08-27）：VRev-046（self）`pass`（0 required；V-F081/V-F082 → 激活事务内 fixed）；架构类 freshness PASS（`ed99e88` → `fddaf638`，不暂挂 `go`）；VP-009/VP-010 无开放阻断。

**边界**：不承接 A3 余项（多实例、水平扩展、就绪探针扩依赖、PG `SKIP LOCKED` vs Redis vs 队列评估——均仍 trigger-gated）；API 与 worker 进程分离（RT-D03）；分布式锁 / Job 租约 / leader election（RT-Q04）；外部消息队列（RT-Q02）；K8s / Helm / Operator（RT-D06）；TLS 终止（RT-D05）；热配置 / 零停机部署产品；PITR；改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；重开 VP-012 / VP-013；业务域。

## 纲领路线图（P-001 · V-F081 fixed）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **合同冻结**：停机顺序（新请求停止 → 存量排空 → Store 排空 → 超时/退出码）；grace period / 超时默认值与配置键（I-002）；运行中 Job 停机语义（I-001：等完成 vs 中断标记重跑或按类型分流）；I-003（Store 排空与迁移窗口重叠）确认语义口径并登记（最晚 R2 关闭） | 起点 | **已关门**（GOAL-002 done 3/3 · A-001 self `pass` · 合同 v0.1.0 = GOAL-002 D-002） |
| R2 | **实现与测试**：`http.Server` Shutdown 合同化（grace / 超时 / 退出码）；Job 停机行为实现（R1 冻结语义）；双方言连接关闭顺序 + 迁移窗口重叠时停机语义（I-003 required，方案冻结前关闭）；复用 VP-015 结构化日志 / correlation | 依赖 R1 | **已关门**（GOAL-003 done 3/3 · A-001 self `pass` · shutdown_timeout 配置键 + main 接线 + compose 15s + 测试锁；I-003 已 verified 承接） |
| R3 | **证据与关门**：SIGTERM / SIGINT → 排空 → 退出码可核对（信号测试或等价 harness，单进程 + Compose）；SQLite / PG 双方言排空一致；checksum 台账不变；退出判据 1～5；开放 required = 0 | 依赖 R2 | **已关门**（GOAL-004 done 3/3 · A-001 self `pass` + A-002 grok independent `conditional`（F-001/F-002 required → 已 fixed）· harness A/B/C + PG 实测 PASS + 新请求拒绝断言） |

`progress` = 已关门纲领阶段数 / 3。当前 **3/3**（2026-08-27：R1、R2、R3 全部关门；Root done）。

## 成功标准（方向级 · 与 VP-021 退出判据镜像）

1. 停机顺序 / 超时 / 退出码合同落盘，且单进程 + Compose 下可核对（信号测试或等价 harness）。
2. 运行中 Job 的停机语义已冻结并有明确行为证据。
3. 双方言（SQLite / PG）Store 排空语义一致可核对；checksum 台账不变。
4. 未进 A3 余项；未改 Charter；未改 Profile 默认集作为本波成功条件。
5. 开放 required finding = 0（或已合法闭合）。

## 信息就绪与未知项

与 VP-021 I-021-00X 同号镜像（I-00X ↔ I-021-00X）。禁止在 R1 关闭前直接改进程生命周期实现 / 迁移台账相关 DDL。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 停机时运行中 Job 语义：等完成 vs 中断标记重跑（或二者按 Job 类型分流） | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：**中断标记重跑**（GOAL-002 D-001 accepted；合同 §4） |
| I-002 | required | grace period / 超时默认值与可配置键（含超时后的强制退出语义） | 方案冻结 | R1 | lead 提案 + 用户确认 | **verified** | — | 2026-08-27 用户裁决：默认 `10s` + `http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`；非法值 fail-closed（GOAL-002 D-001；合同 §6） |
| I-003 | required | Store 排空与迁移窗口重叠时的停机语义（fail-closed？排队？） | 方案冻结 / 实施 | R2 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：**fail-closed 启动期校验**，无运行时迁移窗口（GOAL-002 D-001；合同 §5） |
| I-004 | non-blocking | 停机是否需日志 / 指标断言（消费 VP-015 已交付面） | 验收 | R3 | lead 提案 | **verified** | — | lead 口径：结构化日志三事件断言；指标不进分母（GOAL-002 D-001；合同 §7） |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-021-graceful-shutdown-and-connection-drain）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。