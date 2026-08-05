---
id: E-007-r6-c62-system-data-reconcile
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-007 · R6 C6.2 contribution-driven system-data reconcile（切片 4）

## 已发生事实

- `core.auth-session` 新增独立 `systemdata` owner 包。`Bootstrap` 只在迁移前空库
  `WasFresh()` 为真时创建最小管理员与 system roles；`Reconcile` 在 Provider finalize
  后消费 `ContributionSet.Permissions` / `ContributionSet.Navigation`，两条路径不再由
  `store.Open` 或中心 seed 隐式混合。
- auth-session migration 0009 新增 `system_data_reconcile` 版本/checksum ledger 与
  `system_data_grants` 受管 grant 台账。reconcile 按 module/kind/key/version/checksum
  fail closed，版本降级、同版本 checksum 漂移、identity 冲突均失败并回滚；同版本重放
  不改写 `applied_at`。
- `PermissionContribution` 现在必须声明稳定 `PolicyID` 与正数 `SystemDataVersion`；
  `NavigationContribution` 必须声明同模块 `PageID`、稳定 `Visibility` 与正数版本。
  users/roles/settings/activity Provider 共用 auth-session 的 policy/version 常量，生产
  permissions/menu/grants 的唯一输入是 finalized Provider contributions。
- reconciliation 只创建缺失的系统 permission/menu identity，并只调整
  `system_data_grants` 明确追踪的默认 grant；不覆盖 role name、permission description、
  menu sort/enabled 等 operator-owned 字段，不删除 custom 或 disabled-profile 数据。
  system role 的 `system` 标记为 owner 字段，仅在需要修复时更新，稳定重放不改时间戳。
- composition 在 migration + fresh bootstrap 后构造启用 Provider、finalize contributions、
  执行 reconcile，再将 system-data readiness 纳入 lifecycle Ready；失败关闭 store 并以
  `core.auth-session` lifecycle error fail closed。旧 `store/seed.go`、`seedRBAC` 和
  `Store.seedAdmin` 已删除。
- `testsupport` 仅保留测试夹具贡献，避免业务模块同包测试的反向 import cycle；生产路径
  不使用该夹具。composition 整链测试直接验证真实 Provider：fresh `mvp` 为 5 permissions /
  2 navigation，fresh `admin` 为 8/4；Admin → MVP 后 settings/activity 的 ledger、grant、
  permission 与 menu 数据仍保留。

## 验证

- `apps/api: go test ./...`：通过。
- `apps/api: go vet ./...`：通过。
- `git diff --check`：通过。
- 定向包：`go test ./internal/modules/authsession/systemdata ./internal/kernel ./internal/store ./internal/composition`：通过。
- 生产旧 seed 扫描（排除 `*_test.go`）：`seedRBAC`、`Store.seedAdmin`、
  `store.seedAdmin`、`store.seedRBAC` 均零命中。
- 直接回归覆盖 fresh/replay、用户字段保护、profile downgrade、policy version upgrade、
  checksum drift、version downgrade、identity rollback、ledger/role 时间戳稳定性与 readiness。

## 仍开放

- Root A-010 `F-005`（seed/RBAC 非贡献驱动）已由本切片以 `fixed` 路径闭合。
- C6.2 仍不能勾选：A-010 `F-001` 要求 `internal/store` 中 users/roles/settings/
  operation-log 等领域仓储迁至 owner 模块，store 只保留 DB/tx/runner/ledger 平台职责。
  因此 `F-C62-004` 仅收窄、不关闭；本条不是 VP 退出 #2/#3/#5 的完整证据。
