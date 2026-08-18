---
id: A-003-r5-a002-response
goal: GOAL-006-r5-maintenance-read-only-gate
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# A-003 · R5 A-002 response

## 响应头

| 项 | 值 |
|----|----|
| source | self |
| scope | A-002 F-001/F-002；D-003；I-002/I-004 |
| verdict | pass |
| required findings | 0 |

## 原 finding 保留

- A-002 F-001：required/high，`form.controls.readonly` 不能作为全局 read-only 的 disabled capability；原文与证据保留于 [A-002-r5-s0-design-independent.md](A-002-r5-s0-design-independent.md)。
- A-002 F-002：required/med，degraded/read-only 写拒绝不应指向 bootstrap recovery；原文与证据保留于同一 A-002。

## 响应与修订

1. D-003 §Host 与 status 投影已移除 `disabledCapabilities` 与 `host.readOnly` 依赖；degraded/read-only 均仅使用既有 bootstrap `mode: degraded`，原始区别由 system-monitoring `availabilityMode` 表达。该修复不改变上游 protocol pin，也不会剔除 `form.controls.readonly`，因此关闭 F-001 的候选 disposition 为 `fixed`。
2. D-003 §运行态与写门禁明确 degraded/read-only 是应用内 `503` catalog error，只有 maintenance bootstrap 才触发 Host `MAINTENANCE`；不承诺 Retry-After，且客户端按 error code 分流。该修复关闭 F-002 的候选 disposition 为 `fixed`。
3. D-003 同时吸收 F-003～F-006 的实施门要求：统一使用 503 避免 423 collision，显式空 `RUNTIME_MODE` fail closed，S2 精确匹配注册方法和 allowlist。

## 结论

F-001/F-002 已完成可核对修订；A-004 independent 已复核为 pass，I-002～I-004 已 verified，S1 获准开始。
