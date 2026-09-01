---
id: GOAL-004-r3-seam-and-shared-conventions
title: R3 接缝与共享约定（Redis 接缝声明 / 轨道登记继承）
status: done
parent: GOAL-001-rate-limiter-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-027-rate-limiter-port
primary_plan: VP-027-rate-limiter-port
serves_summary: 承载 VP-027 R3 阶段（判据 #4/#5）：RateLimiter Redis 供应商接缝声明（端口不变 / 原子窗口 INCR+EXPIRE / 连接管理）延伸 owner 短文 cache-redis-seam-and-track.md v1.1.0 + Redis 轨道约定继承登记（rl 命名空间段 · §3.3 首条登记）。
---

# GOAL-004 · R3 接缝与共享约定

## 概述

执行 Root 纲领 **R3**：在 D-002 合同（v0.1.1）与 owner 短文 `docs/architecture/cache-redis-seam-and-track.md`（v1.0.0 · 单一所有者 · VP-026/027 轨道）之上，交付 VP-027 **判据 #4（Redis 供应商接缝声明）** 与 **判据 #5（共享约定登记）**：

- **接缝声明**：RateLimiter Redis 级供应商边界——端口不变（同一 `kernel.RateLimiter` / `RateLimiterProvider`）、key 映射沿用 `<ns>:<key>`（限流段 = `rl`）、**原子窗口原语 = INCR + EXPIRE**（Record = INCR+首 EXPIRE；Allow = 读；Clear = DEL）、滑动窗口实现选择触发立项时裁决、连接管理 = 组合根单一持有 + 构造 PING fail-closed、测试 harness 与缓存轨道相同；**不引入 Redis 客户端依赖**（`go.mod` 无 redis）。
- **登记与继承**：§3.3 命名空间登记表**首条登记** `rl`（7 处使用点 · 归属 VP-027）——履行 workspace-026 关门登记的「命名空间登记义务于 VP-027 激活触达」；§1 端口分母增列 `kernel/ratelimit.go`；修订史 v1.1.0。
- 本目标**零 Go 代码变更**（接缝声明与登记全部落架构文档层）。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **接缝声明落盘**：D-001 决策 + owner 短文 v1.1.0 §2.6（端口不变 / key 与原子窗口 / 连接管理 / harness / 红线） | **已关门**（2026-09-01：短文 v1.1.0 §2.6.1～2.6.5 整节落盘 + D-001 accepted） |
| C2 | **登记与继承**：§3.3 登记表 `rl` 行 + §1 端口分母 + §5 复核 + 修订史；判据 #5 共享约定继承确认（VP-028 不属 Redis 轨道保持）；`go.mod` redis 0 复核 | **已关门**（2026-09-01：`rl` 首条登记（026 义务闭环）· §1/§5/修订史 v1.1.0 · redis 0 · 零 Go 变更） |
| C3 | **审视与关门**：self A-001 + independent（grok build grok-4.6 · high A-002）合并响应；R3 关门、Root 台账回写 | **已关门**（A-001 self `pass` + A-002 grok independent `pass`（0 required · F-001 fixed · F-002/F-003 fixed-recording 落短文 §4 跟踪）；R3 关门 3/3 · Root progress 3/4；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R3 已关门）。

## 成功标准（方向级）

1. 判据 #4（Redis 接缝声明落盘）：供应商边界（端口不变）+ 原子窗口语义（INCR + EXPIRE）+ 连接管理约定写入短文 §2.6；不引入 Redis 客户端依赖（go.mod 复核算）。
2. 判据 #5（共享约定登记）：Redis 轨道约定（key 前缀 / 命名空间 / 连接管理 / 测试 harness）经 owner 文档继承登记；`rl` 命名空间段首条登记（§3.3）；VP-028 不属 Redis 轨道保持。
3. 继承义务闭环：workspace-026 关门的「命名空间登记义务 → VP-027 激活」触达并履行（证据 = §3.3 登记行 + 修订史 v1.1.0）。
4. 未越界：零 Go 代码变更；不改 Profile / Manifest / Charter；不消耗 RT-Q05 trigger。

## 信息就绪与未知项

无需新裁决：轨道条款（key 前缀 / 命名空间 / 连接管理 / harness / 变更流程）为 R1（VRev-059 V-F100）已冻结的 owner 约定，VP-027 激活即继承（短文 §1「继承即同意」）；I-027 四项全 verified。滑动窗口的 Redis 表达（ZSET vs 双桶近似）**触发立项时裁决**（短文 §2.6.2 已登记为未预裁项）。

## 父目标

- `GOAL-001-rate-limiter-port`（Root · 纲领 R3）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001 已定）：本目标为架构文档变更 + owner 登记（兼容性/元规则面，非代码实证）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent（项目级默认执行路径）。
- owner 文档修订依短文 §3.5 变更流程：任一 owner（VP-026/027）决策有效；本次 = VP-027 owner 决策（本目标 D-001）→ 修订史登记。