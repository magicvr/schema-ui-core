---
id: GOAL-009-r4-c3-users-roles-migration
doc: decision
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-009

## 信息需求与阶段门禁

| 编号 | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|------|------|-----------------|----------|--------------|-----------------|------|-------------|
| C3-I001 | required | Users/Roles 当前中心注册/Schema/Manifest/seed ownership | C3.1/C3.3 | C3.1 | 全仓扫描 | verified | E-002 |
| C3-I002 | required | C3 保留行为矩阵枚举 | C3.2/C3.4 | C3.1 | 冻结包 §7 + 现有测试 | verified | E-002 行为矩阵 |
| C3-I003 | required | operationlog 失败注入测试 | C3.4 | C3.4 | 失败注入测试 | verified | store seam + handler test |
| C3-I004 | non-blocking | 双 Profile 矩阵 + Manifest secrecy | C3.4 | C3.4 | GOAL-008 E-004 | open | GOAL-008 E-004 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4-C3 Users/Roles 迁移子目标 | accepted | [01-decision/D-001-r4-c3-stage-scope.md](01-decision/D-001-r4-c3-stage-scope.md) |

## 当前约束

- 承接冻结包 §7 切换顺序与 GOAL-008 Provider 契约；C3 只迁移 admin.users/admin.roles，
  不宣称 C4/C5、不推进 Root progress。
- 保留现有行为矩阵（CRUD、授权、角色分配、最后管理员保护、密码、operationlog
  best-effort）；`0003`/`0006` 迁移账本与历史 operation-log 保留。
- 审计模式 `independent`；迁移切片使用 Grok Build `grok-4.5` / `high`。
