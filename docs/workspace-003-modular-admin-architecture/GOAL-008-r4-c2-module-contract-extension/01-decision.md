---
id: GOAL-008-r4-c2-module-contract-extension
doc: decision
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-008

## 信息需求与阶段门禁

| 编号 | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|------|------|-----------------|----------|--------------|-----------------|------|-------------|
| C2-I001 | required | Provider/Registrar 精确契约与六类 Contribution 字段可实施性 | C2.1 | C2.1 | 冻结包 §2 | verified | freeze package `accepted`；D-003 |
| C2-I002 | required | compiled-global Persistence 收集/校验/reconcile 规则 | C2.2 | C2.2 | 冻结包 §4 | verified | freeze package §4 |
| C2-I003 | required | fail-closed 冲突/生命周期/失败清理语义 | C2.3 | C2.3 | 冻结包 §3 | verified | freeze package §3 |
| C2-I004 | required | 当前 Kernel/Composition 与契约差距 + C2 实施证据 | C2.4 | C2.4 | C2 内核对与落证据 | verified | E-002：契约层 + 测试 + vet 通过 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4-C2 模块契约扩展子目标 | accepted | [01-decision/D-001-r4-c2-stage-scope.md](01-decision/D-001-r4-c2-stage-scope.md) |

## 当前约束

- 承接冻结包精确契约（C1 冻结）；C2 只扩展契约并用冲突/最小模块测试验证，不迁移
  业务、不宣称 C3/C4 完成、不推进 Root progress。
- C2 不得在未记录的情况下改变身份、冲突键、安全语义或注册/发布顺序；
  `ConfigNamespaces` 不新增独立 Registrar 方法。
- 审计模式 `independent`；实现切片使用 Grok Build `grok-4.5` / `high`。
