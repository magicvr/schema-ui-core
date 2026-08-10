---
id: GOAL-007-s6-settings-form-page
title: S6 · 设置页表单/详情页改造（recordSource 预填）
status: done
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.3.0
progress: 4/4
---

# GOAL-007 · S6 · 设置页表单/详情页改造

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S6 阶段**：将 `settings` 页从「1 行 × 9 列单例表格 + 5 个工具条弹窗」重构为 **General / Branding / Localization / Appearance 四类分组就地编辑的表单/详情页**（显示当前值、就地保存、Restore defaults 独立确认按钮）。

实现协议层**既有但渲染器未接线**的 `form.props.recordSource` 预填能力（ADR-0021，since 2.1，capability `form.record.load`）：解析放行、GET 预填、`responseMapping` 映射、capability 门禁、reload 重预填、只读权限门禁。**不新增协议字段**——仅实现已冻结、已验证的协议表面。

**背景**：设置页是单一全局配置，读（表格）与写（分散 modal）断开，且弹窗打开时字段为空（工具条触发 `invokeAction(trigger, null)`），空表单一保存有清空风险。VP-007 已写明产品结构为「四类组织」。用户 2026-08-09 选定方案 A（真·表单/详情页），并要求本目标承载治理上下文、Root 暂时回退关门状态。

**方案依据**：本目标 `01-decision/D-001-settings-form-page-scope.md`（范围、recordSource 路径、审计模式）。

## 成功标准（可验收 · 等权检查点 · 共 4 项）

- [x] **C1**：renderer `recordSource` 预填能力落地——capability `form.record.load` 门禁（缺 → fail-closed）、loading/error 状态、reload 重预填、form `title/titleKey` 标题、`${formId}:submit` 只读门禁、`invokeAction` actionId 回退、actionButton 权限禁用。
- [x] **C2**：settings schema 重构——四类 recordSource 预填内联表单 + Restore defaults actionButton；删除表格/modal；meta 增 `form.record.load`；全键复用现有 catalog。
- [x] **C3**：测试与证据——vitest **727/727**（新增 renderer 用例 + `startup-config` 改写 + `schema-keys.structural` 保持）；`npm run build` exit 0；`go test ./apps/api/...` exit 0；e2e M3 已改写（本机 8080 端口排除区间 → 诚实降级留痕，M3 逻辑由单元覆盖）。
- [x] **C4**：治理收口——evidence 入库（E-001/E-002 + A-001/A-002）；independent 关门审计 A-002 **pass**（required 0，F-001 fixed / F-002 accepted-residual）；**2026-08-09 用户书面确认关门**（D-002）→ GOAL-007 `done` `4/4`；Root `GOAL-001` 恢复 `done`（`7/7`，临时回退解除）。

## 派生进度展示

`progress: 4/4` 由上方 4 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

> P-005：本目标无到期 required 信息门禁。实现方案（recordSource 复用既有协议表面）已由 D-001 冻结。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | recordSource 既有协议语义（capability 门禁、responseMapping、search 禁止）是否与 renderer 接线一致 | C1/C2 | C1 实施前 | 核对 `component-registry.json` + `request-construction.ts` + conformance cases | **verified** | — | D-001 §2 盘点（2026-08-09） |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S6 阶段子目标）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 三目录平铺；索引文件在 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- Root `GOAL-001` 于 2026-08-09 用户指令**暂时回退**关门状态（`done` → `active`，progress `6/6` → `6/7`）以承接本 S6 子目标；历史关门记录 D-002 保留不重写，重新开根决策见 GOAL-001 `01-decision/D-003-reopen-root-for-s6.md`。**C4 用户书面确认（D-002）后 Root 已恢复 `done`（`7/7`），解除临时回退**（GOAL-001 `D-004`）。
- VP-007 保持 `closed`（本目标为已关闭波次上的增量产品化，非新愿景波次；不触碰 vision 状态）。
