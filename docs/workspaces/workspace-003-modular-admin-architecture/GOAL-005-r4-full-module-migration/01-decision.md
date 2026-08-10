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
的 Records/Schema CRUD 冲突已由用户书面裁决为 historical-only，并由 D-003 固定
为当前 R4 范围；不得恢复 Records 产品 CRUD，也不得删除既有迁移账本或历史
operation-log 兼容语义。

## 决策索引

| 编号 | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4 范围、五项检查点与信息门禁 | accepted | [01-decision/D-001-r4-stage-scope.md](01-decision/D-001-r4-stage-scope.md) |
| D-002 | 2026-08-05 | 响应 C1 能力盘点 finding | accepted | [01-decision/D-002-r4-c1-inventory-response.md](01-decision/D-002-r4-c1-inventory-response.md) |
| D-003 | 2026-08-05 | Provider、Records 与 operationlog P-004 裁决 | accepted | [01-decision/D-003-r4-c1-decisions.md](01-decision/D-003-r4-c1-decisions.md) |

## 当前约束

- R4 必须迁移标准 Admin 模块的六项能力；metadata 声明不能替代真实 provider
  注册、Schema ownership、授权和持久化证据。
- C1 能力盘点已由 D-002/E-005 完成事实响应；该响应不等于 provider contract
  已冻结，也不等于业务模块已迁移。
- Users/Roles 的现有行为和协议是兼容基线；任何 operationlog 一致性变化需
  新决策和对应测试。
- Records 已由 `0006 records_retire` 从产品面退役；D-003 将其固定为
  historical-only。GOAL-007 负责证明当前运行面没有 Records 专属实现和误导性演示
  命名；迁移账本、历史 operation-log 和历史治理证据保持不变。
- operationlog 采用 Option A：业务成功后的 best-effort append，R4 不自动
  purge/archive/delete；append 失败可能造成审计缺口的 residual 已由用户接受，
  owner 为 `magicvr`，复核时间为 `2026-08-05 08:32:22 +08:00`，触发条件为合规/
  运营 retention 要求、日志规模阈值、恢复演练缺口或进入 R5 数据生命周期决策。
- R4 关键阶段建议使用 Grok Build independent audit；provider 失败不得静默
  降级为 self。
