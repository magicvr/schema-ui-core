---
id: A-002-r4-c2-contract-review
doc: audit-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: self
date: 2026-08-05
scope: R4-C2 contract implementation slice, conflict/fail-closed semantics, verification
verdict: conditional
---

# A-002 · R4-C2 契约实施 self 审计

## Finding closure

| finding | closure | 证据路径 |
|---------|---------|----------|
| F-C2-001（C2-I004 待实施证据） | `fixed` | E-002：kernel 新增 contribution/provider/persistence 契约层；全量测试 + vet 通过 |

## 已核实成果

- `kernel/contribution.go`：六类 Contribution + `ContributionIdentity` 规范 Key 语义
  （HTTP "METHOD pattern"、Schema PageID、Auth Permission、Navigation NodeID、
  Manifest FragmentID、Persistence Name），字段校验（identity 匹配、owner、tombstone
  互斥、reconcile 一致性）。
- `kernel/provider.go`：`Provider`/`Registrar`（Registrar 无 Persistence 方法）；
  `RegisterContributions` 只注册启用 provider、Descriptor 精确匹配、Register 仅写已
  声明 Kind+Key、finalize 全局冲突/引用/capability/确定性排序校验、失败丢弃整个集合。
- `kernel/persistence.go`：`CollectPersistence` compiled-global catalog（version/name/
  checksum/（module,name）唯一、无缺口、tombstone/reconcile 一致、确定性排序）。
- 测试：happy path、跳过未迁移/非启用模块、Descriptor mismatch、未声明 key、identity
  不匹配、跨模块冲突、navigation 悬空 permission、fragment capability 缺失、确定性、
  双 Profile 契约矩阵；persistence 的 gap/重复/校验全通过。
- fx 静态检查：kernel 与 modules 无 `go.uber.org/fx` import；composition 允许（Fx 根）。

## Open required

无。C2.4 的运行时双 Profile 矩阵（register/conflict/Start/Ready）与 readyz 真实
readiness 需真实模块 provider 接线，登记为 C3/C5 门禁，不阻断 C2 关门。

## Gate

C2 保持 `active 4/4`；C2.1-C2.4 勾选。Grok independent 复审未完成前不关门。
本意见不修改 status/progress。
