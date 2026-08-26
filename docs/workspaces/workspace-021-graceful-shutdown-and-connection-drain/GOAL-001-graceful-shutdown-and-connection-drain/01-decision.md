---
id: GOAL-001-graceful-shutdown-and-connection-drain
doc: decision
status: active
parent: null
created: 2026-08-27
updated: 2026-08-27
version: 0.1.0
---

# 决策记录 · GOAL-001 优雅停机 / 连接排空合同

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 停机时运行中 Job 语义：等完成 vs 中断标记重跑 | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：中断标记重跑（GOAL-002 D-001） |
| I-002 | required | grace period / 超时默认值与可配置键 | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：默认 10s + `http.shutdown_timeout`；非法值 fail-closed（GOAL-002 D-001） |
| I-003 | required | Store 排空与迁移窗口重叠时的停机语义 | 方案冻结 / 实施 | R2 | 用户裁决 | **verified** | — | 2026-08-27 用户裁决：fail-closed 启动期校验，无运行时迁移窗口（GOAL-002 D-001） |
| I-004 | non-blocking | 停机是否需日志 / 指标断言（消费 VP-015） | 验收 | R3 | lead 提案 | **verified** | — | lead 口径：结构化日志断言；指标不进分母（GOAL-002 D-001） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-27 | 激活与开区：lead 绑定、纲领 R1～R3、信息门禁登记、架构类 freshness 留痕 | accepted | `01-decision/D-001-workspace-root-establishment.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。