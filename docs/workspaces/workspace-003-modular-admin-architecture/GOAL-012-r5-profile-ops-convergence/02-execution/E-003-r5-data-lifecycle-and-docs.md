---
id: E-003-r5-data-lifecycle-and-docs
doc: execution-entry
goal: GOAL-012-r5-profile-ops-convergence
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-003 · R5 数据生命周期核验与 fork 文档

## R5.2 · fresh/upgrade/recovery 核验（verified；versioned reconcile 为 residual）

- `store/migrate_test.go`：fresh DB 无快照、迁移应用、seed 不覆盖、reopen 幂等；
  pre-v0002/pre-v0004 升级快照、integrity_check、restore、用户/refresh token 保留。
- `operations_test.go`：既有 v3 ledger 升级到 0004（operation_log）+ 可恢复快照。
- `composition_test.go`：MVP reopen 恢复可选模块数据与 core readiness。
- **版本化 system-data reconcile**：当前 seedRBAC 幂等补齐 system keys；独立版本化
  reconcile 路径（冻结 §4.2「版本化、幂等、不覆盖用户字段」的显式版本载体）登记
  R5 residual（R5.2 续作或 R6）。

## R5.4 · fork 文档（partial）

- `README.md` 新增「模块化 Admin 架构与 Profile」段：模块 Provider、Profile/
  MODULES_ENABLED、同一构建双 Profile、数据边界、healthz/readyz、fork 起点。
- Docker Compose 一键启动（既有）+ QUICKSTART 反映新架构待续作。

## 验证

API `go test ./...` + `go vet` + Web `vitest run`（495）通过；README 文本更新无代码
影响。

## 边界

R5.1（Profile 配置收敛 + Configuration 运行时迁移 + 中心适配器删除）、R5.2 versioned
reconcile、R5.4 完整 fork/升级恢复文档待续作。
