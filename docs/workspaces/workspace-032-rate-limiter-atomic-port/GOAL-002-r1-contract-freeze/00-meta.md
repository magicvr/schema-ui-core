---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（AllowRecord 端口契约）
status: done
parent: GOAL-001-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-03
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-032-rate-limiter-atomic-port
primary_plan: VP-032-rate-limiter-atomic-port
serves_summary: 承载 VP-032 R1：把 VRev-073 已冻结的 I-032-001/002 落成可执行端口合同（D-002）+ kernel.RateLimiter.AllowRecord + 编译期 stub + Memory 单锁实现与合同级测试。14 处使用点迁移归 R2。
---

# GOAL-002 · R1 合同冻结（AllowRecord）

## 概述

执行 Root 纲领 **R1**：在已 closed 的 VP-027 端口（Allow/Record/Clear/RetryAfterSeconds + 滑动窗口谓词）上新增原子方法 `AllowRecord`，消除 Allow→Record TOCTOU。信息项 I-032-001/002 已由 [VRev-073](../../../vision/reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) 激活冻结，本目标**不再裁决**签名或分母，只落盘合同与可编译端口面。内存供应商必须实现新方法（否则接口无法编译）；**生产 14 处调用点仍走 Allow+Record，归 R2 迁移**。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **继承激活冻结**：I-032-001/002 无新 P-004 裁决；记录到 D-001 | **已关门**（2026-09-03 · 用户指令冻结 R1 方案；无开放 required） |
| C2 | **合同正文 + 端口落地**：D-002 冻结；`kernel.RateLimiter` 增 `AllowRecord`；stub + Memory 单锁实现；合同级测试绿 | **已关门**（2026-09-03 · D-002 v0.1.0；kernel + Memory + 顺序等价/并发预算/`-race`） |
| C3 | **审视与关门**：阶段关门 self（Root D-001：R1 default self；independent 留 R3） | **已关门**（2026-09-03 · A-001 pass + A-002 独立审 F-001 fixed 闭合于 A-003；pass 0 required） |

`progress` = 已关门检查点数 / 3。当前 **3/3**。

## 成功标准（方向级）

1. 合同冻结：`AllowRecord(key string, now time.Time) bool` 的单锁语义、与 Allow-then-Record 的顺序等价、false 路径不登记 key，均写入 D-002 且可测试。
2. 兼容：`Allow`/`Record`/`Clear`/`RetryAfterSeconds` 语义不变；Retry-After 仅在 Allow **或** AllowRecord 返回 false 后调用。
3. 可编译：kernel stub 与 Memory 均实现新方法；`go test` kernel + ratelimit 绿。
4. 未越界：不迁移 14 处生产调用点；不实现 Redis；不改 Profile 默认集；不重开 VP-027。

## 信息就绪与未知项

与 Root / VP-032 同号镜像。无开放 required。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-032-001 | required | `AllowRecord` 精确签名与返回值 | 方案冻结 + 判据 1 | C1 | `/vision` VRev-073 | **verified** | — | `AllowRecord(key string, now time.Time) bool`；bool 足够；不返回剩余额度（D-001 / D-002 §1） |
| I-032-002 | required | 使用点是否全迁；Clear 是否需原子变体 | 方案冻结 + 判据 2 | C1 | `/vision` VRev-073 | **verified** | — | 14 处全迁（R2）；Clear 无需原子变体；立即消费 vs 失败预算两口径（D-002 §4/§5） |

## 父目标

- `GOAL-001-rate-limiter-atomic-port`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账。

## 备注

- 审计模式：Root D-001 — 阶段关门 default **self**；R3 证据/关门可按需 independent。本目标 C3 = self，不升格 cross（相对 VP-027 R1：本次是既有端口的加法，不是新公共面从零冻结）。
- R1 不关门生产 TOCTOU：调用点仍是 Allow→Record，直到 R2。
