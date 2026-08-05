---
id: E-002-r4-c5-verification
doc: execution-entry
goal: GOAL-011-r4-c5-acceptance
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R4-C5 验收验证（C5.1-C5.3）

## C5.1 · 双 Profile 行为矩阵（verified）

- `composition_test.go` `TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan`：
  mvp 与 admin 各跑一轮，断言 settings 路由 404 vs 401、branding 404 vs 200、
  settings/activity schema 页与 manifest 页存在性（mvp 禁用 → 页消失）。
- `TestMVPRecoveryRestoresOptionalModuleDataAndCoreReadiness`：mvp 下 Settings/
  Activity surface 关闭但 `settings.update` 行与数据保留（禁用≠删表）。
- Web `vitest run`（495）与 integration tests：同一构建消费不同 Profile 页面集。

## C5.2 · ledger drift/unknown + 双 Profile 失败矩阵（verified）

- `store/migrate_test.go`：unknown applied version（V-MIG-03）、missing intermediate
  version、checksum drift、partial baseline、roles drift 均 fail closed。
- `TestDualProfileRegisterValidationFailClosed`（composition，新增）：mvp/admin 各跑
  一轮，mismatched provider descriptor → MODULE_API_MISMATCH fail closed。
- `kernel/provider_test.go`：register/conflict/descriptor mismatch fail-closed
  （profile-agnostic）。Start 失败清理由 `TestAppStartFailsClosedWhenPortUnavailable`
  覆盖；Ready 失败 → Stop 由 composition `registerLifecycle` 处理（代码审）。

## C5.3 · C5 收尾

| 项 | 状态 |
|----|------|
| PolicyID/Visibility allowlist 深化 | R4 最小 trim 规则 accepted-residual（GOAL-010 A-003 C4-004）；allowlist 随 R5/R6 |
| 中心 RegisterSettings/RegisterActivity 终态删除 | module.go 死适配器待 R6 删除（handler 测试环境仍走中心，文档化测试路径） |
| Schema owner 完全 ContributionSet 驱动 | accepted-residual（GOAL-010 A-003 C4-001），触发 C5/R6 schema 发布接线 |
| readyz 真实 readiness | 当前 store-ping；真实模块图 readiness 属 R5（冻结 §3 诚实边界） |

## 结论

C5.1/C5.2 证据充分；C5.3 收尾项以 accepted-residual 或文档化登记。C5.4 关门审计
（self + Grok independent）待执行。
