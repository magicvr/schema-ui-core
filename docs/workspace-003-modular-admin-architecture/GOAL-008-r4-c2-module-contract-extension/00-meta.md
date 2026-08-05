---
id: GOAL-008-r4-c2-module-contract-extension
title: R4-C2 · 模块契约扩展
status: done
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 4/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 GOAL-005 C1 冻结的 Provider 精确契约，扩展 Kernel/组合根模块贡献边界为可验证的 HTTP/Schema/Authorization/Navigation/Manifest/Persistence provider，先用冲突与最小模块测试验证，不迁移业务。
---

# GOAL-008 · R4-C2 模块契约扩展

## 概述

本子目标是 `GOAL-005-r4-full-module-migration` 的 C2 子目标，承接冻结包
`GOAL-005/attachments/r4-c1-freeze-package-draft.md`（`status: accepted`）的
Provider/Registrar 精确契约。它只扩展 Kernel/组合根模块贡献边界并用最小模块与
冲突测试证明契约成立，**不迁移** Users/Roles/Settings/Activity（属 C3/C4），
不把现有中心注册代码宣称为已完成迁移。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-005-r4-full-module-migration` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `independent`；实现切片使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [x] **C2.1 / Provider 契约实现**：framework-agnostic `Provider` + Plan-owned
  `Registrar` + 六类 Contribution 结构化类型按冻结包 §2 落库；`ContributionIdentity
  .Key` 语义校验；模块公共 API 与 `kernel` 不 import Fx。
- [x] **C2.2 / compiled-global Persistence catalog**：`Provider.CompiledPersistence()`
  为唯一收集入口，`Registrar` 无 Persistence 方法；catalog 静态校验按冻结包 §4
  （全局 version/name/checksum/（module,name）唯一、版本连续无缺口、tombstone 与
  Apply 互斥、reconcile 元数据一致、确定性按版本排序）。ledger 侧「漂移 / 未知已
  应用记录」运行时 fail-closed 与 fresh/upgrade/reconcile 事务路径属 C3/C5 数据
  门禁（冻结 §4.2）。
- [x] **C2.3 / 双检与失败语义**：Plan 解析 → Descriptor 精确匹配（ID/Version/
  KernelAPIRange/DependsOn/Provides/Requires/ContributionKeys 全键集）→ Register
  仅写声明 Kind+Key → 全局冲突/引用完整性/capability/确定性排序校验；register/
  conflict 失败丢弃整个集合，不发布部分 surface。Manifest secrecy 扫描与
  Start/Ready 反向清理的运行时矩阵属 C3/C5（冻结 §3 步骤 4/8）。
- [x] **C2.4 / 最小模块与矩阵**：冲突 + 最小模块测试证明 provider contract；静态检查
  `modules/**` 与 `kernel` 无 Fx import；`mvp`/`admin` 双 Profile 契约矩阵通过；
  readyz 不冒充模块图 readiness（C2 未触及 readyz，未作宣称）。

四个检查点等权；当前 `progress: 4/4`。完成本子目标只表示 C2 关闭，不关闭 GOAL-005、
Root 或 VP-003，不自动放行 C3。双 Profile 的 register/conflict/Start/Ready 运行时
矩阵与 readyz 真实 readiness 需真实模块 provider 接线，属 C3/C5 门禁。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| C2-I001 | required | Provider/Registrar 精确契约与六类 Contribution 字段是否可实施？ | C2.1、所有迁移 | C2.1 | 用户整包接受冻结包；C1 冻结 | verified | freeze package §2 `accepted`；GOAL-005/006 D-003 |
| C2-I002 | required | compiled-global Persistence 收集、校验和 reconcile 规则是否冻结？ | C2.2、C3/C5 数据门禁 | C2.2 | 冻结包 §4；GOAL-006 C1-I001 | verified | freeze package §4 `accepted` |
| C2-I003 | required | fail-closed 冲突/生命周期/失败清理语义是否冻结？ | C2.3、C3/C4 切换 | C2.3 | 冻结包 §3 | verified | freeze package §3 `accepted` |
| C2-I004 | required | 当前 Kernel/Composition 状态与冻结契约的差距，以及 C2 最小模块/冲突测试的实施证据？ | C2.4 验收、C3 起点 | C2.4 | C2 内核对 `kernel/module.go`、`composition.go` 并落实施证据 | verified | E-002：kernel 新增 contribution/provider/persistence 契约层；全量测试 + vet 通过 |

## 阶段路线图

1. 继承冻结包契约，核对当前 Kernel/Composition 与契约差距（信息门禁就绪）。
2. 实现 Provider/Registrar/Contribution 类型与 `CompiledPersistence()` 收集（C2.1/C2.2）。
3. 实现 Plan 双检、全局冲突/引用/secrecy/确定性排序校验与失败清理（C2.3）。
4. 冲突 + 最小模块测试、fx 静态检查、双 Profile 矩阵（C2.4），self + Grok
   independent 复审后关门（已完成：Grok A-003 三条 required 经 A-004 闭合）。

## 范围与非目标

范围包括 Provider/Registrar/Contribution 契约、compiled-global Persistence 收集、
Plan 双检、确定性聚合、失败闭锁、冲突/最小模块测试、fx 静态检查和双 Profile 验证。
非目标包括 Users/Roles/Settings/Activity 业务迁移（C3）、Schema 其他能力迁移（C4）、
旧中心注册/Schema owner map/Manifest `adminModules` 的最终移除（C3/C4）、R5/R6。
