---
title: 执行记录 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.3.0
---

# 执行记录 · GOAL-006

## 2026-08-02 · 目标立项

- 用户在 Root R3 信息取舍中书面确认推荐方案 B、`features` 菜单投影、两步迁移、读写权限边界和恢复证据口径；Root 记录 D-009 并将 `I-003` 置为 `verified`。
- 从核心五件套约定建立本目标，设定 `parent: GOAL-001-production-admin-foundation`、`status: active` 与六个顺序成功检查点；同步更新工作区 `goal-tree.md`。
- 记录 D-001，选择一个端到端目标承载 R3 强耦合闭环。
- 登记 `I-006-001/002` 两个 required 实施细化项；当前均为 `collecting`，尚未到期但会阻断各自列明的实现门禁。
- **未做**：没有产品代码、数据库、API、Web manifest 或测试行为变更；当前进度为 `0/6`。

## 立项时计划（历史；当时不是事实）

1. 形成版本化 DDL、约束、迁移编号与 seed key 计划，关闭 `I-006-001`。
2. 选择首个真实 `page_ref` / `feature_key` 及 admin/viewer 矩阵，关闭 `I-006-002`。
3. 按 S1 → S6 顺序实施并在每个检查点记录可复现证据。

## 2026-08-02 · 关闭实施前信息门禁（D-002 / D-003）

- 按用户“继续下一步”指令，只读核对当前 `store.Open → migrate → seedAdmin`、`users` / `refresh_tokens` schema、auth/store 兼容签名、records 401/403 路径、真实 manifest、表达式求值、`Session.Features` 与现有测试。
- 落盘 [I-006-001-schema-migration-plan.md](attachments/I-006-001-schema-migration-plan.md)，记录 `0001`/`0002`、精确 DDL/FK/delete/索引、备份恢复、两步读路径切换、稳定 seed 及验证矩阵；记录 D-002，`I-006-001` → `verified`。
- 落盘 [I-006-002-menu-projection-matrix.md](attachments/I-006-002-menu-projection-matrix.md)，记录真实 `list-edit-lifecycle` / `menu_list_edit_lifecycle` 选择、admin/viewer/editor 投影与 API/Web 测试矩阵；记录 D-003，`I-006-002` → `verified`。
- 当前 GOAL-006 无开放 required 信息项，也无正式审计意见或 required finding；S1～S6 均未实施，`status` 保持 `active`，派生进度保持 `0/6`。
- **未做**：未修改 API/Web/manifest、未创建数据库表或备份、未执行迁移与应用测试；附件中的 SQL、函数和测试名均为已冻结的实施计划，不是已发生事实。
- **计划（非事实）**：下一步从 S1 开始实现版本迁移 runner 与可恢复起点；S1 事实和验证齐备前不勾选检查点。

## 2026-08-02 · S1 版本迁移与可恢复起点实施

按用户指令实施 S1（D-002 / I-006-001 §2–§4、§7 V-MIG/V-REC）。新增与修改代码见 `apps/api/internal/store/`。

- **migration runner（`migrate.go` 新增）**：编译期迁移清单 `0001 r2_baseline` / `0002 rbac_expand`；`schema_migrations` 台账（单调整数 version、唯一 name、64 位 SHA-256 checksum、Unix 秒 applied_at）。启动时按连续前缀校验：未知版本、缺失中间版本、name 不匹配、checksum 漂移均 fail closed。每个迁移在单事务内完成 DDL + 数据变换 + 台账插入，任一步失败整体回滚。
- **0001 baseline**：空库在一个 bootstrap 事务内创建 `users`/`refresh_tokens`/`idx_refresh_tokens_user_id` + 台账并登记；既有无台账 R2 库先经 `sqlite_master`/`table_info`/`foreign_key_list`/`index_list` 指纹核对再建台账登记，部分结构/缺列/未知结构回滚且不留下空台账。
- **0002 rbac_expand**：单事务创建 `roles`/`user_roles`/`permissions`/`role_permissions`/`menu_items`/`role_menu_items` 及反向索引；按 `^[a-z][a-z0-9_-]*$` 校验并回填每个 `users.roles`（派生 id `role-<key>`、用户内去重）；非数组、非法 key 或约束冲突使 0002 整体回滚（0001 已提交不受影响）。
- **pre-v0002 恢复快照**：对非空文件库，在应用 0002 前用 `VACUUM INTO` 产生 `<db>.pre-v0002-<UTC>.sqlite`（路径经驱动绑定，不经 shell），并对快照做 `integrity_check=ok`；迁移后主库另跑 `integrity_check=ok` 与 `foreign_key_check` 无违规行。
- **store.go 收敛**：移除旧的 `CREATE TABLE IF NOT EXISTS` `schema` 常量与 `migrate()`；`Store` 新增 `path` 字段；`Open` 签名与 `seedAdmin` 行为不变（启动调用与错误返回契约保持）。
- **证据**：`go test ./internal/store/ -count=1` 12 个用例全过（`TestMigrateFreshDB`、`TestMigrateExistingR2DB`〔含 pre-v0002 快照可新路径打开、原 users/refresh 查询复现〕、`TestMigrateExistingR2DedupeRoles`、`TestMigrateFailClosed{UnknownVersion,MissingIntermediate,ChecksumDrift,PartialBaseline,InvalidRoles}`、`TestForeignKeyEnabled`）；`go test ./...` 全仓 API 通过；`go vet ./...` 与 `gofmt -l` 干净。真实服务冒烟：全新临时库启动后台账 {1,2}、9 张用户表、admin 种子存在、`integrity_check=ok`、`/healthz` 正常。
- **边界**：0002 的 DDL/回填作为 S1 迁移链交付（pre-v0002 快照依赖其存在，D-002 §2.5）；S2 保留阶段 A/B 读路径切换（规范化双写、集合比对、权威读）——本轮未实现。全新库的 `user_roles` 在 seedAdmin 前回填为空属中间态，由 S2 双写与 S3 种子闭合。
- **检查点**：S1 成功标准已勾选；派生进度 `0/6 → 1/6`。S2～S6 尚未实现或审计，`status` 保持 `active`。
