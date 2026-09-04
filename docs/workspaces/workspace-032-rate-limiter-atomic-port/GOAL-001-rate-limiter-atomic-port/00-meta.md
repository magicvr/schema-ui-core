---
id: GOAL-001-rate-limiter-atomic-port
title: 限流器端口原子化
status: active
parent: null
created: 2026-09-03
updated: 2026-09-04
version: 0.3.1
progress: 1/3
plan_refs:
  - VP-032-rate-limiter-atomic-port
primary_plan: VP-032-rate-limiter-atomic-port
serves_summary: 架构分支 · 承接 VP-027 residual R-007：kernel.RateLimiter 新增原子 AllowRecord 并迁移冻结 14 处 Allow→Record 使用点，消除 TOCTOU；内存供应商实现；Redis 仍 RT-Q05 trigger-gated。
---

# GOAL-001 · 限流器端口原子化

## 概述

承接 [VP-032-rate-limiter-atomic-port](../../vision/plans/VP-032-rate-limiter-atomic-port.md)（active v0.2.0 · [VRev-073](../../vision/reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) self `pass` · 架构类 freshness PASS `42036a3c`→`b1c03acd`）：消除 `kernel.RateLimiter` Allow/Record 两调用之间的 TOCTOU。**对象面**：端口新增 `AllowRecord(key, now) bool` + 内存供应商单锁实现 + 冻结 14 处生产使用点迁移 + 并发穿透回归。**红线（激活即生效）**：不重开 VP-027；不实现 Redis / **不消耗 RT-Q05 trigger**；不改 Profile 默认集 / Manifest（VP-008 `go`）；`Allow`/`Record` 保留兼容。

## 成功标准（对应 VP-032 五条方向级退出判据）

- [ ] 判据 #1（原子性）：`AllowRecord` 在并发下 check+record 原子，无穿透窗口（有并发回归测试）
- [ ] 判据 #2（行为等价）：冻结 14 处使用点全部迁移；立即消费路径单请求等价；失败预算路径 `Clear` 后净状态等价
- [ ] 判据 #3（兼容）：`Allow`/`Record` 保留；文档标注 `AllowRecord` 为推荐路径
- [ ] 判据 #4（边界保持）：未重开 VP-027；未实现 Redis；未改 Profile 默认集
- [ ] 判据 #5（审计闭合）：开放 required finding = 0（或已合法闭合）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。`progress` = 已完成纲领阶段 / 3。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同落盘（GOAL-002）：D-002 冻结 + kernel.AllowRecord + Memory 单锁实现 + 合同级测试。I-032-001/002 已由 VRev-073 冻结 | **已关门**（2026-09-03 · GOAL-002 合同与端口落地 + A-003 关门） |
| R2 | 14 处迁移 + handler 回归（判据 2/3）：按立即消费 / 失败预算两口径迁生产调用点；handler 既有限流测试仍绿 | **进行中**（GOAL-003：14 处已全迁；失败预算初版被 A-002 证伪 → D-002 令牌化修复完成、回归全绿，待复审关门） |
| R3 | 证据与关门（判据 4/5；依赖 R1–R2）：证据矩阵 / 越界核账 / 审计闭合 | 待 R2 关门 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-032-001 | required | `AllowRecord` 精确签名与返回值语义（bool 是否足够，是否需返回剩余额度）。 | 方案冻结 + 退出判据 1 | R1 | `/vision` 激活冻结（VRev-073） | **verified** | — | 2026-09-03：`AllowRecord(key string, now time.Time) bool`；bool 足够；不返回剩余额度；`RetryAfterSeconds` 独立（VRev-073） |
| I-032-002 | required | 是否所有使用点都应迁移（Clear-on-success 是否需要原子变体）。 | 方案冻结 + 退出判据 2 | R1 | `/vision` 激活冻结（VRev-073）；**2026-09-04 回流重审（GOAL-003 A-002）** | **revised** | — | 2026-09-03：14 处生产 Allow→Record 全迁；Clear 无需原子变体；立即消费 vs 失败预算两口径（VRev-073）。**2026-09-04 修正**：键级 Clear 无法回滚当次占槽 → 失败预算改为令牌化 Reserve/Cancel（GOAL-003 D-002 · I-032-003） |

R1 无开放 required 信息门禁（两项均已 vision 冻结）。R1 工作是把冻结写入 kernel 合同与测试，不是再裁决签名。

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- **开区（2026-09-03 · 用户指令）**：VP-032 `planned → active` v0.2.0（VRev-073 self `pass` 0 required · 架构类 freshness PASS `42036a3c`→`b1c03acd` · 不暂挂 `go`）；lead `workspace-032-rate-limiter-atomic-port`。
- 审计模式（D-001）：阶段关门 default self；实证门禁（R3 证据 / 关门）按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点见 D-001：消费候选 = HEAD `b1c03acd`。
- V-F117 recommended：VP-030 Root 已 done 但 VP 仍 `active`；另轮 `/vision` 关门，不构成本区实施门禁。
