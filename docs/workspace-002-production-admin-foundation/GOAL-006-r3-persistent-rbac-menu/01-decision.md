---
title: 决策 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.3.0
---

# 决策 · GOAL-006

## D-001 · 用一个端到端目标实施 Root D-009

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 以单一 R3 子目标承载版本迁移、规范化 RBAC、增量种子、records 读写授权、`features` 菜单投影和恢复/回归；六个成功标准按依赖顺序推进。
  2. API 权限 key 固定为 `records.read` 与 `records.write`；viewer 仅获得前者，admin 获得两者。权限判断来自持久化 role-permission 关系，不再直接判断 `admin` 角色字符串。
  3. `menu_items` 使用唯一 `page_ref` 和显式唯一 `feature_key`；`/api/accounts/me.features` 输出布尔值，静态 manifest 用 `visibleWhen` 消费。Web 投影只控制展示，API 仍独立强制授权。
  4. 迁移分为“建立/回填/双读核对”和“切换规范化读写”两步；旧 `users.roles` 的删除或停用不纳入本目标的不可逆切换。
  5. R3 恢复口径固定为迁移前副本恢复 + `PRAGMA integrity_check` + 身份/授权/菜单/refresh 关键查询；完整生产备份流程留给 R5。
- **理由**：这些变更共享 schema、seed、身份快照和端到端验证，拆成多个并列目标会制造跨目标中间态与重复门禁；顺序检查点能保留可核对的阶段边界。
- **实施门禁**：`I-006-001` 在首个代码变更前关闭；`I-006-002` 在 S5 前关闭。两项只细化已选模型，不重开 Root D-009 的方案裁决。

### 未选方案

- **按 migration / authorization / menu 拆成三个并列目标**：强耦合中间态难以独立验收，且容易让菜单投影先于真实授权或迁移契约落地。
- **先写代码、再补 DDL/feature key 决策**：会把约束、迁移顺序和真实菜单项选择变成隐式事实，违反 P-005。

## D-002 · 冻结版本迁移、规范化 DDL 与增量种子计划

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. `schema_migrations` 使用单调整数 `version`、唯一 `name`、64 位 SHA-256 `checksum` 与 Unix 秒 `applied_at`；已应用的未知版本、缺号或校验和漂移一律 fail closed。
  2. `0001_r2_baseline` 负责空库创建现有 `users` / `refresh_tokens`，或对无版本台账的既有 R2 库执行结构指纹核对后登记；任何部分结构或不兼容列/索引不静默接管。
  3. `0002_rbac_expand` 在单事务内创建 `roles`、`user_roles`、`permissions`、`role_permissions`、`menu_items`、`role_menu_items`，校验并回填 `users.roles`，再记录迁移版本；精确 DDL、FK/delete 语义与反向索引见 [I-006-001-schema-migration-plan.md](attachments/I-006-001-schema-migration-plan.md)。
  4. 每个连接显式启用并断言 `PRAGMA foreign_keys=ON`。既有文件库应用首个未执行迁移前，以 SQLite 一致性快照产生带目标版本和 UTC 时间的备份；快照与迁移后主库都必须通过 `integrity_check`，主库另跑 `foreign_key_check`。
  5. 两步兼容期：阶段 A 继续以旧 JSON 为对外读值，同时规范化双写并按角色集合比对；阶段 B 以规范化关系为权威读值、按 role key 排序输出，同时继续双写和比对。旧列删除只允许后续显式迁移，本目标不执行。
  6. 稳定 seed 角色为 `admin`、`editor`、`viewer`；permission 为 `records.read`、`records.write`。admin 获读写，viewer 获读，editor 为 R2 兼容仅获读；admin 种子用户保留 `admin` + `editor` 角色。任意既有用户不得使关系 seed 整体跳过。
- **理由**：该计划在不破坏 R2 身份/refresh 契约的前提下提供可诊断升级、可回退快照和稳定授权关系；`editor` 只读避免因角色名称推断而扩大原有写权限。
- **信息门禁**：`I-006-001` → `verified`；证据为附件中的当前代码事实、精确结构、迁移状态机、seed 清单与测试矩阵。本决策只放行 S1/S2/S3 实施，不构成检查点完成。

### 未选方案

- **继续启动时执行一段 `CREATE TABLE IF NOT EXISTS`**：无法区分历史版本、校验漂移或部分升级。
- **在 `0002` 同时删除 `users.roles`**：破坏 D-009 的核对和恢复窗口。
- **把 editor 自动升级为写权限**：现有写门禁只认 admin；该变化会扩大权限且没有用户裁决。
- **仅复制数据库主文件作为备份**：可能遗漏 SQLite 日志状态；采用 SQLite 自身的一致性快照机制。

## D-003 · 首个真实菜单 gate 采用 list-edit-lifecycle

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 首个持久化菜单项固定为 `page_ref: list-edit-lifecycle`、`feature_key: menu_list_edit_lifecycle`、`id: menu-list-edit-lifecycle`；feature key 只用下划线，不使用点号。
  2. 真实 manifest 的该导航子项增加 `visibleWhen.when: "$context.features.menu_list_edit_lifecycle == true"`；页面、路由、标签与导航结构仍由静态 manifest 决定。
  3. `/api/accounts/me.features` 对已登记菜单 key 输出完整布尔投影：admin 为 `true`；viewer/editor 为 `false`。匿名请求仍由认证边界返回 `401`，不生成导航上下文。
  4. `catalog` 保持未加菜单 gate，使 viewer 的 `records.read` 有对应读入口；`list-edit-lifecycle` 含编辑表单且依赖 records 写链，admin-only 投影与 `records.write` 边界一致。
  5. 菜单隐藏只控制导航投影；直接深链不是安全边界，records API 仍独立执行 `records.read` / `records.write`。
- **理由**：`Session.Features` 是 flat `map[string]bool`，而表达式中的点号表示嵌套对象，因此下划线 key 可避免序列化后无法命中；选择真实且已存在 fixture 的编辑生命周期页，同时不剥夺 viewer 的只读入口。
- **信息门禁**：`I-006-002` → `verified`；证据与正反矩阵见 [I-006-002-menu-projection-matrix.md](attachments/I-006-002-menu-projection-matrix.md)。本决策只放行 S5 实施，不构成菜单投影已实现。

### 未选方案

- **`admin.catalog` 等带点号 key**：flat Go map 与 Web 点路径求值语义不一致。
- **对 `catalog` 做 admin-only gate**：会让 viewer 拥有 `records.read` 却失去明确读入口。
- **选择 `settings` 或 `activity`**：manifest 有引用但当前后端没有对应 checked-in schema fixture，不能作为首个端到端证据。
- **选择 `overview`**：它是 `homePageRef`；首个 gate 会引入根路由/首页可达性问题。

## D-004 · S1 交付包含 0002 迁移链；S2 独有交付为读路径切换

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. S1 的交付物包括迁移 runner、`schema_migrations`、`0001 r2_baseline`、`0002 rbac_expand`（DDL + 角色回填）以及 pre-v0002 恢复快照与验证。
  2. S2 的独有交付物是阶段 A/B 读路径切换：`CreateUser` 规范化双写、`UserByID/UserByUsername` 集合比对与规范化权威读、按 role key 升序输出；旧 JSON 列保留至后续显式迁移。
  3. 全新库在 `seedAdmin` 之前完成 0002 回填，`user_roles` 为空属预期中间态，由 S2 双写与 S3 增量种子闭合。
- **理由**：成功标准 S1 的"升级前可恢复数据库副本"需要一个真实的数据变换迁移来产生与验证快照；D-002 §2.5 将快照固定命名为 `pre-v0002`，因此 0002 必须已在迁移链内。`0002` 只在 0001 之后运行，不影响空库 bootstrap；0002 不涉及任何对外读路径，未抢先占用 S2 的范围。
- **未选方案**：**S1 只交付 runner + 0001，0002 推迟到 S2**——pre-v0002 快照将无可验证的迁移目标，S1 的可恢复起点标准无法在代码层面达成。
