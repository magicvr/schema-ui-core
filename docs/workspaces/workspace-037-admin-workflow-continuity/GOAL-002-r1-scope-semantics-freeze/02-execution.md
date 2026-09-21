---
id: GOAL-002-r1-scope-semantics-freeze
doc: execution
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-17
version: 0.5.0
---

# 执行台账 · GOAL-002 R1

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-16 | 建立 R1 子目标并完成初始扫描 | recorded | [E-001-initial-r1-scan.md](02-execution/E-001-initial-r1-scan.md) |
| E-002 | 2026-09-17 | 形成 R1 分母与状态/反馈矩阵 | recorded | [E-002-r1-matrix-record.md](02-execution/E-002-r1-matrix-record.md) |
| E-003 | 2026-09-17 | 记录用户确认 Saved View 使用 localStorage | recorded | [E-003-user-saved-view-choice.md](02-execution/E-003-user-saved-view-choice.md) |
| E-004 | 2026-09-17 | 冻结 dirty-state 与反馈/恢复实现语义 | recorded | [E-004-r1-semantics-freeze.md](02-execution/E-004-r1-semantics-freeze.md) |

## 当前事实

- 已建立本目标五件套、三个 ledger 目录和 `attachments/`。
- 本轮 R1 文档响应基于 Git HEAD `92cf8582`；另有既有未跟踪 `.claude/settings.local.json`，未触碰。
- 已从 `apps/api/modules/**/schema/*.json` 盘点出当前全部 table/form 节点；精确矩阵见附件。
- 已核对 4 个 profile 的注册/可发现/隐藏路由与稳定 `page:`/`action:` ID；当前矩阵包含 24 个列表表面与 58 个表单节点。
- 已记录当前查询状态、列配置缺口、dirty-state 基线及 API/Toast/maintenance 错误分类。
- 用户已确认 Saved View 采用浏览器 `localStorage`，按 `user.id + pageId + tableId` 隔离；D-003～D-005 已把 R1 信息项 I-037-002～004 的取舍冻结。
- 当前工作树另有 Saved Views 与 dirty-state 的实现切片及对应测试；这些变更是在语义记录后形成的未提交实现事实，不被本 R1 审计用作 R1 完成证据，后续分别由 R2/R3 目标承接并验收。

## 事实边界

只写已发生的扫描、文件和验证结果。R1 已记录信息/方案冻结；Saved View、dirty guard 与跨页面反馈的阶段完成仍需各自目标中的实现事实和回归证据。
