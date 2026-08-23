---
id: D-001
doc: decision-entry
goal: GOAL-001-observability
status: accepted
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

# D-001 · 开区 scaffold 与 A4 纲领路线图

## 背景

用户确认：对 VP-015 做 self Review；通过后激活并开区。未另写工作区 slug。VRev-033 self `pass`（0 required）。A4 退出分母已由 VP-015 v0.1.0 / VR-034 冻结，激活时 editorial 收口退出 4 与 I-015-003。

工作区 slug 按 VP-013 / VP-014 惯例取 `workspace-015-observability`（`workspace-NNN` + VP slug）。Root slug 取 `GOAL-001-observability`，与 VP / 工作区后缀对齐。用户确认的是「开区」动作；slug 为惯例推导并在此留痕，不是另一次口头点名。

## 决策

1. 激活 `VP-015-observability`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-015-observability`；Root = `GOAL-001-observability`。
3. 纲领路线图 R1～R5：导出合同冻结 → 指标 scrape → OTel traces → 与 request-id 关联 → 双路径证据。
4. 配置面：缺省无 Prometheus / collector / Jaeger；导出为显式配置的生产/验收路径。不改 Compose 默认依赖。
5. 开区审计模式 **none**（可逆文档 scaffold）。R1 合同冻结起按内核门禁走 **self**；指标/tracing 生产路径实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
6. 本回合**不**创建 R1 子目标、**不**改 `apps/api` 代码。

## 架构类 freshness（V-F065）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是架构 A4，按自身激活门闩做轻量复核：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10） |
| 现行 HEAD | `323c00a` |
| VP-009 / VP-010 | 无开放中高危暂挂宣称 |
| Vision open required | 0 |
| F-007 | 上传授权深度仍 deferred；本 VP 不扩张授权 scope |
| 是否消费业务解锁 | **否**。不改 Profile / 模块矩阵 / Manifest 为意图 |
| 结果 | **PASS**。不暂挂 `go` |

`consumer_vp` = VP-015；`last_freshness_review_at` = 2026-08-21；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-014。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-005 满足 V-F064；freshness 表满足 V-F065。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求激活并开区。
- 重开 workspace-014：VP-014 已关门且默认不接新区。
- 一开区就引入 OTel SDK / Prometheus：R1 合同未冻结。
- 默认改为必须有收集器：与 A4 内嵌默认冲突。
- 把本 VP 当业务 VP 走 VP-011 式完整 freshness 矩阵：解锁 scope 不匹配。
- 等待用户另写 slug 再开区：用户本轮已指令「开区」；惯例与 VP-013/014 同构，记录推导以免静默。
