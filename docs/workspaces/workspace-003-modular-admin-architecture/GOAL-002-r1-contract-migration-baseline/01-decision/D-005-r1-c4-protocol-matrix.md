---
id: D-005
title: 冻结 C4 协议继承与模块候选三态矩阵
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
---

# D-005 · C4 协议继承与模块候选三态矩阵

## 决定

1. R1 继承 Q2 canonical `I-PROTO-001 v0.1.3`：`include` 为 D-NODE/D-EXPR/D-DATA/D-PERM/D-APP/D-VER/D-VAL，`include-partial` 为 D-COMP/D-ACT/D-TABLE/D-FORM，`exclude` 为 D-UPLOAD。
2. 以 [r1-c4-protocol-matrix.md](../attachments/r1-c4-protocol-matrix.md) 作为候选模块到 protocol domain、fixture suite 与 Profile candidate 的可追踪矩阵；该矩阵不代表模块已实现，也不冻结 I-004 的精确 Profile 集或 precedence。
3. 维持 D-EXPR 的 `$context`/visible-when 子集、D-ACT 的非批量子集、D-TABLE 的排序/搜索子集、D-FORM/D-COMP 的固定边界；D-UPLOAD 继续整域排除。
4. 新增 domain、扩大 partial、改变 D-UPLOAD 排除或引入新上游协议版本时，必须追加新决策、递增覆盖表版本、完成验证并更新 affected required information gate；不得静默改写 v0.1.3。

## 修正

此前 C1 摘要漏列 D-EXPR 与 D-VER；本决定以冻结表的 7 个 include domain 为准，已在 C1 evidence 中修正。

## 约束

本决定只冻结协议继承范围和变更门槛，不修改 Q2 覆盖表、不读取 workspace-001 过程状态、不将 fixture/test presence 解释为 R1 模块化实现完成。
