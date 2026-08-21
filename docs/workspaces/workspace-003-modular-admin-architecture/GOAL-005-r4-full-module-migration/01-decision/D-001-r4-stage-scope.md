---
id: D-001-r4-stage-scope
doc: decision-entry
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: accepted
---

# D-001 · R4 阶段范围与信息门禁

## 决定

建立 R4 主实施子目标，按 C1→C5 渐进推进。第一阶段只收集并冻结范围、
provider contract 和行为边界；C2/C3/C4 实施必须等待 C1 required information
关闭。

R4 直接承接 Users、Roles 和 VP-003 指定的 `records/Schema CRUD` 范围。当前
Records 已被 `0006 records_retire` 删除产品表/权限，只保留历史 operation-log
事件；该现状与 VP-003 用语冲突，登记为 R4-I003，等待用户裁决，不预先选择恢复
或永久退役路径。

## 审计安排

R4 属于 migration/production 高影响范围。C1 方案冻结和 C5 close-out 计划
执行 self + Grok Build independent audit，指定 `grok-4.5`、`high`；独立意见
只写入 R4 audit ledger，由主治理记录响应。
