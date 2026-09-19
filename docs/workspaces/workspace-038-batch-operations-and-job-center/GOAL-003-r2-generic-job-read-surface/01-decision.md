---
id: GOAL-003-r2-generic-job-read-surface
doc: decision
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-003

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-007 | required | 管理列表索引决策（O-3） | C2 | C2 前 | 索引矩阵 + 查询形状静态判定 | open | — | R1 矩阵 §4 |
| I-038-008 | required | 结果 URL 泛化口径 | C3 | C3 前 | 读 wallet 结果路由与映射 | open | — | `handler/wallet.go:989-1006` |
| I-038-009 | required | 前端触发机制与 capability 声明口径（O-1/O-2） | C4、R3/R4 | R2 方案冻结前 | 对照先例与守卫 | open | — | R1 `D-001` §1.3 |
| I-038-010 | non-blocking | 导航分组与 i18n 键位 | R4 | R4 前 | 读分组测试与 fragment | open | — | R1 侦察 §6 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 暂无（R2 方案待侦察与用户裁决后落盘） | — | — |

## 约束输入（来自 R1 冻结，非本目标可改）

`GOAL-002/01-decision/D-001-r1-contract-and-denominator-freeze.md`：

- **C1**：管理作用域 + 新增 `jobs.read` 权限（`PolicyAdmin`）；既有 `GetForActor` actor 隔离语义与其冻结测试**不动**，通用读面走新方法 + 新路由。
- **C2**：方案 B —— 异步契约归本地模块自有端点；ADR-0022 同步语义逐字冻结；协议 pin 零改动。
- **C3**：首波 = 仅「新建批量导出所选」1 条。
- **K-1～K-7**：异步提交不得走 `batchMapping` 成功路径；202 + job 投影；409/410 结果语义；写操作须有写权限门；新 kind 须在 runner 启动前注册且不可重复；进度须由新 handler 细粒度上报；无需新建 `jobs` 表。
- **R-1**：`jobRuntime.enabled` 仅在 `admin.wallet` 下置 true —— **R2 必须显式处理**（否则含 `admin.jobs` 但不含 `admin.wallet` 的 Profile 下 runner 不启动）。

## 待冻结（本目标方案项）

1. **索引决策**（`I-038-007`）：是否新增迁移版本。
2. **结果 URL 泛化**（`I-038-008`）。
3. **前端触发机制与 capability 声明口径**（`I-038-009`）。
4. **写权限键名**（`jobs.write` 或其他）与策略（`PolicyAdmin` vs `PolicyAdminEditor`）。
5. **`jobRuntime.enabled` 的 R2 处理方式**（R-1）。
6. **列表默认排序与排序白名单**（须与索引决策一致）。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
