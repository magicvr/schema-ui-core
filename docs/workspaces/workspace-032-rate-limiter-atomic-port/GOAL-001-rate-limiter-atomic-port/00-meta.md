---
id: GOAL-001-rate-limiter-atomic-port
title: 限流器端口原子化
status: done
parent: null
created: 2026-09-03
updated: 2026-09-04
version: 0.4.0
progress: 3/3
plan_refs:
  - VP-032-rate-limiter-atomic-port
primary_plan: VP-032-rate-limiter-atomic-port
serves_summary: 架构分支 · 承接 VP-027 residual R-007：kernel.RateLimiter 新增原子 AllowRecord 与令牌化 Reserve/Cancel，迁移冻结 14 处 Allow→Record 使用点（4 立即消费 + 10 失败预算），消除 TOCTOU；内存供应商实现；Redis 仍 RT-Q05 trigger-gated。
---

# GOAL-001 · 限流器端口原子化

## 概述

承接 [VP-032-rate-limiter-atomic-port](../../vision/plans/VP-032-rate-limiter-atomic-port.md)（active v0.2.0 · [VRev-073](../../vision/reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) self `pass` · 架构类 freshness PASS `42036a3c`→`b1c03acd`）：消除 `kernel.RateLimiter` Allow/Record 两调用之间的 TOCTOU。**对象面**：端口新增 `AllowRecord(key, now) bool` + 内存供应商单锁实现 + 冻结 14 处生产使用点迁移 + 并发穿透回归。**红线（激活即生效）**：不重开 VP-027；不实现 Redis / **不消耗 RT-Q05 trigger**；不改 Profile 默认集 / Manifest（VP-008 `go`）；`Allow`/`Record` 保留兼容。

## 成功标准（对应 VP-032 五条方向级退出判据）

- [x] 判据 #1（原子性）：`AllowRecord`/`Reserve` 在并发下 check+record 原子，无穿透窗口（并发预算 + 无穿透回归测试全绿）
- [x] 判据 #2（行为等价）：冻结 14 处使用点全部迁移；立即消费路径单请求等价；失败预算路径按 GOAL-003 D-002 §3 逐路径语义冻结（每种结果 = 旧计数行为，非计数 `Cancel` 只回滚当次、保留历史）
- [x] 判据 #3（兼容）：`Allow`/`Record` 保留；文档标注 `AllowRecord` 为推荐路径
- [x] 判据 #4（边界保持）：未重开 VP-027；未实现 Redis；未改 Profile 默认集
- [x] 判据 #5（审计闭合）：开放 required finding = 0（全部已合法闭合）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。`progress` = 已完成纲领阶段 / 3。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同落盘（GOAL-002）：D-002 冻结 + kernel.AllowRecord + Memory 单锁实现 + 合同级测试。I-032-001/002 已由 VRev-073 冻结 | **已关门**（2026-09-03 · GOAL-002 合同与端口落地 + A-003 关门） |
| R2 | 14 处迁移 + handler 回归（判据 2/3）：按立即消费 / 失败预算两口径迁生产调用点；handler 既有限流测试仍绿 | **已关门**（GOAL-003：14 处全迁；失败预算初版被 A-002 证伪 → D-002 令牌化 Reserve/Cancel 修复 + 混合历史回归全绿；A-003 self + A-004 grok independent 复审 pass · 0 required） |
| R3 | 证据与关门（判据 4/5；依赖 R1–R2）：证据矩阵 / 越界核账 / 审计闭合 | **已关门**（2026-09-04 · E-004 证据矩阵 · A-001 self + A-002 grok independent 双 pass · 开放 required = 0） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-032-001 | required | `AllowRecord` 精确签名与返回值语义（bool 是否足够，是否需返回剩余额度）。 | 方案冻结 + 退出判据 1 | R1 | `/vision` 激活冻结（VRev-073） | **verified** | — | 2026-09-03：`AllowRecord(key string, now time.Time) bool`；bool 足够；不返回剩余额度；`RetryAfterSeconds` 独立（VRev-073） |
| I-032-002 | required | 是否所有使用点都应迁移（Clear-on-success 是否需要原子变体）。 | 方案冻结 + 退出判据 2 | R1 | `/vision` 激活冻结（VRev-073）；**2026-09-04 回流重审（GOAL-003 A-002）** | **revised** | — | 2026-09-03：14 处生产 Allow→Record 全迁；Clear 无需原子变体；立即消费 vs 失败预算两口径（VRev-073）。**2026-09-04 修正**：键级 Clear 无法回滚当次占槽 → 失败预算改为令牌化 Reserve/Cancel（GOAL-003 D-002 · I-032-003） |
| I-032-003 | required | 令牌化保留契约（Reserve/Cancel 签名与逐路径语义冻结） | 实施门禁 + 退出判据 2/5 | R2 | 用户裁决（2026-09-04 · 方案 A）→ GOAL-003 D-002 冻结 → 实施验证 | **verified** | — | `Reserve(key, now) (token uint64, ok bool)` + `Cancel(key, token)`；10 处失败预算逐路径冻结（GOAL-003 D-002 §3）；GOAL-003 A-004 独立核对一致 |

全工作区无开放 required 信息门禁（I-032-002 已 revised，结论由 I-032-003 承接）。

R1 无开放 required 信息门禁（两项均已 vision 冻结）。R1 工作是把冻结写入 kernel 合同与测试，不是再裁决签名。

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- **开区（2026-09-03 · 用户指令）**：VP-032 `planned → active` v0.2.0（VRev-073 self `pass` 0 required · 架构类 freshness PASS `42036a3c`→`b1c03acd` · 不暂挂 `go`）；lead `workspace-032-rate-limiter-atomic-port`。
- **关门（2026-09-04）**：R1–R3 全链条完成；Root `done` 3/3。关门双审 A-001 self `pass` + A-002 grok build independent `pass`（0 required）。VP-032 文案承接（§首波冻结/判据 #2 的「失败预算 = 入口 AllowRecord + Clear」表述被 GOAL-003 D-002 令牌化取代；判据意图达成）留 `/vision` VP-032 关门/VRev 登记（E-004 §4 · A-002 R-002），不构成本区门禁。
- 审计模式（D-001）：阶段关门 default self；实证门禁（R3 证据 / 关门）按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点见 D-001：消费候选 = HEAD `b1c03acd`。
- V-F117 recommended：VP-030 Root 已 done 但 VP 仍 `active`；另轮 `/vision` 关门，不构成本区实施门禁。
