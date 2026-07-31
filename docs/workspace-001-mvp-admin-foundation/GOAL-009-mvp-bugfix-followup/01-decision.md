---
id: GOAL-009-mvp-bugfix-followup
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.1.0
---

# 决策记录 · GOAL-009

## 信息需求与阶段门禁

与 [00-meta.md](00-meta.md) 同表：`I-009-001` / `I-009-002`（recommended 是否升格必做）仍 open，不阻断 required 五项的实施启动。

## D-001 · 立项为 Root 关门后的修正跟随子目标

**决定**：

在 `workspace-001-mvp-admin-foundation` 新建 `GOAL-009-mvp-bugfix-followup`，`parent: GOAL-001-mvp-admin-foundation`。Root 保持 `status: done` / 纲领 `progress: 6/6`；VP-001 保持 `closed`。本目标不重开 VP 范围，只修正审视认定的 bug 与失真。

**为什么**：

- 用户明确要求：工作区 1、GOAL-001 子目标、审视内容进 attachments、目的为补修 bug。
- VP/Root 已关门的交付事实不因后续修正而回滚；修正以独立子目标可审计、可关门。

**未选方案**：

- **改写已 done 的 GOAL-006/007 执行史冒充当时已修**：违反“只记事实”。
- **新开工作区 / 新 VP**：范围过小，且用户指定工作区 1 + 父 GOAL-001。
- **静默改代码不立项**：无法对照 findings 关门。

## D-002 · 范围边界：required 五条 + recommended 可选

**决定**：

- **必做（成功标准）**：F-009-001～F-009-005（见附件）。
- **建议（默认尝试，可 residual）**：F-009-006 路由级 `Allow`、F-009-007 body/pageSize 上限。
- **排除**：`schemaUrl` 通用管线、host 补全、真实 IAM、upload/batch、多浏览器 e2e 矩阵。

**为什么**：

- 附件已区分真 bug vs 有意 MVP 边界；成功标准必须可在本目标内验证关闭。
- 006/007 改善安全与稳健性，但非审视当日崩溃类缺陷；升格与否由用户/后续决策。

**未选方案**：

- **把架构债一并纳入本目标**：范围膨胀，与“bug 修正”目的不符。
- **只修 API 不碰 Web 集成失真**：F-009-002 是权限演示可信度的核心问题。

## D-003 · 审视意见落盘方式

**决定**：

独立长文写入 `attachments/audit-code-review-bugs-2026-08-01.md`；本目标 `03-audit.md` 以 **A-001**（`source: independent`）索引摘要 + findings，不把长文只留在聊天。

**为什么**：

- 符合 P-003：正式意见在被审目标台账；长文可附件。
- 用户明确要求审计内容作独立附件。
