---
title: 审计台账 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.9.0
---

# 审计台账 · GOAL-006

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-02 | S1 execution-facts + 目标定义/信息门禁 | pass | responded（recommended：F-001/F-002 deferred → S2/S6，F-003 fixed） |
| A-002 | independent | 2026-08-02 | S2 execution-facts + R2 兼容性 + A-001 finding closure review | conditional | responded（F-004 required → fixed；复审建议：/audit 复核或 self 阶段审计） |
| A-003 | independent | 2026-08-02 | finding-closure：F-004 | pass | responded（F-004 fixed 关闭证据独立复核通过） |
| A-004 | self | 2026-08-02 | S1/S2 事实 + F-004 闭合 + S3 门禁（stage 自审） | pass | responded（S3 门禁就绪；F-101 recommended 跟踪 → S6） |

## 当前审计边界

- A-001 覆盖：目标定义与 D-001～D-004、`I-006-001/002` 门禁状态，以及 **S1**（版本迁移与可恢复起点）执行事实对照。
- A-002 覆盖：**S2**（规范化权威读、双写、集合核对与 R2 身份/refresh 兼容）执行事实，以及 A-001 中 F-001/F-002/F-003 的响应状态；不审 S3～S6 或 Root R3 阶段关门。
- A-002 的 F-004（required）已由 `/govern` 响应为 **fixed**（集合语义比较 + 迁移后读取/认证回归）；A-002 verdict 仍为 `conditional`，其 required 门禁已合法闭合。
- A-003 独立复核 F-004 关闭证据：`sameRoleSet` 的集合语义、迁移后双 lookup、Login/Refresh 回归均已在当前工作树复核通过；本条不扩展至 S3～S6。
- `I-006-001/002` 保持 `verified`；当前无到期未关的 required 信息项。
- 审计意见不直接修改 `status` / `progress`；响应和推进由 `/govern` 与用户裁决维护。

---

## A-001 · S1 执行事实与目标定义交叉审计（2026-08-02）

- **source**：independent
- **auditor**：GitHub Copilot（Grok 4.5）· `/audit`
- **类型**：execution-facts（主）+ goal-definition / info-gate（辅）
- **scope**：`workspace-002-production-admin-foundation` / `GOAL-006-r3-persistent-rbac-menu` · S1 检查点主张；D-001～D-004 与 `I-006-001/002`；不审 S2～S6 产品行为
- **verdict**：**pass**
- **工作区**：`workspace-002-production-admin-foundation`；`root_goal=GOAL-001-production-admin-foundation`；`canonical_scope` 匹配；`shared_materials_catalog: none`（本审未将共享资料作证据）

### 范围与区间

| 项 | 结论 |
|----|------|
| 工作区绑定 | 显式工作区已校验；目标 `parent` 指向 Root；未读取其他工作区上下文 |
| 信息门禁 | `I-006-001`（S1/S2/S3 实施前）与 `I-006-002`（S5 前）均为 `verified`，附件与 D-002/D-003 可互指；无开放 required |
| S1 成功标准 | 「schema_migrations + 顺序/校验和 + 事务化 fail-closed + 升级前可恢复副本」 |
| 代码证据根 | `apps/api/internal/store/migrate.go`、`migrate_test.go`、`store.go`、`store_test.go` |
| 文档证据 | `00-meta` S1 勾选与 `1/6`；`01-decision` D-001～D-004；`02-execution` S1 时间线；附件 I-006-001/002 |

### 成果（有证据）

1. **版本台账与 runner**：`schema_migrations(version,name,checksum,applied_at)`；编译期 `0001 r2_baseline` / `0002 rbac_expand`；checksum = 规范化 SQL + `transformID` 的 SHA-256 小写 hex；未知版本 / 缺号（台账不自 1 起或中间断层）/ name 不符 / checksum 漂移均 fail closed（`validateApplied`）。
2. **事务化迁移**：`applyMigration` 单事务执行 `up` + 台账插入；0002 非法 roles 时 0001 保留、0002 整步回滚（`TestMigrateFailClosedInvalidRoles`）。
3. **0001 指纹**：空库建 R2 表+台账；既有 R2 经 table/column/FK/index 指纹后登记；部分基线拒绝且不留空台账（`TestMigrateFailClosedPartialBaseline`）。
4. **0002 DDL + 回填**：六张 RBAC 表及反向索引与 I-006-001 §4 / D-002 一致；role key 正则与用户内去重（`TestMigrateExistingR2DedupeRoles`）。
5. **pre-v0002 快照**：非空文件库在首个 `version>=2` 待迁移前 `VACUUM INTO <db>.pre-v0002-<UTC>.sqlite`，快照 `integrity_check=ok`；空库不产生快照；迁移后主库 `integrity_check` + `foreign_key_check`（`TestMigrateExistingR2DB` / `verifyIntegrity`）。
6. **Open 契约**：`Open(path, adminUsername, adminPasswordHash, seedAdmin)` 签名保留；旧 `CREATE TABLE IF NOT EXISTS` 路径已移除；`seedAdmin` 仍仅写 `users` JSON（S2/S3 范围，D-004 已界定中间态）。
7. **回归（本轮独立复跑，2026-08-02）**：
   - `go test ./internal/store/ -count=1` → PASS（含 9 个迁移相关 + 3 个既有 store 用例）
   - `go test ./... -count=1`（`apps/api`）→ PASS
   - `go vet ./...` 干净；`gofmt -l ./internal/store/` 无输出
8. **目标定义**：单一端到端子目标 + 六顺序检查点；D-004 明确 S1 含 0002 链与 pre-v0002、S2 独有读路径切换，避免范围抢跑；进度 `1/6` 与仅 S1 勾选一致。

### 对照成功标准（S1）

| 标准要素 | 判定 | 证据 |
|----------|------|------|
| `schema_migrations` | 满足 | DDL + 台账读写 + fresh/existing 测试 |
| 顺序 / 校验和检查 | 满足 | `validateApplied` + Unknown/Missing/Checksum 测试 |
| 事务化 fail-closed | 满足 | `applyMigration` 回滚 + InvalidRoles/PartialBaseline |
| 升级前可恢复副本 | 满足 | `snapshotPreV0002` + ExistingR2 快照查询复现 |
| 未越权勾选 S2～S6 | 满足 | meta 仅 S1；读路径/seed grants/授权/菜单仍属后续 |

### Findings

#### F-001 · recommended · 低

- **主题**：`TestMigrateFailClosedMissingIntermediate` 通过删除 version=1、仅留 version=2 触发「台账不自 1 起」，并未构造真正的中间缺号（如已应用 1 与 3）。
- **证据**：`migrate_test.go` `TestMigrateFailClosedMissingIntermediate`；对比 `validateApplied` 中 `a.version != applied[i-1].version+1` 分支缺少对应用例。
- **影响**：fail-closed 行为仍被覆盖；「缺中间版本」分支的回归保护偏弱。
- **建议**：S6 或下一迁移相关改动前，补一条 ledger=`(1,3)`（或等价）的用例。

#### F-002 · recommended · 低

- **主题**：I-006-001 `V-MIG-04` 中 unique / ON DELETE CASCADE|RESTRICT / 反向索引存在性未在 S1 测试中逐项覆盖；当前仅有 `TestForeignKeyEnabled`（refresh_tokens → 缺失 user）。
- **证据**：`migrate_test.go`；`rbacExpandDDL` 已声明约束但无对应正反测试。
- **影响**：不否定 S1 成功标准（标准未要求完整 V-MIG-04 矩阵）；S2～S4 开始依赖这些删除/唯一语义时证据不足。
- **建议**：在 S2/S3 或 S6 回归中补 RBAC 表约束与索引断言，勿把「DDL 已写出」当作已验证行为。

#### F-003 · recommended · 低

- **主题**：`02-execution` 写「路径经驱动绑定」；实现为 SQL 字符串字面量 + 单引号转义（`VACUUM INTO '…'`），非 `?` 参数绑定。不经 shell，路径来自 `s.path`，安全性方向正确，表述略过满。
- **证据**：`migrate.go` `snapshotPreV0002`；`02-execution.md` S1 节。
- **影响**：无功能缺陷；后续若有路径注入审查，应以实现为准修正行文。
- **建议**：`/govern` 响应执行记录时改为「驱动内 SQL 字面量转义，不经 shell」。

### 必改项汇总

- **required / 必改**：无
- **recommended**：F-001、F-002、F-003（均不阻断 S1 事实认定，也不阻断进入 S2 实施；F-002 建议在依赖 RBAC 约束的阶段前闭合）

### 与既有意见的异同

- 本目标此前无正式 A-00N 意见（索引原为「尚无正式审计意见」）。本条为首条 `independent` 意见。

### 结论 + 建议给编排器/用户的下一步

- **结论**：在声明的 S1 scope 内，实现、测试与文档主张一致；S1 勾选与 `progress: 1/6` **有可重复证据支撑**；信息门禁关闭方式合规；未发现 high/required 缺口或名不副实的完成声明。
- **建议 `/govern`**：
  1. 响应 A-001（可记录 recommended 的接纳/延期）。
  2. 在用户确认后推进 **S2**（阶段 A/B 读路径切换、规范化双写与集合比对）——注意全新库 `seedAdmin` 后 `user_roles` 仍可能为空，属 D-004 已承认的中间态，应由 S2/S3 闭合。
  3. **不要**将本 pass 解读为 Root R3 或 S2～S6 放行。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。  
响应、finding 闭合与阶段推进由用户通过 **`/govern`** 处理。

---

## 响应 A-001（2026-08-02 · /govern）

- **模式**：response
- **响应意见**：A-001（independent · verdict `pass` · 无 required）
- **裁决**：采纳 `pass` 结论与证据认定；三项 recommended 按用户裁决处置如下。

### 关闭证据表

| Finding | 状态 | 处置 / 证据路径 |
|---------|------|-----------------|
| F-001 · recommended · 低 | **deferred**（open 跟踪） | 用户裁决延期至 **S6（或下一次迁移相关改动）** 前补 ledger=`(1,3)` 缺中间版本用例，增强 `validateApplied` 该分支的回归保护；不阻断 S1 事实认定与 S2 实施。 |
| F-002 · recommended · 低 | **deferred**（open 跟踪；S2 部分已闭合） | 用户裁决延期：**S2** 已补齐 `user_roles` FK / RESTRICT / CASCADE 正反断言（`TestUserRolesFKAndCascade`，`normalize_test.go`）；完整 V-MIG-04 的 unique / CASCADE\|RESTRICT / 反向索引矩阵仍在 **S6** 回归闭合。不把「DDL 已写出」当作已验证行为。 |
| F-003 · recommended · 低 | **fixed** | `02-execution.md` S1 节表述已由「路径经驱动绑定，不经 shell」修正为「驱动内 SQL 字面量转义（单引号转义），不经 shell」；实现（`snapshotPreV0002`）与行文一致。 |

### 仍开放项

- F-001 保持 open 跟踪至 **S6（或下一次迁移相关改动）**。
- F-002 的 **S2 部分已闭合**（`user_roles` FK / RESTRICT / CASCADE 断言，`TestUserRolesFKAndCascade`）；完整 V-MIG-04 矩阵仍 open 至 **S6**。
- 均为 `recommended`，不构成 required 门禁，不阻断推进。
- 无 high / required 开放项。

### 冲突裁决

- 无冲突：本目标此前无正式意见，A-001 为唯一 independent 条目；无同范围相反 verdict。

---

## A-002 · S2 执行事实与 R2 兼容性交叉审计（2026-08-02）

- **source**：independent
- **auditor**：GitHub Copilot · `/audit`
- **类型**：execution-facts（主）+ finding-closure（辅）
- **scope**：`workspace-002-production-admin-foundation` / `GOAL-006-r3-persistent-rbac-menu` · S2 检查点主张、D-005、`I-006-001` 的两步读写要求、R2 身份/refresh 兼容，以及 A-001 的 F-001/F-002/F-003 响应；不审 S3～S6、权限/菜单产品行为或 Root R3 关门
- **verdict**：**conditional**
- **工作区**：`workspace-002-production-admin-foundation`；`root_goal=GOAL-001-production-admin-foundation`；`canonical_scope` 匹配；`shared_materials_catalog: none`（本审未将共享资料作证据）

### 范围与区间

| 项 | 结论 |
|----|------|
| 工作区与愿景绑定 | 显式工作区、Root、`VP-002-production-admin-foundation` 与 Charter 引用链可解析；目标 `parent` 指向本区 Root |
| 信息门禁 | `I-006-001` 为 `verified`，是 S2 实施输入；未发现影响 S2 的到期 required 信息项 |
| S2 成功标准 | 规范化 RBAC 两步兼容：双写、双读核对、规范化权威读，且保持 R2 身份与 refresh 契约 |
| 代码证据根 | `apps/api/internal/store/store.go`、`migrate.go`、`normalize_test.go`、`migrate_test.go`；认证调用链为 `apps/api/internal/auth/auth.go` |
| 本轮验证 | 2026-08-02 独立运行：`go test ./internal/store -count=1`、`go test ./... -count=1` 与 `go vet ./...`（工作目录 `apps/api`）均通过 |

### 成果（有证据）

1. **阶段 B 读路径**：`UserByID` 与 `UserByUsername` 均经 `userWithRoles` 读取 legacy JSON 与规范化关联；一致时采用 join 关系按 `roles.key` 升序的输出，分歧返回可诊断错误。
2. **事务双写**：`CreateUser` 在同一事务写入 `users.roles` 和 `user_roles`；输入角色先去重。`seedAdmin` 对新建或既有 admin 确保 `admin` / `editor` 关联而不覆盖密码。
3. **失配与约束证据**：`TestReadDetectsRoleMismatch` 覆盖两种 lookup 的 fail-closed 行为；`TestUserRolesFKAndCascade` 覆盖 `user_roles` 的未知 role FK 拒绝、在用 role RESTRICT 与删用户 CASCADE，满足 A-001 F-002 约定的 S2 部分闭合。
4. **R2 主链回归**：认证的 Login、Refresh、JWT middleware 都由 `UserByUsername` / `UserByID` 取得身份；本轮全仓 API 测试通过，现有种子和新建用户的身份、refresh 生命周期未见回归。
5. **A-001 跟踪项**：F-003 的执行记录措辞已与实际 SQL 字面量转义实现对齐；F-001 仍按用户记录延期至 S6（或下一次迁移相关改动），F-002 完整约束/反向索引矩阵仍延期至 S6，均保持 recommended，不构成当前 required 门禁。

### 对照成功标准（S2）

| 标准要素 | 判定 | 证据 |
|----------|------|------|
| 规范化关系与回填 | 满足 | `0002` 建立关系表；`backfillRoles` 从 R2 JSON 建立 role / user-role 关系 |
| 规范化权威读、确定排序 | 满足 | `rolesForUser` 按 key 排序；`TestNormalizedReadSortedByKey` |
| CreateUser / seedAdmin 双写 | 满足 | `CreateUser` / `seedAdmin`；`TestCreateUserDoubleWritesRoles`、`TestSeedAdminDoubleWrites` |
| 两源分歧 fail-closed | 满足 | `userWithRoles`；`TestReadDetectsRoleMismatch` |
| R2 身份与 refresh 兼容 | **不满足** | F-004：重复角色的既有 R2 数据被迁移去重后会被比较器拒绝 |

### Findings

#### F-004 · required · 中

- **主题**：历史 R2 `roles` JSON 含重复值时，迁移去重与 S2 读取比较的语义不一致，造成已迁移用户无法加载身份。
- **证据**：`backfillRoles` 在 `migrate.go` 中明确对同一用户 role key 去重；`TestMigrateExistingR2DedupeRoles` 以 `["admin","admin","editor"]` 验证该迁移结果。随后 `store.go` 的 `setEqual` 以 slice 长度和计数比较多重集合，而非 D-005 / `I-006-001` 所要求的集合比较。因此该 fixture 升级后 legacy 为三项、规范化关系为两项，`UserByID` / `UserByUsername` 返回 role mismatch。
- **影响**：认证 Login 依赖 `UserByUsername`；Refresh 和 Bearer middleware 依赖 `UserByID`。受影响的合法历史用户会在迁移成功后无法登录、刷新或通过请求身份加载，违反 S2 的 R2 身份/refresh 兼容主张；现有去重测试仅断言关系计数，未覆盖升级后的读路径。
- **影响门禁**：S2 事实无条件认定与进入 S3。
- **建议修正**：令两源比较遵循已冻结的集合语义（例如比较前分别去重），或在受控迁移中同步规范化旧 JSON；补充“重复 legacy roles → Open → UserByID / UserByUsername 可读且角色按 key 排序”的迁移回归，并至少覆盖 Login 或 Refresh 其中一条受影响调用链。

### 必改项汇总

- **required / 必改**：F-004。修复并留下可复核测试证据前，不得把 S2 作为无条件通过，也不得推进 S3；若拟接受残余或驳回，须由用户按 P-003 / P-004 书面裁决并留痕。
- **recommended**：F-001（S6 或下次迁移改动前补真正中间缺号 ledger 用例）与 F-002 的完整 V-MIG-04 矩阵仍按 A-001 响应跟踪；F-003 已 fixed。

### 与既有意见的异同

- A-001 的 `pass` 仅覆盖 S1，且明确不审 S2；本条不与 A-001 冲突。
- A-001 F-002 的 S2 部分关闭证据有效，但它只涉及 `user_roles` FK / RESTRICT / CASCADE，不能覆盖 F-004 的读路径兼容问题。

### 结论 + 建议给编排器/用户的下一步

- **结论**：S2 的主要实现路径、约束子集和既有回归具备证据，但一个由现有迁移测试明确承认的历史数据形态会破坏身份兼容，因此 verdict 为 `conditional`，不能无条件认可 S2 完成主张。
- **建议 `/govern`**：先响应 F-004。建议选择 `fixed`：修正集合比较或受控数据规范化，增加迁移后读路径与认证链回归，再请求复审；在此之前不要推进 S3。若不修复而要继续，需用户书面选择 `accepted-residual` 或 `user-overruled`，并记录影响用户范围、期限、缓解和复审触发。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。响应、finding 闭合与阶段推进由用户通过 **`/govern`** 处理。

---

## 响应 A-002（2026-08-02 · /govern）

- **模式**：response
- **响应意见**：A-002（independent · verdict `conditional` · required F-004）
- **裁决**：按用户裁决选择 **fixed** 闭合 F-004——修正读路径集合比较语义，并补充迁移后读取与认证回归测试。

### 关闭证据表

| Finding | 状态 | 处置 / 证据路径 |
|---------|------|-----------------|
| F-004 · required · 中 | **fixed** | `store.go` `userWithRoles` 的两源比较由多重集合（长度+计数）改为 `sameRoleSet` **集合语义**（忽略顺序与重复），对齐 I-006-001 §5 / D-005 已冻结语义；历史 R2 重复 role key 经 0002 回填去重后不再被误判为分歧。证据：① `TestMigrateExistingR2DuplicateRolesReadable`（`migrate_test.go`）：`["admin","admin","editor"]` → Open → `UserByID` / `UserByUsername` 可读且返回去重升序 `["admin","editor"]`；② `TestLoginAndRefreshAfterMigrateDuplicateRoles`（`auth_test.go`）：迁移后 `Login`（走 `UserByUsername`）与 `Refresh`（走 `UserByID`）均成功，access subject=`u-alice`。`go test ./...`、`go vet ./...`、`gofmt -l` 全绿。 |

### 仍开放项

- F-004 已按 `fixed` 合法闭合；S2 无条件事实认定与进入 S3 的 required 门禁解除。
- F-001 保持 open 跟踪至 **S6（或下一次迁移相关改动）**；F-002 的完整 V-MIG-04 矩阵保持 open 至 **S6**（S2 部分已闭合）。均为 `recommended`。
- 无未合法闭合的 required 开放项。

### 冲突裁决

- A-001（pass，仅覆盖 S1）与 A-002（conditional，覆盖 S2）无同范围相反 verdict；A-002 明确不推翻 A-001 的 S1 结论。F-004 修复未改变 D-002/D-005 的方案正文，不构成决策冲突。

### 复审建议

- A-002 建议「修复后请求复审」：可由 `/audit` 复核 F-004 关闭证据，或由 `/govern` 在推进 S3 前视需要补一次 `self` 阶段审计（GOAL-006 目前无 self 意见）。

---

## A-003 · F-004 关闭证据独立复核（2026-08-02）

- **source**：independent
- **auditor**：GitHub Copilot · `/audit`
- **类型**：finding-closure
- **scope**：`workspace-002-production-admin-foundation` / `GOAL-006-r3-persistent-rbac-menu` · 仅复核 A-002 的 `F-004`（required，中）关闭证据；不审 S3～S6、其他 recommended finding 或 Root R3 关门
- **verdict**：**pass**
- **工作区**：`workspace-002-production-admin-foundation`；`root_goal=GOAL-001-production-admin-foundation`；`canonical_scope` 匹配；`shared_materials_catalog: none`（本审未将共享资料作证据）

### 范围与区间

| 项 | 结论 |
|----|------|
| 关闭路径 | 用户已在 A-002 响应中选择 `fixed`，并记录了修正、测试和验证命令；本复核仅判断这些证据是否真实、充分、可重复 |
| F-004 原主张 | 历史 R2 重复角色经 0002 去重后，S2 读路径不得误判角色分歧并阻断身份加载 |
| 信息门禁 | `I-006-001` 仍为 `verified`；本复核没有发现新的 required 信息项或 residual 裁决需求 |
| 审计范围 | 仅涉及 `store.go`、`migrate_test.go`、`auth_test.go` 与对应执行记录；不将 `progress: 2/6` 当作关闭依据 |

### 成果（有证据）

1. **根因修正已存在**：`userWithRoles` 现在调用 `sameRoleSet`；该函数以 map 比较角色集合，忽略 legacy JSON 中的重复值和顺序，同时仍能识别真实集合差异。规范化关联仍作为权威输出，并由 SQL 按 role key 升序返回。
2. **迁移后双 lookup 覆盖**：`TestMigrateExistingR2DuplicateRolesReadable` 构造 `roles=["admin","admin","editor"]` 的既有 R2 数据，执行 `Open` 后同时验证 `UserByID` 和 `UserByUsername` 成功，并返回去重且升序的 `["admin","editor"]`。
3. **认证链覆盖**：`TestLoginAndRefreshAfterMigrateDuplicateRoles` 使用迁移后的重复角色用户验证 Login（`UserByUsername`）与 Refresh（`UserByID`），并核对 access token subject 为迁移用户 id。
4. **独立复跑结果**：聚焦 store 测试、聚焦 auth 测试、`go test ./... -count=1`、`go vet ./...` 和 `gofmt -l ./internal/store ./internal/auth` 均通过；格式检查输出 `gofmt clean`。
5. **记录一致性**：`02-execution.md` 与 A-002 响应均指向同一 `sameRoleSet` 修正和两条新增回归；`00-meta` 与 `goal-tree` 保持 `active / 2/6`，未以修复 finding 伪造新的检查点完成。

### 对照关闭条件（F-004）

| 条件 | 判定 | 证据 |
|------|------|------|
| 修正与 finding 根因对应 | 满足 | `userWithRoles` 改用 `sameRoleSet` 集合语义 |
| 重复 legacy roles 迁移后可读 | 满足 | `TestMigrateExistingR2DuplicateRolesReadable` |
| Login / Refresh 受影响调用链恢复 | 满足 | `TestLoginAndRefreshAfterMigrateDuplicateRoles` |
| 真实集合差异仍 fail closed | 满足 | 既有 `TestReadDetectsRoleMismatch` 保留并通过 |
| 证据可重复核对 | 满足 | 本轮聚焦测试、全仓测试、vet、gofmt 均通过 |

### Findings

- 本 scope 内无开放 finding。
- F-004 的 `fixed` 关闭路径满足 P-003 的「可核对修正 + 产物/测试证据」要求；A-002 的 `conditional` 历史 verdict 不被改写，本复核只确认其 required finding 已合法闭合。

### 必改项汇总

- **required / 必改**：无（F-004 已由 `fixed` 合法闭合）。
- F-001 与 F-002 完整矩阵仍按 A-001 响应延期至 S6，超出本次复核 scope，未据此改变其状态。

### 与既有意见的异同

- A-003 与 A-002 同向：A-002 提出的 F-004 修复建议已由当前代码和测试证据兑现。
- A-003 不构成 S2 全量关门审计，也不对 S3～S6 作 pass 判断；S2 后续推进仍须由 `/govern` 汇总相关意见。

### 结论 + 建议给编排器/用户的下一步

- **结论**：F-004 的集合语义修正、迁移后双 lookup、Login/Refresh 回归和独立复跑证据均充分，`F-004 → fixed` 可被独立复核确认；本 scope verdict 为 `pass`。
- **建议 `/govern`**：记录 A-003 对 F-004 的独立复核通过，并基于当前无开放 required finding 的事实评估是否推进 S3；不得将本条 pass 解读为 S3～S6 或 Root R3 关门通过。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。响应、后续推进与关门由用户通过 **`/govern`** 处理。

---

## 响应 A-003（2026-08-02 · /govern）

- **模式**：response
- **响应意见**：A-003（independent · finding-closure `F-004` · verdict `pass`）
- **裁决**：采纳 A-003 `pass` 与证据认定——F-004 的 `fixed` 关闭证据（`sameRoleSet` 集合语义、迁移后双 lookup、Login/Refresh 回归、独立复跑）经 independent 复核确认，required 门禁解除。

### 关闭证据表

| Finding | 状态 | 处置 / 证据路径 |
|---------|------|-----------------|
| F-004 · required · 中 | **verified**（fixed → 独立复核确认） | A-003 `pass`：`store.go` `sameRoleSet` 集合语义 + `TestMigrateExistingR2DuplicateRolesReadable`（迁移后 `UserByID`/`UserByUsername`）+ `TestLoginAndRefreshAfterMigrateDuplicateRoles`（Login/Refresh 调用链）+ 聚焦/全仓测试、vet、gofmt 独立复跑通过。 |

### 仍开放项

- F-004 已合法闭合（`fixed` + independent 复核 `verified`）；S2 的 required 门禁无剩余。
- F-001 保持 open 跟踪至 **S6（或下一次迁移相关改动）**；F-002 完整 V-MIG-04 矩阵保持 open 至 **S6**（S2 部分已闭合）。均为 `recommended`，不阻断 S3。
- 无未合法闭合的 required 开放项。

### 冲突裁决

- A-003 与 A-002 同向（F-004 修复建议已兑现）；A-001 pass 仅覆盖 S1，与 A-002/A-003 无同范围相反 verdict。不构成冲突。

### 推进评估（S3）

- 当前无开放 required finding、无到期 required 信息项（`I-006-001/002` verified）、无意见冲突；S2 required 门禁已解除，具备推进 S3 的门禁条件。
- 按 P-004 §3.1：GOAL-006 目前仅有 independent 意见（A-001/A-002/A-003），无 `self` 意见；是否在 S3 前补一次自审由用户裁决（见下方问题）。
- A-003 `pass` 仅确认 F-004 关闭，**不**构成 S3～S6 或 Root R3 关门通过。

---

## A-004 · S1/S2 事实与 S3 门禁自审（2026-08-02）

- **source**：self
- **auditor**：Claude Code `/govern` · `04` 自审
- **类型**：execution-facts（主）+ finding-closure（辅）
- **mode**：stage（推进 S3 前阶段检查，P-004 §3.1 用户裁决先补自审）
- **scope**：`workspace-002-production-admin-foundation` / `GOAL-006-r3-persistent-rbac-menu` · S1 与 S2 检查点主张、A-001/A-002/A-003 意见与 F-004 关闭证据、S3 门禁就绪；不审 S3～S6 实现或 Root R3 关门
- **verdict**：**pass**
- **工作区**：`workspace-002-production-admin-foundation`；`root_goal=GOAL-001-production-admin-foundation`；`canonical_scope` 匹配；`shared_materials_catalog: none`

### 范围与区间

| 项 | 结论 |
|----|------|
| 意见台账 | A-001 pass（S1）、A-002 conditional（S2；F-004 required → fixed）、A-003 pass（F-004 独立复核）；全部已 responded |
| 开放必改 | 无（F-004 已 `fixed` + A-003 独立 `verified`；F-001/F-002 recommended → S6） |
| 信息门禁 | `I-006-001/002` verified；无到期开放 required |
| S3 门禁 | S1/S2 无未合法闭合 required；无 required 信息项；无意见冲突 |

### 成果（有证据）

1. **S1**：`schema_migrations` 台账 + 编译期 `0001/0002` 迁移链 + 顺序/缺号/未知版本/checksum 漂移 fail-closed + 单事务回滚 + pre-v0002 `VACUUM INTO` 快照与恢复验证（`migrate.go`、`migrate_test.go`，A-001 pass 复核）。
2. **S2**：阶段 B 终态——`userWithRoles` 双源集合核对 + 规范化权威读（`sameRoleSet`、按 key 升序）、`CreateUser`/`seedAdmin` 事务双写、派生 role 自建；`normalize_test.go` 5 用例覆盖双写/排序/分歧报错/seed 双写/FK-CASCADE（A-002 复核 + F-004 修正）。
3. **F-004 闭合**：`sameRoleSet` 集合语义修正 + `TestMigrateExistingR2DuplicateRolesReadable` + `TestLoginAndRefreshAfterMigrateDuplicateRoles`；A-003 independent 复核 `pass`。
4. **R2 契约保持**：auth Login/Refresh/JWT middleware、handler、account 全仓 API 回归通过；对外 `account.User {id,name,roles}` 形状不变。
5. **本轮独立复跑（2026-08-02）**：`go test ./... -count=1`（apps/api）全绿；`go vet ./...` 干净；`gofmt -l` 无输出。

### 对照成功标准

| 阶段 | 标准 | 判定 | 证据 |
|------|------|------|------|
| S1 | 台账 + 顺序/校验和 + 事务化 fail-closed + 可恢复副本 | 满足 | migrate_test.go 9 用例 + A-001 pass |
| S2 | 规范化关系 + 双写/双读核对 + 规范化权威读 + R2 契约 | 满足 | normalize_test.go 5 用例 + A-002/A-003 + F-004 fixed |

### Findings

- **F-101 · recommended · 低（自审新增，跟踪）**：F-001 与 F-002 完整 V-MIG-04 矩阵仍按既有响应延期至 S6；本轮自审确认不构成 S3 阻断，但建议在 S6 关门回归前闭合。
- 无 required / 必改项。

### 必改项汇总

- **required / 必改**：无。
- **recommended**：F-001、F-002（→ S6）、F-101（跟踪）。

### 结论 + 建议给编排器/用户的下一步

- **结论**：S1 与 S2 检查点主张有可重复证据支撑；F-004 已由 `fixed` + independent 复核合法闭合；无未闭合 required finding 或 required 信息项。S3 门禁就绪，可推进 **S3（增量幂等种子）**。
- **建议**：在用户已确认的 seed 接线方案（Open 内 seedAdmin=true 时运行 `seedRBAC`）下实施 S3；S3 完成后按需再审计。
- **限制**：本 `pass` 不构成 S3～S6 或 Root R3 关门通过。