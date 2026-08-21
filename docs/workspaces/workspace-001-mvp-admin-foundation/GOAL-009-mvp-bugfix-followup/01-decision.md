---
id: GOAL-009-mvp-bugfix-followup
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.3.0
---

# 决策记录 · GOAL-009

## 信息需求与阶段门禁

与 [00-meta.md](00-meta.md) 同表：`I-009-001` / `I-009-002`（recommended 是否升格必做）已 **resolved**——2026-08-01 用户书面裁决「都纳入实施」（见下方 D-004），F-009-006/007 已实施并 `fixed`。

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

## D-004 · F-009-006/007 纳入实施（2026-08-01）

**决定**：

用户裁决 I-009-001/002 为「都纳入实施」：F-009-006（records 写路由 fail-closed 鉴权）与 F-009-007（body/pageSize 上限）从 recommended 升格为本目标实施项，不做 accepted-residual。

**为什么**：

- F-009-006 严重度 med（误用时 high）：演示 API 存在鉴权库却不挂写路由，易与「后端独立鉴权」叙事混淆；挂 fail-closed 检查成本低、收益明确。
- F-009-007 低风险高性价比：`MaxBytesReader` + `pageSize` 上限 + 400 测试改动极小。

**未选方案**：

- **对 006/007 作 accepted-residual**：用户否决；两条均可在 MVP 范围内低成本落地。
- **只做一条**：用户选择两条都做，避免留半截鉴权/上限。

**执行证据**：`records.go` `writeGate()`（401/403）+ `MaxBytesReader` 4 KiB + `pageSize ≤ 100`；`TestRecordsWriteRequiresSession`、`TestRecordsWriteDeniedWithoutAdminRole`、`TestRecordsUpdateBodyTooLarge`、`TestRecordsListPageSizeCap`。

## D-005 · 关门授权（2026-08-01）

**决定**：

用户 P-004 裁决：**不补 self 关门审**，接受 A-002（independent 关门复审，`verdict: pass`）作为关门审计依据；关闭前修复 A-002 两条 recommended（F-A002-001 README 措辞、F-A002-002 台账卫生）；随后授权 **GOAL-009 → `status: done`**。

**为什么**：

- A-002 已独立复核 A-001 七条 `fixed` 证据并复跑 `go test ./...` + 聚焦 vitest（7/7 pass）；open required = 0。
- 两条 recommended 为低成本文档修订，关闭前修复使关门干净，避免把 hygiene 债带到后续。
- 关门审计要求（self 或 independent）由 A-002 满足；补 self 审无额外价值。

**未选方案**：

- **补 self 关门审后再关**：用户选择不需，A-002 独立 pass 已足够。
- **对 002 作 accepted-residual**：改动极小，直接修复优于留残余。

**关门范围**：不改 Root `GOAL-001` `done` / 纲领 `6/6` / VP-001 `closed`；浏览器手测保持 optional。
