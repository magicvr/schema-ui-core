---
id: GOAL-012-w12-multi-instance-rate-limiting
doc: decision
status: done
parent: GOAL-001-production-hardening
created: 2026-08-26
updated: 2026-08-26
version: 0.2.0
---

# 决策记录 · GOAL-012

## 信息需求与阶段门禁

> 本文件是稳定索引；信息台账与 `00-meta.md` 同号同状态镜像。`accepted-residual` 必须指向用户书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 产品部署拓扑意图：是否官方支持多实例水平扩展部署？ | S2 方案冻结 | S2 前 | 用户裁决 | **verified**（D-002 §1） | 复审触发 =「多实例部署形态出现」 | 维持单实例官方边界 |
| I-002 | required | 共享限流状态载体选型：内核 Store 端口新表 vs 进程外依赖（Redis 等） | S2 方案冻结 / S3 实施 | S2 前 | 用户裁决 | **verified**（D-002 §2） | 预登记 ≠ 实施承诺；A3 触发时正式冻结 | 预登记方向 = Redis 等进程外依赖（用户裁决） |
| I-003 | non-blocking | login/recovery 桶键空间与预算统一性、多实例 Retry-After 语义、窗口参数可配置性 | S2 条款细化 | S2 | 随 D-002 一并裁决 | **closed**（D-002 §3） | 多实例语义归 A3 触发时冻结 | 单实例边界下语义保持现状 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | W12 立项：承接 workspace-019 E-009 §F-002 登记项（评估先行路线图） | accepted | [D-001-w12-intake-and-roadmap.md](01-decision/D-001-w12-intake-and-roadmap.md) |
| D-002 | 2026-08-26 | S2 方案冻结：维持单实例官方边界 · 载体预登记 Redis 方向 · 零码收官 | accepted | [D-002-w12-s2-freeze-single-instance.md](01-decision/D-002-w12-s2-freeze-single-instance.md) |
