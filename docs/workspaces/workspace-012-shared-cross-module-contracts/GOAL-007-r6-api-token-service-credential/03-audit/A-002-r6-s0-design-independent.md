---
id: A-002-r6-s0-design-independent
goal: GOAL-007-r6-api-token-service-credential
doc: audit-entry
record_id: A-002
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R6 S0 design-plan；D-002、I-002～I-004、机器 principal 与用户自助隔离、secret/hash/expiry/revoke、scope ceiling、管理 API、审计敏感数据、0044 migration、Profile/Manifest/protocol 不变式
audit_type: design-plan
verdict: conditional
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews: A-001
---

# A-002 · R6 S0 independent 设计审计（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：design-plan / S0
- **scope**：GOAL-007 R6 S0 方案冻结前交叉复核。重点：D-002、I-002～I-004、机器 principal 与用户自助隔离、secret/hash/expiry/revoke、scope ceiling、管理 API、审计敏感数据、0044 migration、Profile/Manifest/protocol 不变式。
- **verdict**：**conditional**
- **required findings**：3（F-001 high；F-002 / F-003 med）

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 覆盖本目标；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：GOAL-007 `00-meta`、D-001、D-002、E-001、E-002、A-001；只读 Root 路线图指针与 VP-012 / Charter `schema-ui-core-admin-foundation@0.2.0`；对照现行代码：`auth.Middleware` / `NewOpaqueToken` / `HashToken`、`account.User`、`requirePermission` / 自助面 IdentityFrom、`operation_log` 0043 CHECK、authsession migration 编号、`DeleteUser`、R5 `WithOperationalGate`、`manifest.rejectFragmentSecrets`、`kernel.BuiltinModules`、composition system-data 注入。
- **excluded**：不改 `status` / `progress` / `00-meta` / D-002 正文 / goal-tree；不读取或比较其他工作区上下文；不审 S1 实现或 S3 关门证据；不把 A-001 的 recommended 实施门禁升为本条 required。
- **共享资料**：无固定引用；不得当作事实或 finding 关闭依据。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-007 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R6 / VP-012 / Charter 的明显冲突 | VP-012 将 R6 限定为机器凭据管理面/作用域/吊销/审计且与用户会话分离；GOAL-007 非目标与 VP-012「不改装配语义、不承载 Tier D」一致 |
| Vision Review required | 本 scope 未见开放 required | 本意见不写 `docs/vision/reviews.md` |
| P-005 信息项 | I-001 / I-005 / I-006 可维持 verified；I-002～I-004 仍为 `collecting`，最晚阶段 = S0 结束前 | `00-meta.md` 信息表；本条 F-001～F-003 阻断其关闭 |
| P-004 冲突 | 无互否必改项 | A-001 为 pass / required=0；本条新增 required，不是对同一 finding 一要一否。编排器须响应本条后才能关 S0 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 用户 JWT 与 refresh 边界扫描属实，可作机器凭据对照 | `auth.go`：HS256 JWT 仅 `Subject`+`token_version`；`NewOpaqueToken` 32-byte `RawURLEncoding` + `HashToken` SHA-256 hex；禁用/改密走 `token_version` | 通过。I-001 维持 verified |
| `sui_sc_` 前缀可与 JWT 无歧义分流 | JWT compact 头为 `eyJ…`；D-002 §1 要求大小写敏感前缀先分流，失败统一 `401 UNAUTHENTICATED` | 通过（设计层）。S2 必须前缀优先，见 F-004 |
| hash-only、一次性 secret、独立生命周期与现有 refresh 模式一致且不复用 session 行 | D-002 §1–2；`refresh_tokens` 已是 hash-only / 过期 / 原子撤销 | 通过。完整 raw（含前缀）入 hash 可避免截断碰撞 |
| 0044 作为下一全局 version 编号正确 | 现行最大为 operationlog `43`、jobs `42` | 通过编号；DDL/事件面不完整，见 F-001 |
| principal 投影可复用 `requirePermission` | `account.User.Permissions` 是唯一授权解释器；`requirePermission` 只读该切片，不读 Roles | 通过。`Roles=[]` 不破坏 permission gate |
| 用户自助面可被 `UserIdentityFrom` 收口 | Identity-only 入口现为 `/api/accounts/me`、`account_self`、`account_avatar`、`mfa.requireIdentity`、`notifications.identity`、`wallet_self.selfIdentity` | 通过方向。清单需在 S2 锁死，见 F-006 |
| scope ceiling 三重约束可 fail closed | catalog 存在性 + 创建者 `PermissionsForUser` 子集 + 保留 `service-credentials.read/write`；未知 key 复用已有 `INVALID_PERMISSION_REF` | 通过。`errorcatalog` 已有该码，无需扩 catalog |
| 管理 API 可挂在 composition 而不改 Profile/Manifest | 既有 `files.write` 在 `RegisterContributions` 之后追加 system-data；upload 走中央 mux 而非 `BuiltinModules()` 声明 | 通过。按此路径可不改 `kernel.BuiltinModules()` / fragment / nav / page |
| R5 会自动罩住 POST create/revoke | `operational.go`：已注册 POST 为 business mutation；allowlist 仅 login/refresh/logout/mfa-verify/改密 | 通过。GET metadata 在受控态保持可读，与 D-002 一致 |
| 敏感字段与 Manifest 不变式可核对 | `manifest.secretKeyNames` 含 `token`/`secret`/`authorization`/`apikey`；D-002 禁止 secret/hash/header 进 metadata、operation detail、service log、public Manifest | 通过设计意图。I-005 维持 verified（S3 仍须复证 bytes/pin） |
| I-006 审计模式 | meta 已冻结 `cross` + grok-build independent；A-001 self 已落盘 | 通过。本条补上独立审 |

## 对照成功标准（S0 设计可冻结性）

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. 机器凭据与用户 JWT/refresh 分离；secret 只创建返回一次；持久层只存不可逆 hash | **设计可达** | D-002 §1–2 与现行 opaque/hash 原语一致；不签 JWT、不碰 `token_version` |
| 2. 有效/过期/吊销可判定；scope 只从既有 permission key 解析；越权 fail closed | **部分** | 状态机与 ceiling 清楚；名称唯一/用户删除副作用未冻，见 F-002/F-003 |
| 3. 管理 mutation 受权限与 R5 保护；创建/使用/吊销审计含 actor/credential/correlation 且不泄密 | **不可按原文落地** | 事件名已定，但 0043 `operation_log.event` CHECK 不含新事件，见 F-001 |
| 4. Profile 默认集、模块矩阵、Manifest 装配与协议 pin 可保持不变 | **设计可达** | 无新 module/page/nav/fragment；permission 走 system-data；secret 不得进 Manifest |

## Findings

### F-001 · 0044 未冻结 `operation_log.event` CHECK 扩展

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| 关联 | I-003、I-004；成功标准 3 |
| evidence | D-002 §2 只定义 `0044 service_credentials`（`core.auth-session`）；D-002 §5 新增 `service-credentials.create/use/revoke`。现行 0043 DDL：`apps/api/internal/modules/operationlog/migration/migration.go` L490–495，`CHECK (event IN (…))` 无上述三事件。`operationlog/repository_test.go` `TestRepositoryRejectsUnknownEventAndExposesFailureSeam` 证明未知 event 被拒绝。事件扩展历来由 `core.operationlog` rebuild（0043 还要备份/恢复 `operation_log_correlation`）。 |
| closure | — |

D-002 同时写「operation log 失败…均为 best-effort，不改变已经通过的认证结果」。若 S1 只建 `service_credentials`、不扩 CHECK，create/revoke/use 插入会失败并被吞掉，成功标准 3 的审计主张名存实亡。S0 必须补：下一 version 如何扩 CHECK（0044 兼做或 0045 归属 `core.operationlog`）、三事件字面量、以及 correlation 表在 rebuild 时的备份/恢复。在此之前不得把 I-003/I-004 标为 verified，也不得放行 S1。

### F-002 · `created_by ON DELETE RESTRICT` 与现行用户删除未对齐

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| 关联 | I-003 |
| evidence | D-002 §2：`created_by` FK `users` + `ON DELETE RESTRICT`。`authsession/users_repository.go` `DeleteUser` / `DeleteUsersBatch` 只删 `refresh_tokens` 再 `DELETE FROM users`，无 credential 预检或映射。SQLite RESTRICT 会在仍有凭据（含已吊销）时让用户删除失败，现路径会变成未冻结的存储错误。 |
| closure | — |

独立生命周期已排除「删用户级联删凭据」和「凭据随创建者禁用而失效」，这与 D-001/D-002 一致，但 S0 仍须写明删除语义：拒绝删除（稳定错误码）、先要求吊销、或仅未吊销行 RESTRICT。未写则 S1 会破坏已有 users.delete。

### F-003 · 名称大小写不敏感唯一缺少 DDL 与冲突码

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| 关联 | I-003、I-004 |
| evidence | D-002 §2/§4：`name` 非空、大小写不敏感唯一，1～100 chars；复用错误码列表无 taken/conflict。本模块对 username/role 的不敏感比较已用 `COLLATE NOCASE`（`users_repository.go` L436–439），但 UNIQUE 默认大小写敏感。D-002 未写规范化存储、`COLLATE NOCASE` 唯一索引，也未把重名映射到已有 `INVALID_CREATE_FIELD`（或明确不扩 catalog 的其他码）。并发双写会落到 `INTERNAL`。 |
| closure | — |

### F-004 · `sui_sc_` 分流必须先于 JWT / `devSession` 回落

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| 关联 | I-002；S2 |
| evidence | `auth.go` `Middleware`：JWT 解析失败时，若 `devSession` 则 `injectDevSession`（全权限 StaticDevSession），否则 401。D-002 §1.3 写了前缀分支与统一 401，但未点名「前缀匹配在 `ParseAccessToken` 之前；credential 失败不得落入 JWT 或 `injectDevSession`」。 |
| closure | — |

S2 若把 credential 查找放在 JWT 失败之后，dev/测试会把任意 `sui_sc_*` 提升为 Dev Admin。应写入 D-002 或 S2 验收：前缀优先，credential 失败只回 401。

### F-005 · create/revoke 审计不应与 use/`last_used_at` 同为 best-effort

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| 关联 | I-004 |
| evidence | D-002 §5 末句把 operation log 失败与 `last_used_at` 失败一并 best-effort。`Recorder.RecordOperation` 注释写明由调用方决定策略。 |
| closure | — |

`use` 与 `last_used_at` 不阻断已通过认证是合理的。create/revoke 是管理 mutation，审计失败应 fail closed（回滚或 500），否则补上 CHECK 后仍可能丢掉唯一创建记录。建议 D-002 拆开两条策略。

### F-006 · 自助隔离清单与 data-permission `self` 作用对象

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| 关联 | I-002；与 A-001 F-002 同向 |
| evidence | D-002 §3 按面描述 `UserIdentityFrom`。现行 identity-only 入口见上表。`resources.go` `resolveScope` 用 `user.ID` 取 data-permission `self`；机器 ID 为 `service-credential:<id>`，不会继承 `created_by` 的行范围。 |
| closure | — |

S2 应枚举全部 IdentityFrom 自助入口并加 helper 测试；并写明 `self` scope 作用于机器 principal，不模拟创建者。这与「不继承创建者实时角色」一致，但现在只隐含。

### F-007 · list/detail 对已吊销行、分页与主键形状未写

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 关联 | I-003、I-004 |
| evidence | D-002：吊销幂等 204 / 不存在 404；未说 GET list/detail 是否仍返回已吊销行、是否分页、`id` 字符集。 |
| closure | — |

建议：list/detail 默认含 `revokedAt`；list 复用既有 page/pageSize；`id` 冻结为稳定 hex/UUID，避免路径歧义。不阻断 S0。

## 必改项汇总

1. **F-001（high / required）**：在 D-002 冻结 `service-credentials.create/use/revoke` 进入 `operation_log.event` CHECK 的 migration 归属、编号与 correlation 重建步骤；未完成前不关闭 I-003/I-004、不放行 S1。
2. **F-002（med / required）**：冻结 `created_by` 与 `DeleteUser` / batch-delete 的交互及错误码。
3. **F-003（med / required）**：冻结名称大小写不敏感唯一的 DDL/规范化，以及重名的既有错误码映射。

## 与既有意见的异同

| 来源 | 结论 | 本条 |
|------|------|------|
| A-001 self · pass · required=0 | 认为 D-002 已可交叉复核后关闭 I-002～I-004 | **不同意其「设计已完整」**。分流、ceiling、hash-only、human-only 管理、不变式方向成立，但漏了现行 operation log CHECK 与用户删除/重名门禁 |
| A-001 F-001～F-003 recommended | S1 直接查库证 hash-only / 并发吊销；S2 覆盖全部 user-only consumer；S3 锁 Profile/Manifest bytes 与 R5 mutation | **同意**，保持 recommended，不升格。本条 F-006 补 data-permission `self` |
| I-001 / I-005 / I-006 | verified | **同意**（I-005 的 bytes/pin 仍须 S3 复证） |
| I-002～I-004 | collecting，待 D-002 + cross | **维持 collecting**。I-002 主体可冻，但仍受 F-004/F-006 建议修订；I-003/I-004 被 F-001～F-003 阻断 |

无 P-004 必改互否。编排器应修订 D-002 并再走 self/independent 复审，而不是用 A-001 pass 覆盖本条。

## 结论 + 建议给编排器/用户的下一步

D-002 作为 R6 精确契约**方向正确**：前缀分流、hash-only、独立 principal、human-only 管理、scope 三重上限、R5 自动覆盖、不新增 Profile/module/Manifest fragment。独立复核**不能**无条件放行 S0——F-001 使「可审计」在当前 schema 下无法按原文落地。

建议 `/govern`：响应 A-002 F-001～F-003（修订 D-002，保持 `proposed` 直到复审）；不要把 I-002～I-004 标 verified，不要开 S1；F-004～F-007 可一并写入 D-002 以免 S2 返工。修订后对本 scope 再跑 `/audit` 复审关闭证据。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
