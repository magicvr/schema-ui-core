---
id: GOAL-002-r1-denominator-and-contract-freeze
doc: decision
status: done
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 与本目标 `01-decision/D-NNN-*.md` 维护；`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-001 | required | Job 种类×作用域矩阵与读面缺口 | R1 C1、R2 | R1 | 侦察 + C1 矩阵 | **verified**（2026-09-19） | — | `attachments/R1-recon-I-038-001-job-kinds-and-scopes.md`；`attachments/r1-job-kind-scope-matrix.md`；D-001 §2 |
| I-038-002 | required | 批量异步契约形态与协议面影响 | R1 C2、R3 | R1 | 侦察 + 用户 P-004 裁决 | **verified**（2026-09-19 用户裁决方案 B） | — | `attachments/R1-recon-I-038-002-batch-async-contract.md`；D-001 §1 |
| I-038-003 | required | 首波批量/长操作分母 | R1 C3、R3 | R1 | 侦察 + 用户 P-004 裁决 | **verified**（2026-09-19 用户裁决首波 = 新建批量导出所选） | — | `attachments/R1-recon-I-038-003-batch-operation-inventory.md`；`attachments/r1-first-wave-denominator-matrix.md`；D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | R1 分母与契约冻结（C1/C2/C3；用户 P-004 裁决） | accepted | `01-decision/D-001-r1-contract-and-denominator-freeze.md` |

## P-004 裁决记录（2026-09-19）

| 决策点 | 用户裁决 | 备选（未选） |
|--------|---------|-------------|
| C2 契约形态（`I-038-002`） | **方案 B** · 另立本地模块自有异步契约；ADR-0022 同步语义完全冻结；协议 pin 零改动 | A（扩展 ADR-0022 异步变体，触碰 pinned schema）；C（本地扩展位挂批量工具栏） |
| C3 首波分母（`I-038-003`） | **仅「新建批量导出所选」**（1 条） | 回收站 `purge-all` 改异步；CSV 导入改异步；既有同步 `batch-delete` 改异步 |
| C1 作用域模型（`I-038-001`） | **管理作用域 + 新增 `jobs.read` 权限**（`PolicyAdmin`） | 复用既有权限；仅 actor 作用域 |

**未定项（不冒充已裁决，见 D-001 §1.3）**：O-1 前端触发机制；O-2 capability 声明口径；O-3 管理列表索引决策。三项均归 R2 方案冻结。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
