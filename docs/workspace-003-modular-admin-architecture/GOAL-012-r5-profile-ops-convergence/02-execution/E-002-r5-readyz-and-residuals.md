---
id: E-002-r5-readyz-and-residuals
doc: execution-entry
goal: GOAL-012-r5-profile-ops-convergence
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R5 readyz 真实 readiness 与 residual 状态

## 已发生事实

- **R5.3 readyz 真实模块图 readiness**：`handler/RegisterWithReadiness` 接受
  `ready func() bool`；composition `readinessGate`（atomic）在 module Start+Ready
  全部成功后置位；readyz 在 store ping OK 且 gate ready 时返回 200，否则 503
  （`not-ready`）。`TestReadyzGatedOnModuleReadiness` 覆盖两态。
- **R4 residual 状态**：
  - Schema 完全 ContributionSet 驱动：accepted-residual（R5 深化 pending）
  - 中心 RegisterSettings/RegisterActivity 适配器终态删除：R5.1 pending（module.go
    死适配器删除被工具权限拦截，登记 R6）
  - PolicyID/Visibility allowlist 深化：R5.1 pending（R4 最小 trim residual）
  - 双 Profile Start/Ready 矩阵：R5.2 pending（R4 收窄 residual）
  - Configuration 运行时迁移：R5.1 pending

## 验证

API `go test ./...`（14 包）+ `go vet` + Web `vitest run`（495）通过。

## 边界

R5.1（Profile 配置收敛 + residual）、R5.2（fresh/reconcile/升级恢复）、R5.4
（Docker/代理/fork 文档）待续作。
