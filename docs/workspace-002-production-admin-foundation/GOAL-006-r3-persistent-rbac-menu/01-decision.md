---
title: 决策 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 决策 · GOAL-006

## D-001 · 用一个端到端目标实施 Root D-009

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 以单一 R3 子目标承载版本迁移、规范化 RBAC、增量种子、records 读写授权、`features` 菜单投影和恢复/回归；六个成功标准按依赖顺序推进。
  2. API 权限 key 固定为 `records.read` 与 `records.write`；viewer 仅获得前者，admin 获得两者。权限判断来自持久化 role-permission 关系，不再直接判断 `admin` 角色字符串。
  3. `menu_items` 使用唯一 `page_ref` 和显式唯一 `feature_key`；`/api/accounts/me.features` 输出布尔值，静态 manifest 用 `visibleWhen` 消费。Web 投影只控制展示，API 仍独立强制授权。
  4. 迁移分为“建立/回填/双读核对”和“切换规范化读写”两步；旧 `users.roles` 的删除或停用不纳入本目标的不可逆切换。
  5. R3 恢复口径固定为迁移前副本恢复 + `PRAGMA integrity_check` + 身份/授权/菜单/refresh 关键查询；完整生产备份流程留给 R5。
- **理由**：这些变更共享 schema、seed、身份快照和端到端验证，拆成多个并列目标会制造跨目标中间态与重复门禁；顺序检查点能保留可核对的阶段边界。
- **实施门禁**：`I-006-001` 在首个代码变更前关闭；`I-006-002` 在 S5 前关闭。两项只细化已选模型，不重开 Root D-009 的方案裁决。

### 未选方案

- **按 migration / authorization / menu 拆成三个并列目标**：强耦合中间态难以独立验收，且容易让菜单投影先于真实授权或迁移契约落地。
- **先写代码、再补 DDL/feature key 决策**：会把约束、迁移顺序和真实菜单项选择变成隐式事实，违反 P-005。
