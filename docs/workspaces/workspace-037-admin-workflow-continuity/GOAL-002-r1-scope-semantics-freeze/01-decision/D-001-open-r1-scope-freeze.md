---
id: D-001-open-r1-scope-freeze
doc: decision-entry
status: accepted
parent: GOAL-002-r1-scope-semantics-freeze
created: 2026-09-16
updated: 2026-09-16
version: 0.1.0
source: supervisor
---

# D-001 · 建立 R1 事实与语义冻结子目标

## 决定

在 workspace-037 内建立 `GOAL-002-r1-scope-semantics-freeze`，承载 R1 的分母盘点、状态/反馈场景矩阵和方案冻结；R2、R3、R4 继续保持 pending，直到 I-037-001～004 关闭。

## 理由

R1 具有独立的证据交付、门禁和审计范围。把它从 Root 拆出可以保持 Root 路线图的阶段顺序，同时避免在列表分母、用户隔离或离开保护尚未明确时直接改 runtime。

## 不包含的决定

本条不决定 Saved View 的持久化介质、跨设备语义、列配置格式或 dirty-state 的最终文案；这些属于 I-037-002～004 的待确认方案决策。
