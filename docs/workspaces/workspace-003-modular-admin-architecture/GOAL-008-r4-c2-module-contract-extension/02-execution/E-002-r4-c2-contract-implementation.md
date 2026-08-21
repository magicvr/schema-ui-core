---
id: E-002-r4-c2-contract-implementation
doc: execution-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R4-C2 模块契约实现切片

## 已发生事实

按冻结包 §2/§3/§4 在 `apps/api/internal/kernel/` 实现契约层：

- `contribution.go`：六类结构化 Contribution（Route/Page/Permission/Navigation/
  Fragment/Migration）+ `ContributionIdentity`；`Key` 规范语义（HTTP = "METHOD
  pattern"，Schema = PageID，Auth = Permission，Navigation = NodeID，Manifest =
  FragmentID，Persistence = Name）；字段校验（identity 匹配、owner、tombstone/
  Apply 互斥、reconcile 一致性）。
- `provider.go`：`Provider`/`Registrar` 接口（`Registrar` 无 Persistence 方法）；
  `RegisterContributions` 只注册启用 provider（计划中无 provider 的模块保持中心注册，
  冻结 §7 过渡态），Descriptor 必须与 Plan 模块精确匹配，`Register` 只能写已声明
  `Kind + Key`，finalize 前做全局冲突/引用完整性/capability/确定性排序校验，任一
  失败丢弃整个集合。
- `persistence.go`：`CollectPersistence` compiled-global catalog（唯一收集入口
  `CompiledPersistence()`），全局 version/name/checksum/（module,name）唯一，版本
  连续无缺口，tombstone/reconcile 一致性，确定性按版本排序。
- `module.go`：`ContributionKeys` 增加 `Fragments` 声明字段（加法，不破坏既有模块）。

## 验证证据

- 新增测试：`provider_test.go`（happy path、跳过未迁移模块、跳过非启用 provider、
  Descriptor mismatch、未声明 key、identity 不匹配、跨模块冲突、navigation 悬空
  permission、fragment capability 缺失、确定性排序、双 Profile 契约矩阵）、
  `persistence_test.go`（happy、version gap、重复 version/name/checksum、tombstone
  with Apply、reconcile 越界/缺 checksum、provider error）。
- `go test ./...` 全量 API 通过；`go vet ./...` 干净。
- fx 静态检查：`internal/kernel/` 与 `internal/modules/` 均无 `go.uber.org/fx`
  import（仅 `composition/` 作为 Fx 组合根允许）。

## 边界

双 Profile 的 register/conflict/Start/Ready 运行时矩阵（冻结 §3）与 readyz 真实
readiness 需真实模块 provider 接线，属 C3/C5 门禁；C2 覆盖契约层 + 最小模块/冲突
测试 + 双 Profile 契约矩阵。业务模块迁移（C3）未开始。
