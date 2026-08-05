---
id: A-004-r4-c2-audit-response
doc: audit-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: self
date: 2026-08-05
scope: Response to Grok A-003 required findings F-IND-C2-001/002/003 and recommended dispositions
verdict: conditional
---

# A-004 · Grok A-003 必改项响应

## required 闭合

| finding | closure | 证据路径 |
|---------|---------|----------|
| F-IND-C2-001 · Descriptor 完全匹配 | `fixed` | `provider.go` `descriptorsMatch` 比对全规范字段（ID/Version/KernelAPIRange/DependsOn/Provides/Requires/ContributionKeys 全键集）；`TestRegisterContributionsDescriptorKernelAPIMismatch` 覆盖 |
| F-IND-C2-002 · Fragments 声明冲突 | `fixed` | `module.go` `validateContributions` 增加 fragment 全局唯一；`TestFragmentsDeclarationConflict` 覆盖 |
| F-IND-C2-003 · C2.2 条文与实现不对齐 | `fixed`（收窄条文） | meta C2.2 改为「compiled-global Persistence **catalog** 静态校验」；ledger 侧 drift/unknown/fresh/upgrade/reconcile 事务路径显式挂 C3/C5（冻结 §4.2） |

## recommended 处置

| finding | 处置 |
|---------|------|
| F-IND-C2-004 · Manifest secrecy | 延至 C3：公开 Manifest 禁 secret 规则 + 测试登记为 C3 门禁 |
| F-IND-C2-005 · Ready 失败反向清理 | 延至 C3/C5：运行时生命周期矩阵；C2 未接 composition 生产路径 |
| F-IND-C2-006 · PolicyID/Visibility/JSON 校验器 | 延至 C3：真实业务数据接线时实现最小校验 |
| F-IND-C2-007 · navigation parent 顺序敏感 | `fixed`：finalize 两遍扫描；`TestNavigationParentOrderIndependent` |
| F-IND-C2-008 · 冲突错误稳定 code | `fixed`：`*kernel.Error` `CodeModuleContributionConflict`；测试断言 code |

## 验证

`go build`/`go test ./internal/kernel/`/全量 `go test ./...`（apps/api）/`go vet`
均通过。fx 静态检查不变。

## 结论

三条 required 已合法闭合，recommended 的 C2-004/005/006 显式登记 C3/C5 门禁。
C2.1-C2.4 检查点成立；GOAL-008 具备关门条件。关门与 GOAL-005 放行 C3 由
`/govern` 在确认后执行。
