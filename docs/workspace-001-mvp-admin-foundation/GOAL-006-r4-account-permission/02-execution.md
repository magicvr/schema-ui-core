---
id: GOAL-006-r4-account-permission
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# 执行记录 · GOAL-006

## 时间线

### 2026-07-31 · R4 目标立项与范围登记

- `/govern` 复核显式工作区、Charter/VP 对齐、Root 路线图、R3 关门（GOAL-005 A-007/A-008）与父目标信息门禁。
- 创建本目标五件套和 `attachments/` 目录，并将 `GOAL-006-r4-account-permission` 挂到 `GOAL-001-mvp-admin-foundation`。
- 将 R4 范围记录为账号权限最小 API 设计、`D-PERM` 映射冻结与前后端鉴权链路；明确排除 R5 Renderer/业务范例、完整权限继承产品化与完整协议支持。
- 登记 `I-006-001` 为 required/open（R4 方案冻结前验证）；父目标 `I-PROTO-002` 保持 open、作为 R4 **实施**门禁。
- 同步工作区 `goal-tree.md`；Root `status` 保持 `active`、`progress` 保持 `3/6`。
- 本次没有修改 `apps/*`，没有收集或验证任何 `I-006-*`，没有放行方案冻结、实施或 `done`；`I-PROTO-002` / `I-PROTO-003` 未改变。

### 2026-07-31 · R4 契约收集与方案冻结（D-004）

- 从固定 commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`（artifact `2.7.0`）拉取并落盘 D-PERM 资料至 `attachments/dperm/`：
  - `permissions-inheritance/cases.json`（fixtureVersion 1.0，17 cases：13 valid + 4 invalid；target kinds：formField 7 / formSubmit 7 / rowAction 5 / actionButton 5 / toolbarTrigger 2 / column 1）
  - `node.schema.json`（`Permissions` view/edit/delete；`PermissionCascade` keys enum edit|delete、unique、minItems 1）
  - `0023-container-permission-inheritance.md`（ADR-0023：effectivePermission AND 公式、4 条结构边、5 类 cascade type 白名单、columns 不参与、执行时序 fail-closed）
  - `permission-inheritance.md`（v2.7 场景：编辑/删除继承扩展示例，R5 范例页候选）
  - SHA-256 已核验（见 D-004 表）。
- 对照语义规范原文（固定 commit 下 `01-node-protocol.md` §3.9/§3.9.1、`03-component-registry.md` intent 矩阵与 D4a 表单 edit 目标、`08-renderer-spec.md` §7.1 执行时序）与覆盖表 v0.1.3（`D-PERM=include`、`permissions-inheritance=include`），完成最小 API 设计与映射结论。
- 用户确认「按此方向冻结」且「I-PROTO-002 在方案冻结时一并闭合」；D-004 落盘。
- `I-006-001` → `verified`；父目标 `I-PROTO-002` → `verified`（Root meta 同步留痕）。
- 本阶段**未**实施代码（`apps/*` 未修改）；R4 实施仍需用户指令并记实施事实；`I-PROTO-003` / `I-PROTO-004` 未改变。

## 完成后边界

1. R4 方案已冻结（D-004）；`I-006-001` 与父目标 `I-PROTO-002` 已 `verified`。
2. R4 实施（前端权限求值/显隐禁用、Go 会话与鉴权模型、fixture 对照测试）按 02-execution 后续时间线推进并记事实；不把「方案冻结」写成「已实现」。
3. 开放 required 信息项到期前不得越过对应门禁；`I-PROTO-003`（R5 验收/关门）不属本目标处理。

## 进度评估

R4 完成「契约发现与信息就绪」与「方案冻结」两阶段：`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`（证据：D-004 + `attachments/dperm/` 固定资料 SHA-256）。**无实施事实**（`apps/*` 未修改）；Root `progress` 仍为 `3/6`。
