---
id: E-003-r4-c2-audit-response
doc: execution-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-003 · Grok 审计必改项响应与修复

## 响应（Grok A-003）

| finding | closure | 证据路径 |
|---------|---------|----------|
| F-IND-C2-001 · Descriptor 完全匹配仅 ID+Version | `fixed` | `provider.go` `descriptorsMatch`：比对 ID/Version/KernelAPIRange/DependsOn/Provides/Requires/ContributionKeys 全键集；新增 KernelAPIRange mismatch 测试 |
| F-IND-C2-002 · Fragments 未纳入 Plan 声明冲突 | `fixed` | `module.go` `validateContributions` 增加 fragment 全局唯一校验；新增 `TestFragmentsDeclarationConflict` |
| F-IND-C2-003 · C2.2 成功标准与实现不对齐 | `fixed`（收窄条文） | meta C2.2 收窄为 catalog 静态校验；ledger 侧 drift/unknown/fresh/upgrade/reconcile 事务路径明确挂 C3/C5 |
| F-IND-C2-004 · Manifest secrecy 未实现（recommended） | 延至 C3 | 登记 C3/C5 门禁（公开 Manifest 禁 secret 规则 + 测试），不改 C2 勾选措辞 |
| F-IND-C2-005 · Ready 失败不反向 Stop（recommended） | 延至 C3/C5 | 运行时生命周期矩阵属 C3/C5；C2 未接 composition 生产路径 |
| F-IND-C2-006 · PolicyID/Visibility/JSON 无校验器（recommended） | 延至 C3 | 真实业务数据校验器在 C3 接线时实现；C2 契约类型已含字段 |
| F-IND-C2-007 · navigation parent 引用顺序敏感（recommended） | `fixed` | `provider.go` finalize 改两遍扫描（先收齐 NodeID 再查 parent）；新增 `TestNavigationParentOrderIndependent` |
| F-IND-C2-008 · 冲突错误无稳定 error code（recommended） | `fixed` | `contributionConflictError` → `*kernel.Error` `CodeModuleContributionConflict`；测试断言 code |

## 验证

- `go build ./internal/kernel/` ok；`go test ./internal/kernel/` 全通过（含新增
  KernelAPIRange mismatch、Fragments 声明冲突、navigation 顺序、conflict code）。
- 全量 `go test ./...`（apps/api）与 `go vet ./...` 通过。
- fx 静态检查不变（kernel/modules 无 Fx import）。

## 边界

C2 主体（契约层 + 冲突/最小模块测试 + 双 Profile 契约矩阵 + fx 静态检查）完成。
Manifest secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、ledger 侧
drift/unknown 与运行时双 Profile 矩阵登记为 C3/C5 门禁，不阻断 C2 关门。
