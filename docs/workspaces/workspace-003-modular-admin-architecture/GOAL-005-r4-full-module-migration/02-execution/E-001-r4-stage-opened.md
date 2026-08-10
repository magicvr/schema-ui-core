---
id: E-001-r4-stage-opened
doc: execution-entry
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: recorded
---

# E-001 · R4 阶段建立与初始边界扫描

- Root D-009/E-009/A-008 已确认 R3 `done 4/4`、I-006 verified、Root progress
  `3/6`，因此建立 R4 子目标。
- 当前 API 组合根仍以 `handler.Register(..., plan)` 作为中央边界；Settings/
  Activity 已有试点模块包，Users/Roles 尚未拥有 module-owned registration。
- Kernel contribution metadata 当前声明 Routes/Pages/Navigation/Permissions/
  ConfigNamespaces，但没有结构化 provider 字段；这形成 R4-I002。
- Users/Roles 的 handler/store/schema/renderer/e2e 行为基线已经定位；Records
  表和权限已由 `0006 records_retire` 清除，历史 operation-log 事件仍保留。
- 因 VP-003 R4 仍写明 `records/Schema CRUD`，R4-I003 保持 collecting；在该信息
  冲突关闭前不推进 C2。
