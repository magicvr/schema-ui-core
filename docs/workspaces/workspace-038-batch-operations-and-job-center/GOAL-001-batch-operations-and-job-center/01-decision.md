---
id: GOAL-001-batch-operations-and-job-center
doc: decision
status: active
parent: null
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md`（Root 层）维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-001 | required | Job 种类分母与可见作用域 | R1 范围冻结、R2 读面 | R1 | 扫描 `internal/jobs` 与各模块 Job 消费点，形成种类×作用域矩阵 | open | — | 待确认 |
| I-038-002 | required | 批量异步契约与协议面影响 | R1 方案冻结、R3 实施 | R1 | 对照 `apps/web/src/protocol` 与上游 v2.9.0 契约 | open | — | 待确认 |
| I-038-003 | required | 首波批量操作分母（保持同步 / 改异步 / 不进首波） | R1 范围冻结、R3 实施 | R1 | 盘点 `resources.go` 批量面、`data-transfer`、wallet reconcile | open | — | 待确认（承接 `V-F126`） |
| I-038-004 | required | Profile / 模块矩阵边界 | 激活、R1、VP-008 `go` | 激活前 | 读 `kernel/profile.go` 模块矩阵与 Profile 集合 | **verified** | 2026-09-19 用户 P-004 裁决 | D-001 |
| I-038-005 | required | 激活前 Admin 类 freshness | 激活与开区 | 激活前 | 五域 freshness review | **verified** | 2026-09-19 已完成 | D-001；VRev-099 |
| I-038-006 | non-blocking | 历史作业保留与清理策略 | 后续运维波次 | 关门后或容量触发 | 出现真实需求时由 `/vision` 复核 | deferred | 理由：首波聚焦可见性；责任人：`/vision`；触发：容量/保留期需求 | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | 工作区与 Root 建立：VP-038 激活落盘、`I-038-004` 裁决与 freshness 记录 | accepted | `01-decision/D-001-workspace-root-establishment.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
