---
id: GOAL-004-r3-async-batch-operation
doc: decision
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-011 | required | 导出数据面分母与权限口径 | C1 | C1 前 | 读既有导出实现与列集 | open | — | `internal/handler/export.go` |
| I-038-012 | required | 自定义组件能否读到表格选择集 | C2 | C2 前 | 读 render selection 与 custom context | open | — | R3 侦察 |
| I-038-013 | non-blocking | 导出文件名/格式一致性 | C2/R4 | R4 前 | 对照既有 CSV 约定 | open | — | `export.go` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 暂无（R3 方案待侦察与用户裁决后落盘） | — | — |

## 约束输入（来自 R1/R2 冻结，非本目标可改）

**R1 `GOAL-002/01-decision/D-001-…`**：

- **C2 = 方案 B**：异步契约归 `admin.jobs` 本地端点；ADR-0022 同步语义**逐字冻结**；协议 pin 零改动。
- **C3 = 首波仅 1 条**：「新建批量导出所选」；同步 `batch-delete` 保持（VP-038 显式非目标）。
- **K-1**：异步提交**不得**走 `batchMapping` 成功路径（202 体被丢弃）。
- **K-2**：202 + job 投影（与 wallet reconcile 同形）。
- **K-5**：新 kind 须在 runner 启动前注册，kind 不可重复。
- **K-6**：进度须由新 handler 细粒度上报（`reporter.Progress(n)` 语义为 0..99，终态由 `CompleteWithCommit` 置 100）。
- **K-7**：无需新建迁移（`jobs` 表已存在）。

**R2 `GOAL-003/01-decision/D-001-…`**：

- **§3 O-1**：前端触发 = **自定义组件**（`registerCustomComponent`），形态照 `monitoring-auto-refresh.tsx`。
- **§3 O-2**：页面**不声明** `actions.batch.request`；用既有 `table.selection` / `actions.page.trigger` 标记可覆盖的 capability。
- **§4**：`jobs.write` 归本阶段声明（`PolicyAdmin`）；路由前缀 `/api/jobs`。
- **§2**：结果地址由共享 `jobs.ResultURL(basePath, id)` 派生。

## 待冻结（本目标方案项）

1. **导出数据面分母**（`I-038-011`）：哪些资源可批量导出；列集与敏感字段口径；权限键。
2. **前端选择集获取路径**（`I-038-012`）：自定义组件如何取得当前选中行。
3. **触发入口位置**：列表页 body 内的 custom 节点 vs 其他可行位置。
4. **进度语义**：按行数上报的分档与节流。
5. **Job kind 命名与导出结果形状**。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
