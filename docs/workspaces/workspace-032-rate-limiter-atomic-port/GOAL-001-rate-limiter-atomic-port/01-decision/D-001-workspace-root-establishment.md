---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-rate-limiter-atomic-port
date: 2026-09-03
status: active
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户 2026-09-03 指令「/vision 走流程激活 VP-032（架构原子限流），然后交 /govern 开设工作区」；slug 按 VP-013～030 惯例且 VP-032 正文预登记（`workspace-032-rate-limiter-atomic-port` / `GOAL-001-rate-limiter-atomic-port`）。激活门禁已满足：VRev-073 self `pass`（0 required；I-032-001/002 冻结）+ 架构类轻量 freshness PASS（`42036a3c` → `b1c03acd`，不暂挂 `go`）。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-032-rate-limiter-atomic-port`（canonical `docs/workspaces/workspace-032-rate-limiter-atomic-port/`） |
| Root | `GOAL-001-rate-limiter-atomic-port`（`parent: null`；primary_plan = `VP-032-rate-limiter-atomic-port`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 合同落盘 → R2 内存实现+14 处迁移+并发回归 → R3 证据与关门（串行；阶段内可并行子目标） |
| 审计模式 | 阶段关门 default **self**；R3 证据/关门实证门禁可按需 **independent**（grok build 先例 · 项目级默认执行路径） |
| 红线 | 不重开 VP-027；不预制 Redis（不引入客户端依赖 / 不消耗 RT-Q05 trigger）；不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`）；`Allow`/`Record` 保留兼容；GOAL-014 分层锁定与 recyclebin 领域 `Record` 显式排除 |

## 继承的激活冻结（VRev-073 · 不可在 R1 重开裁决，除非用户书面改分母）

| ID | 冻结结论 |
|----|----------|
| I-032-001 | `AllowRecord(key string, now time.Time) bool`；bool 足够；不返回剩余额度；`RetryAfterSeconds` 独立 |
| I-032-002 | 14 处生产 Allow→Record 全迁；Clear 无需原子变体；立即消费 vs 失败预算两口径 |

## freshness 三字段（VRev-073 · 先例执行惯例）

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-032-rate-limiter-atomic-port`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-03（`42036a3c` → `b1c03acd` · 架构类轻量 PASS · 协议 pin / 依赖锁 / Profile 默认集 / provenance 零变更；区间 = VP-030 已审结目） |
| next_freshness_review_trigger | 首个 C 端业务域 VP 激活（H-002 发现机制）或 多实例部署评估 |

## 未选方案

- 不在本波关闭 VP-030（V-F117 recommended；alignment §7 须用户确认；不阻断本区开区）。
- 不把本意图做成 workspace-027 子目标（VP-027 已 closed，禁止重开）。
- 不在本波实现 Redis 供应商（RT-Q05 触发条件未成立）。
- 不返回剩余额度（超出 residual 范围）。
