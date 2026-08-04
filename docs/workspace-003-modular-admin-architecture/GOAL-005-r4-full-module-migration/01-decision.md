---
id: GOAL-005-r4-full-module-migration
doc: decision
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-005

## 信息与门禁规则

R4 继承 Root D-009 的阶段入口和 I-PROTO-001 v0.1.3 范围约束。C1 未完成前，
不得冻结全量迁移方案或把 Users/Roles 的现有实现升级为 R4 迁移事实。R4-I003
的 Records/Schema CRUD 冲突属于 P-004 信息冲突，必须以用户书面裁决或新的
canonical 决策记录关闭，不得用本地代码现状静默代替 VP 意图。

## 决策索引

| 编号 | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4 范围、五项检查点与信息门禁 | accepted | [01-decision/D-001-r4-stage-scope.md](01-decision/D-001-r4-stage-scope.md) |

## 当前约束

- R4 必须迁移标准 Admin 模块的六项能力；metadata 声明不能替代真实 provider
  注册、Schema ownership、授权和持久化证据。
- Users/Roles 的现有行为和协议是兼容基线；任何 operationlog 一致性变化需
  新决策和对应测试。
- Records 当前由 `0006 records_retire` 退役，VP-003 又将 `records/Schema CRUD`
  列为 R4；该冲突在 C1 前保持 collecting。
- R4 关键阶段建议使用 Grok Build independent audit；provider 失败不得静默
  降级为 self。
