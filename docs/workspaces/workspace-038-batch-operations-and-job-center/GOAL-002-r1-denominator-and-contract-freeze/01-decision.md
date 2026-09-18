---
id: GOAL-002-r1-denominator-and-contract-freeze
doc: decision
status: active
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
| I-038-001 | required | Job 种类×作用域矩阵与读面缺口 | R1 C1、R2 | R1 | 侦察 + C1 矩阵 | open | — | `attachments/R1-recon-I-038-001-job-kinds-and-scopes.md` |
| I-038-002 | required | 批量异步契约形态与协议面影响 | R1 C2、R3 | R1 | 侦察 + 用户 P-004 裁决 | open | — | `attachments/R1-recon-I-038-002-batch-async-contract.md` |
| I-038-003 | required | 首波批量/长操作分母 | R1 C3、R3 | R1 | 侦察 + 用户 P-004 裁决 | open | — | `attachments/R1-recon-I-038-003-batch-operation-inventory.md` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 暂无（R1 冻结决策待用户 P-004 裁决后落盘） | — | — |

## 待裁决（P-004 · 禁止静默自动裁）

R1 的 C2 与 C3 属**方案选型**，须经用户裁决后落盘：

1. **C2 契约形态**：选项 A（扩展 ADR-0022 异步变体，触碰 pinned schema）/ B（另立本地模块自有异步契约）/ C（混合：本地扩展位承载）。
2. **C3 首波分母**：异步首波纳入哪些操作（导出 / 导入 / `purge-all` / 其他），以及是否需要在首波内**新建**批量操作（当前生产页面批量 UI 分母为 0）。
3. **C1 作用域口径**：管理作用域（跨 actor）Job 读面的权限模型。

裁决结论与未选方案将写入 `01-decision/D-001-r1-contract-and-denominator-freeze.md`。

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
