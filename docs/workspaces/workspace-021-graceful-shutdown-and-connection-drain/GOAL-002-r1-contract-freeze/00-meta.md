---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（停机顺序 / 超时与配置键 / Job 语义 / Store 排空）
status: active
parent: GOAL-001-graceful-shutdown-and-connection-drain
created: 2026-08-27
updated: 2026-08-27
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-021-graceful-shutdown-and-connection-drain
primary_plan: VP-021-graceful-shutdown-and-connection-drain
serves_summary: 承载 VP-021 R1 阶段：把现行进程生命周期收成可核对的优雅停机/连接排空合同——停机顺序、HTTP drain、运行中 Job 停机语义、双方言 Store 排空与迁移窗口关系、超时默认与配置键、退出码语义。
---

# GOAL-002 · R1 合同冻结

## 概述

本目标执行 Root 纲领 **R1**：在既有实现事实（`cmd/server/main.go` 信号处理 → `composition.registerLifecycle` OnStop 顺序 → `jobs.Runner.Stop` 取消语义 → `store.Open` 启动期迁移）之上，冻结 VP-021 合同正文——**停机顺序 / HTTP drain / 运行中 Job 语义 / 双方言 Store 排空与迁移窗口 / 默认超时与配置键 / 退出码语义**。合同正文 = GOAL-002 D-002 产物；不在此目标内改生命周期实现或迁移台账（V-F081 约束）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-001（Job 停机语义）、I-002（grace/超时默认与配置键）、I-003（Store 排空 × 迁移窗口）三条 required 由用户裁决；I-004（日志/指标断言）确认口径 | 进行中（D-001 提案 → 用户裁决） |
| C2 | **合同正文**：停机顺序、HTTP drain（`http.Server` Shutdown 语义 / grace / 超时 / 强制退出）、Job 停机行为、双方言连接关闭顺序、迁移窗口关系、退出码语义落盘（D-002） | 待定 |
| C3 | **审视与关门**：合同自审（self）+ 可选独立审；R1 关门、Root 信息台账回写 | 待定 |

`progress` = 已关门检查点数 / 3。当前 **0/3**。

## 成功标准（方向级）

1. 合同正文明确停机顺序：停止接收新请求 → 存量请求排空（grace）→ Job / 后台任务处理 → Store 排空 → 超时与退出码语义。
2. HTTP drain 与既有本地双进程 / Compose 路径（RT-D01 `delivered`）对齐；预算与配置键可核对。
3. 运行中 Job（VP-012 六态）停机行为冻结（等完成 / 中断标记重跑其一），并有现有 runner 事实佐证。
4. SQLite / PG 双方言排空语义一致可核对；迁移窗口与停机的边界明确；checksum 台账不变。
5. 未越界：不进 A3 余项；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；无生命周期实现变更。

## 信息就绪与未知项

与 Root / VP-021 同号镜像（I-001 ↔ I-021-001，…）。C1 未关闭前不得冻结合同正文（方案冻结门禁）。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 停机时运行中 Job 语义：等完成 vs 中断标记重跑（或按 Job 类型分流） | 方案冻结 | C1 | 用户裁决 | **collecting** | — | 建议「中断标记重跑」（D-001 证据） |
| I-002 | required | grace period / 超时默认值与可配置键（含超时后强制退出语义） | 方案冻结 | C1 | 用户裁决 | **collecting** | — | 建议默认 10s + `http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`（D-001 证据） |
| I-003 | required | Store 排空与迁移窗口重叠时的停机语义（fail-closed？排队？） | 方案冻结 / 实施 | C1（口径）/R2（关闭） | 用户裁决 | **collecting** | — | 建议 fail-closed = 启动期校验（D-001 证据） |
| I-004 | non-blocking | 停机是否需日志 / 指标断言（消费 VP-015） | 验收 | C2/C3 | lead 提案 | **collecting** | — | 建议：结构化日志断言；指标不进分母 |

## 父目标

- `GOAL-001-graceful-shutdown-and-connection-drain`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。