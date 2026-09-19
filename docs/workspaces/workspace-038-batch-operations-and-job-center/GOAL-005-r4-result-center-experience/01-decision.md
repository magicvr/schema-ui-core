---
id: GOAL-005-r4-result-center-experience
doc: decision
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-014 | required | 结果中心结构选型（既有 jobs 页收敛 vs 独立页） | C2 | C2 前 | 读页面结构/行操作/导航约定 | open | — | R4 侦察 |
| I-038-015 | required | 管理作用域取消/重试的合同口径 | C1 | C1 前 | 读 repository actor 实现与 runner 取消路径 | open | — | `repository.go` |
| I-038-016 | non-blocking | 结果过期后的呈现与下载失效语义 | C2 | R4 前 | 对照既有错误呈现约定 | open | — | R2 `jobs.go` |
| I-038-013 | non-blocking | 导出文件名/格式与两入口 UI 文案一致性（承接自 R3） | C3 | R4 内 | 对照既有导出约定 | open | — | R3 `D-001` §1.1 |
| I-038-010 | non-blocking | 导航分组与 i18n 键位（承接自 R2） | C3 | R4 内 | 读分组测试与 fragment | open | — | R2 `D-001` §4 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 暂无（R4 方案待侦察与用户裁决后落盘） | — | — |

## 约束输入（来自 R1～R3 冻结，非本目标可改）

- **R1 `D-001`**：契约归本地（方案 B）；ADR-0022 同步语义冻结；首波 1 条；`jobs.read` 管理作用域；K-3 结果三段语义（409/410/终态）。
- **R2 `D-001`**：O-1 自定义组件触发；O-2 不声明 `actions.batch.request`；路由前缀 `/api/jobs`；结果地址由 `jobs.ResultURL` 派生；`jobs.write` 归 R3。
- **R3 `D-001`**：`jobs.batch-export` 数据面（users/roles 同分母、上限 500、结果 JSON 含 csv）；`jobs.write` + `data.export` 双重门禁；进度语义；前端不 `reloadList`。

## 待冻结（本目标方案项）

1. **结构选型**（`I-038-014`）：结果中心落在既有 `jobs` 页还是独立页。→ **须询问用户**
2. **取消/重试合同口径**（`I-038-015`）：可取消/可重试状态集合、取消进行中的语义、重试上限。
3. **过期与下载失效的 UI 呈现**（`I-038-016`）。
4. **i18n/主题/可访问性收敛范围**（`I-038-013`、`I-038-010`）。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
