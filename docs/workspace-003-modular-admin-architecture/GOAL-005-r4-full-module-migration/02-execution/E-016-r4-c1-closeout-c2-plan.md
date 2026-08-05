---
id: E-016-r4-c1-closeout-c2-plan
doc: execution-entry
goal: GOAL-005-r4-full-module-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-016 · R4-C1 门禁闭合与 C2 计划边界

## C1 闭合事实

- GOAL-006 `done 4/4`（Grok A-006 `pass`）：Provider 精确契约整包冻结、Records
  historical-only、operationlog Option A + bounded residual 均经用户 D-003 裁决。
- GOAL-007 `done 4/4`（Grok A-003 `pass`）：Records 运行面核验完成，产品面无残留、
  兼容边界保留、防复活测试补齐。
- R4-I001/I002/I003 `verified`，R4-I004 `accepted-residual`，R4-I005 non-blocking
  open。A-008 按 ID 汇总全部 C1 required finding 闭合。
- C1 检查点成立（范围与信息冻结）；GOAL-005 progress `0/5 → 1/5`。

## C2 计划边界（冻结包 §8 要求写入 execution ledger）

C2 只做**模块契约扩展**，不迁移业务（C3）：

1. **Provider/Registrar/Contribution 类型**：按冻结包 §2 落库 framework-agnostic
   `Provider` + Plan-owned `Registrar` + 六类 contribution 结构化字段；模块公共 API
   不 import Fx；`ContributionIdentity.Key` 语义校验。
2. **compiled-global Persistence**：`Provider.CompiledPersistence()` 为唯一收集入口，
   `Registrar` 无 Persistence 方法；迁移/seed/reconcile 规则按冻结包 §4（版本/name/
   checksum/tombstone/reconcile 全局校验）。
3. **Plan 解析与运行时双检**：按冻结包 §2.3（composition 建 registry → 解析 Profile
   → Descriptor 匹配 → Register 仅写声明 Kind+Key → 全局校验 → finalize）；未声明/
   重复/字段身份不匹配 fail closed。
4. **确定性聚合与失败语义**：HTTP/Schema/Auth/Nav/Manifest 冲突、引用完整性与
   Manifest secrecy fail closed；任一 register/conflict/Start/Ready 失败不得留下部分
   surface，按顺序反向清理（冻结包 §3）。
5. **最小模块 + 冲突测试**：先用冲突与最小模块验证 provider contract；静态检查
   `modules/**` 与 `kernel` 不 import Fx；`mvp`/`admin` 双 Profile 各跑一次矩阵。
6. **C2 不迁业务**：Users/Roles/Settings/Activity 迁移属 C3/C4；readyz 不得冒充模块
   图 readiness。

## C2 信息门禁（承接 C2-I001..C2-I003，见 GOAL-008 meta）

| 编号 | 级别 | 状态 |
|------|------|------|
| C2-I001 | required | Provider/Registrar 精确契约与 Contribution 字段可实施性（冻结包 verified） |
| C2-I002 | required | compiled-global Persistence 收集规则（冻结包 verified） |
| C2-I003 | required | fail-closed 冲突/生命周期语义（冻结包 verified） |
| C2-I004 | required | 当前 Kernel/Composition 状态与冻结契约差距的 C2 实施证据（collecting，C2 内收集） |

## 提交

C1 闭合与 C2 开设已 git 提交（标题见提交日志）。
