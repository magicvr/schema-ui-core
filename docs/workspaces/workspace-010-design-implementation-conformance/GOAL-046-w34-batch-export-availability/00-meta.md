---
id: GOAL-046-w34-batch-export-availability
title: W34 · 「导出所选」可用性（operator config 漂移根因 + 触发面门禁）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-046 · W34 · 「导出所选」可用性

## 概述

用户在 VP-038 关门后报告（2026-09-19）：

> 现在无论是用户列表页还是角色列表页，选中列表项并点击「导出所选」的时候，都会报错宣示「未找到」，而并不会正常导出。需要修正一下这个问题。

诊断结论：**两个独立缺陷叠加**，均属 VP-038 已交付面的可用性缺口，而非 R3/R4 的导出逻辑错误。

1. **operator config 漂移（用户环境的直接原因）**：`apps/api/configs/config.yaml` 声明 `profile: custom`，其注释自称「= the full admin preset PLUS channel.telegram」，但 VP-038 把 `admin.jobs` 加入 admin preset（`kernel/profile.go`）时**漏改这份列表**。模块未装配 → `POST /api/jobs/batch-export` 未挂载 → 路由回落 `{"error":"NOT_FOUND","messageKey":"error.notFound"}`（zh-CN「未找到」）。schema 节点属于 `admin.users`/`admin.roles`，在所有 profile 照常渲染，于是按钮看似可用、点击即失败。
2. **触发面缺少可用性门禁（mvp/demo 同类隐患）**：`mvp`/`demo` preset 既无 `admin.jobs`（`jobs.write`）也无 `admin.data-transfer`（`data.export`），而「导出所选」节点随 users/roles 模块发布到所有 profile，因此在这些 profile 下同样只有 404/403 一条路。

本目标于 2026-09-19 立项，随独立审计 `A-002` 的 required 闭合（`A-003` 响应）当日关门：`done · 4/4`。

## 交付摘要

- **① operator config 修复**：`configs/config.yaml` 补 `- admin.jobs`（含成因注释）。修复后真实服务端到端：users/roles 均 `202 queued → succeeded`，结果可下载为 CSV（2155 字节实测）。
- **② 触发面可用性门禁**：`jobs-batch-export.tsx` 按路由真实的两道门禁（`jobs.write` + `data.export`）判断可用性；不可用时渲染为 **disabled + 文案说明**，而非隐藏——遵循 W33 `D-001` §3 冻结的 fail-open 取向（静默移除操作入口是更严重的失败模式），并且空插槽宿主会破坏 page-actions 行的控件高度契约（该回归被 `list-visual-surface` e2e 实测捕获后在本波次修正）。上下文未提供 permissions 时保持 fail-open。
- **③ 404 语义纠正**：提交命中 404 时报「当前部署未启用批量导出（未装配 admin.jobs 模块）」，不再把服务端裸 `NOT_FOUND` 当作「数据不存在」抛给用户。新增 i18n 键 `schema.jobs.batchExport.unavailable`（en/zh）。
- **④ 永久守卫（防同类漂移）**：`internal/config` 新增 `TestOperatorConfigCoversAdminPreset` —— 断言 operator config 的内联模块列表是 admin preset 的**超集**（显式允许 `channel.telegram` / `biz.digital-offer` 等额外项）且必须含 `admin.jobs`。变异验证：删除 `- admin.jobs` → 守卫失败并指名模块。
- **⑤ 回归锁**：`jobs-batch-export.test.tsx` 新增 4 例可用性用例（两门禁缺一即不可用 / 齐备可用 / 上下文未知 fail-open / 404 文案）。变异验证：去掉门禁判断或把 404 文案改回裸 message，对应用例即红。

## 范围与非目标

### 本目标范围

- `configs/config.yaml` 模块列表修正；批量导出触发面的可用性渲染与 404 文案；上述两项的守卫/回归与愿景层登记。

### 明确非目标

- 不改 VP-038 `status`（保持 `closed`）与 workspace-038 台账正文；不改 Job 六态合同、批量导出契约、两道服务端门禁的语义（本波只做**客户端可用性镜像**，服务端仍是唯一授权方）；不把导出搬到 mvp/demo（是否进 preset 属 Profile 内容决策，本波不动）；不触碰 pinned 协议工件。

## 成功检查点

- [x] **C1 诊断与方案冻结**：`D-001` 落盘（两道根因、fail-open 决策、守卫设计、授权范围）；`I-046-001`/`002`/`003` 关闭。
- [x] **C2 实施**：operator config 补 `admin.jobs`；触发面按两道门禁渲染；404 文案纠正；i18n 补齐。证据：`02-execution/E-001` §2。
- [x] **C3 验证与证据**：config 守卫（含变异）+ 可用性用例（含双向判别性变异）+ 真实服务端到端（202→succeeded→CSV）+ 全量回归（vitest 121 files / 1482 tests、typecheck/build、e2e mvp 18/5/0 与 admin 19/4/0、Go 全绿）+ 双 profile 可用性 e2e 契约入仓。
- [x] **C4 审计与投影**：self `A-001` `pass`（0 required）→ independent `A-002`（grok build）`conditional`（2 required）→ `A-003` 响应（2 `fixed`，开放 required = 0）；`goal-tree.md`/`workspace.md`/`roadmap.md` 同步；Root 保持 active 程序容器。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-046-001 | required | 用户环境「未找到」的确切来源：路由未挂载、权限拒绝，还是数据缺失 | C1 | C1 | 原样复现 operator config 启动 + `curl -i` 观察状态/`messageKey`；与 admin preset 对照 | **verified** | — | `404 NOT_FOUND`/`error.notFound`（路由未挂载）；admin preset 同请求为 403 路由存在；`E-001` §1 |
| I-046-002 | required | 不可用时应隐藏还是可见但禁用（与 W33 `D-001` §3 fail-open 冻结取向的关系） | C2 | C1 | 复核 W33 冻结文本 + page-actions 高度契约的 e2e 约束 | **verified** | — | `D-001` §3：可见 + disabled + 说明；隐藏既违背冻结取向，也会留下高度 0 的空插槽宿主 |
| I-046-003 | non-blocking | 是否有既有 e2e/单测覆盖「节点已发布但路由未挂载」这一形态 | C3 | C2 | 全仓检索批量导出相关断言 | **verified** | — | 原覆盖只到「admin profile 下可用」；缺口由本波 config 守卫 + 4 例可用性用例补齐 |

## 父目标

- `[workspace-010-design-implementation-conformance]` `GOAL-001-design-implementation-conformance`（长期程序容器，保持 active）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标是 VP-010 持续符合性程序的一个**波次子目标**，不是新 VP、不改变 Charter；Root 保持 active 程序容器。
- VP-038 保持 `closed`：本波只闭合「已交付面在 operator config 下不可用」这一缺口，并在 `roadmap.md`「未决项统一登记」登记，不重开 VP-038、不改 workspace-038 台账。
- 审计模式 `self` + independent（`cross`）：变更为**生产/运维配置**（模块装配）+ 数据外带入口的可用性面，按 `AGENTS` §6b 风险表取 cross；independent provider 按项目级决策 = 本地 grok build（grok-4.6 · high · `/audit`）。
