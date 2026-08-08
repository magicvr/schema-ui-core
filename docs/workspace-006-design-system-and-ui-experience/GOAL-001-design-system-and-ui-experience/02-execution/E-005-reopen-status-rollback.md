---
id: E-005-reopen-status-rollback
title: 回退 Root/工作区 done 状态并同步 goal-tree（A-006/D-006）
date: 2026-08-09
status: done
parent: GOAL-001-design-system-and-ui-experience
---

# E-005 · 状态回退事实

## 触发

- 用户：打开应用后样式与参考页/先前界面无差别 → 要求核查。
- 核查结论落盘：A-006（fail）。
- 用户指令：落盘审计意见并回退工作区完成状态。

## 已写入 / 已修改路径

| 动作 | 路径 |
|------|------|
| 新建审计 | `03-audit/A-006-visual-fidelity-premature-closeout.md` |
| 新建响应 | `03-audit/A-007-response-a006-reopen.md` |
| 更新索引 | `03-audit.md` |
| 新建决策 | `01-decision/D-006-reopen-after-visual-fidelity-audit.md` |
| 废止标记 | `01-decision/D-005-root-closeout-user-confirmed.md` → `superseded` |
| 决策索引 | `01-decision.md` |
| 本执行记录 | `02-execution/E-005-reopen-status-rollback.md`；`02-execution.md` |
| Root 元数据 | `00-meta.md`：`active`，`progress: 2/5`，S2/S3/S5 取消勾选 |
| 工作区 | `workspace.md`：`active` |
| 目标树 | `goal-tree.md` 同步 |
| GOAL-003 | `00-meta.md` / `03-audit.md`：回 `active`，检查点取消勾选，登记 under-delivery |

## 未改

- **未**修改 `apps/web` 业务/样式实现（本条目仅为治理回退）。
- **未**改 VP-005 `status`（仍 `active`）。
- **未**删除 S1–S5 历史实施 commit 或 GOAL-002/004/005 完成记录。

## 回退后开放门禁（事实）

- Root 开放 required：**F-VUI-001**、**F-VUI-002**（见 A-006）。
- 在其合法闭合前，禁止再次勾选 S2/S3、禁止 Root/`workspace` `done`。
