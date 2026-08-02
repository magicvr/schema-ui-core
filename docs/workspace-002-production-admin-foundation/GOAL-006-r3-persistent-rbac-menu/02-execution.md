---
title: 执行记录 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.6.0
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
- **pre-v0002 恢复快照**：对非空文件库，在应用 0002 前用 `VACUUM INTO` 产生 `<db>.pre-v0002-<UTC>.sqlite`（驱动内 SQL 字面量转义——单引号转义、不经 shell），并对快照做 `integrity_check=ok`；迁移后主库另跑 `integrity_check=ok` 与 `foreign_key_check` 无违规行。
- **store.go 收敛**：移除旧的 `CREATE TABLE IF NOT EXISTS` `schema` 常量与 `migrate()`；`Store` 新增 `path` 字段；`Open` 签名与 `seedAdmin` 行为不变（启动调用与错误返回契约保持）。
- **证据**：`go test ./internal/store/ -count=1` 12 个用例全过（`TestMigrateFreshDB`、`TestMigrateExistingR2DB`〔含 pre-v0002 快照可新路径打开、原 users/refresh 查询复现〕、`TestMigrateExistingR2DedupeRoles`、`TestMigrateFailClosed{UnknownVersion,MissingIntermediate,ChecksumDrift,PartialBaseline,InvalidRoles}`、`TestForeignKeyEnabled`）；`go test ./...` 全仓 API 通过；`go vet ./...` 与 `gofmt -l` 干净。真实服务冒烟：全新临时库启动后台账 {1,2}、9 张用户表、admin 种子存在、`integrity_check=ok`、`/healthz` 正常。
- **边界**：0002 的 DDL/回填作为 S1 迁移链交付（pre-v0002 快照依赖其存在，D-002 §2.5）；S2 保留阶段 A/B 读路径切换（规范化双写、集合比对、权威读）——本轮未实现。全新库的 `user_roles` 在 seedAdmin 前回填为空属中间态，由 S2 双写与 S3 种子闭合。
- **检查点**：S1 成功标准已勾选；派生进度 `0/6 → 1/6`。S2～S6 尚未实现或审计，`status` 保持 `active`。

## 2026-08-02 · S2 规范化 RBAC 两阶段兼容实施

响应 A-001 后按用户确认「直接交付 B 终态」实施 S2（D-002.5 阶段 A/B；I-006-001 §5、§7 V-MIG-05）。改动见 `apps/api/internal/store/`。

- **读路径（阶段 B 终态）**：`UserByID/UserByUsername` 经 `userWithRoles` 读取——取 legacy JSON 与 join `user_roles`+`roles` 的规范化角色做集合比对，不一致返回可诊断错误（双读核对，fail loudly）；一致时以规范化关系为权威读值、按 role key 升序输出。`account.User {id,name,roles}` 形状、`refresh_tokens.user_id`、JWT subject 与密码字段均不变；旧 `users.roles` 列保留（删除属后续显式迁移）。
- **写路径（阶段 A/B 双写）**：`CreateUser` 单事务写 legacy JSON 与 `user_roles`（输入 roles 先去重，保证两源集合一致）；`seedAdmin` 幂等确保 admin+editor 关系，闭合 D-004 点出的 S1 全新库 `user_roles` 空中间态，且不覆盖已有密码。
- **派生 role 自建**：新增 `ensureRole`/`linkUserRole` 共享助手，按需幂等创建 `role-<key>`（system=0），供 0002 回填、CreateUser、seedAdmin 共用，避免 S3 种子前 `user_roles` FK 悬空；0002 `backfillRoles` 重构复用该助手。
- **F-002（S2 部分）闭合**：`TestUserRolesFKAndCascade` 断言 `user_roles` FK（未知 role 拒绝）、RESTRICT（在用 role 删除拒绝）、CASCADE（删用户级联清理）；完整 V-MIG-04 unique/CASCADE|RESTRICT/反向索引矩阵仍留 S6。
- **证据**：新增 `normalize_test.go` 5 用例（`TestCreateUserDoubleWritesRoles`、`TestNormalizedReadSortedByKey`、`TestReadDetectsRoleMismatch`、`TestSeedAdminDoubleWrites`、`TestUserRolesFKAndCascade`）全过；`go test ./...` 全仓 API 通过（auth/account/handler R2 契约保持）；`go vet ./...` 与 `gofmt -l` 干净。真实服务冒烟：全新库 `admin/admin` 登录成功，`user.roles=["admin","editor"]` 来自规范化源（升序），seed `user_roles=2 roles=2`（gap 闭合），`integrity_check=ok`。
- **边界**：S2 只切读路径与双写；权限/菜单 grant 与稳定角色升级（system=1）属 S3 增量种子，本轮未实现。
- **检查点**：S2 成功标准已勾选；派生进度 `1/6 → 2/6`。S3～S6 尚未实现或审计，`status` 保持 `active`。

## 2026-08-02 · A-002 F-004 fixed：读路径集合比较修正 + 迁移后读取/认证回归

响应 A-002（independent · `conditional` · F-004 required），按用户裁决 **fixed**。改动见 `apps/api/internal/store/` 与 `apps/api/internal/auth/`。

- **根因**：0002 `backfillRoles` 对同一用户 role key 去重，但 `store.go` 的 `userWithRoles` 比较按多重集合（长度+计数），导致历史 R2 `roles` JSON 含重复值（如 `["admin","admin","editor"]`）的已迁移用户两源长度不一致而被误判为分歧，`UserByID` / `UserByUsername` 无法加载身份。
- **修正**：`userWithRoles` 改用 `sameRoleSet` **集合语义**比较（忽略顺序与重复），对齐 I-006-001 §5 / D-005 已冻结语义；返回仍以规范化关系为权威、按 role key 升序。真正的集合分歧仍 fail closed。
- **回归证据**：`TestMigrateExistingR2DuplicateRolesReadable`（`migrate_test.go`：重复 legacy roles → Open → `UserByID`/`UserByUsername` 可读且返回 `["admin","editor"]`）；`TestLoginAndRefreshAfterMigrateDuplicateRoles`（`auth_test.go`：迁移后 `Login` 走 `UserByUsername`、`Refresh` 走 `UserByID` 均成功，access subject=`u-alice`）。
- **验证**：`go test ./...`（apps/api）全绿、`go vet ./...` 干净、`gofmt -l` 无输出。
- **检查点**：F-004 为 S2 内 required 修正，不构成新检查点；派生进度保持 `2/6`，`status` 保持 `active`。

## 2026-08-02 · S3 增量幂等种子实施

A-004 自审（S1/S2 事实 + F-004 闭合 + S3 门禁）`pass` 后，按用户确认的接线方案实施 S3（D-002.6 / I-006-001 §6，V-SEED-01）。改动见 `apps/api/internal/store/`。

- **`seed.go` 新增 `seedRBAC`**：幂等 ensure 稳定实体与关系——roles `admin`/`editor`/`viewer`（升级 `system=1`，upsert name/system）、permissions `records.read`/`records.write`（`perm-records-read`/`perm-records-write`）、菜单项 `list-edit-lifecycle`（`menu-list-edit-lifecycle` / `menu_list_edit_lifecycle`）、grants（admin → read+write+menu；editor、viewer → 仅 read）。
- **Open 接线**：`seedAdmin=true` 时在 `migrate()` + `seedAdmin()` 之后运行 `seedRBAC()`（D-006，用户确认）。生产 `main.go` 恒为 true，服务启动自愈；`Open(seedAdmin=false)` 既有语义不变。任意既有用户不使关系 seed 整体跳过；重复启动无重复；不覆盖非种子用户字段；editor 不升级为写。
- **验证**：`TestSeedRBACEntitiesAndGrants`（实体存在、system=1、admin 2 项 permission + 1 菜单、editor/viewer 仅 read、records.write 仅 admin）；`TestSeedRBACIncrementalWithExistingUsers`（先 `seedAdmin=false` 建既有用户 → 重开 `seedAdmin=true` 关系补齐、非种子用户字段不变、三次启动无重复）。`go test ./...` 全仓 API 通过；`go vet ./...` 与 `gofmt -l` 干净。真实服务冒烟：全新库 `admin/admin` 登录成功，`roles(system=1)=3`、`permissions=2`、`menu_items=1`、`role_permissions=4`、`role_menu_items=1`、`integrity_check=ok`。
- **边界**：S3 只建稳定种子关系；后端读写授权 gate（permission 判断）属 S4，本轮未实现；menu 投影与 `/api/accounts/me.features` 属 S5。
- **检查点**：S3 成功标准已勾选；派生进度 `2/6 → 3/6`。S4～S6 尚未实现或审计，`status` 保持 `active`。
