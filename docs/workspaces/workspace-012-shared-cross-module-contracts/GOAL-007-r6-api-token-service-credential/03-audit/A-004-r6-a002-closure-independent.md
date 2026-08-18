---
id: A-004-r6-a002-closure-independent
goal: GOAL-007-r6-api-token-service-credential
doc: audit-entry
record_id: A-004
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-002 F-001～F-007 closure；复核 D-003、A-003 是否足以 fixed；0045 operation_log correlation-safe rebuild、created_by 无 FK、NOCASE unique/稳定错误、prefix-before-devSession、管理审计同 transaction、user-only 清单/data-permission self、分页/ID
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
responds_to: A-002
reviews: A-003
---

# A-004 · A-002 F-001～F-007 关闭复核（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：ad-hoc / finding-closure
- **scope**：GOAL-007 R6 S0；仅复核 A-002 F-001～F-007 是否已由 D-003 / A-003 给出可重复核对的 `fixed` 证据。重点：0045 `operation_log` correlation-safe rebuild、`created_by` 无 FK、NOCASE unique/稳定错误、prefix-before-devSession、管理审计同 transaction、user-only 清单/data-permission `self`、分页/ID。
- **verdict**：**pass**
- **required findings**：0

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 覆盖本目标；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：A-002 原文 F-001～F-007；A-003 候选响应；D-002（`superseded`）与 D-003（`proposed`）；E-003；对照现行代码：`operationlog` 0043 rebuild、`DeleteUser`、NOCASE/`INVALID_CREATE_FIELD`、`auth.Middleware`/`injectDevSession`、`IdentityFrom` 全调用点、`resolveScope`、`Store.WithTx`、既有 `page`/`pageSize`。
- **excluded**：不改 `status` / `progress` / `00-meta` / D-003 / goal-tree / 业务代码；不读取或比较其他工作区上下文；不审 S1 实现或 S3 关门证据；不把 A-001 recommended 实施门升为本条 required；不运行 `go test` / e2e（S1 未开工，符合预期）。
- **共享资料**：无固定引用；不得当作事实或 finding 关闭依据。
- **P-005**：本意见不改信息表。I-002～I-004 登记仍为 `collecting`；本轮结论是设计证据已足够，编排器可改为 `verified`。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-007 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R6 / VP-012 / Charter 的明显冲突 | D-003 只修订 A-002 指出的契约缺口，不改 Profile/module/Manifest 不变式 |
| Vision Review required | 本 scope 未见开放 required | 本意见不写 `docs/vision/reviews.md` |
| 既有 Goal 审计 | A-001 self = pass（0 required）；A-002 independent = conditional（F-001 high / F-002–F-003 med required；F-004–F-007 recommended）；A-003 self = pass（proposed fixed，不自闭） | `03-audit.md` |
| P-004 冲突 | 无待裁冲突 | A-003 对 A-002 全部意见走 `fixed` 路径，无 residual/overrule，也无一要一否 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 现行最大 migration 仍为 0043，0044/0045 作为下一对全局号仍空 | `apps/api/internal/modules/operationlog/migration/migration.go` L361–364 Version 43；jobs 42；collector 登记无 44/45 | 通过。与 A-002「0044 为下一号」一致；D-003 把 CHECK 扩到 0045 避免 auth-session 越权改 operationlog |
| 0043 已是 correlation-safe rebuild 先例 | 同文件 `migrateOperationLogWalletJobs` L423–447：TEMP `operation_log_correlation_backup` → DROP correlation → `rebuildOperationLog` → 重建 correlation DDL/index → INSERT 恢复 → DROP backup。`operationLogCorrelationDDL` L483–488 含 `REFERENCES operation_log(id) ON DELETE CASCADE`，故必须先拆 sidecar | 通过。D-003 §1.3 四步与 0043 同构 |
| 0043 CHECK 不含三件 service-credentials 事件 | `operationLogWalletJobsDDL` L495；`repository_test.go` `TestRepositoryRejectsUnknownEventAndExposesFailureSeam` | 通过。F-001 原文缺口仍在现码中，须由尚未开工的 0045 关闭；设计已冻结归属与测试 |
| `DeleteUser` / batch 仍只清 `refresh_tokens` 再删 `users` | `authsession/users_repository.go` L264–268、L282–285 | 通过。D-003 去掉 FK 后该路径不再被 RESTRICT 打成未冻存储错误 |
| 名称不敏感唯一在现网不是列 UNIQUE 默认行为 | `users`/`roles` DDL 为大小写敏感 `UNIQUE`；`COLLATE NOCASE` 仅用于查询（`users_repository.go` L436–439，`roles_repository.go` L435–438） | 通过。D-003 把 NOCASE 放到列 UNIQUE，比现网 username 更强，且正是 A-002 要求 |
| `INVALID_CREATE_FIELD` + `fieldErrors` 已存在、不扩 catalog | `errorcatalog.go` L58、L215–254；`writeLocalizedFieldError`（`localize.go` L35–39）；`resources.go` L575 | 通过。D-003 映射与既有 field-validation 包络一致 |
| `devSession` 会在 JWT 解析失败后抬升任意 Bearer | `auth.go` `Middleware` L459–463：`ParseAccessToken` 失败且 `devSession` → `injectDevSession`（全权限 StaticDevSession） | 通过。F-004 原文风险仍在现码；D-003 §4 把 `sui_sc_` 判据移到该回落之前 |
| `RecordOperation` 自开事务，且 `Store.WithTx` 不可重入 | `operationlog/repository.go` L133–166 内部 `withTx`；`store.go` L68–80 每次 `BeginTx`，无 savepoint/嵌套 | 通过。D-003 §5「既有 tx 上的 recorder + authsession callback」是唯一可落地的同事务写法 |
| 现行 IdentityFrom 自助入口与 D-003 六项清单一致 | 调用点：`account.go` `/api/accounts/me`；`account_self.go` identity helper（profile/password/sessions）；`account_avatar.go` POST upload；`mfa.go` `requireIdentity`；`notifications.go` `identity`；`wallet_self.go` `selfIdentity`。其余 `IdentityFrom` 仅在 `resources`/`export`/`import`/`upload` | 通过。permission-gated 路径按 D-003 保留通用 identity |
| data-permission `self` 以 `user.ID` 为 actor | `resources.go` `resolveScope` L303–307：`ScopeFor(user.ID, …)`；`scopeOwned` L312–316 比较 `constraint.ActorID` | 通过。机器 `ID=service-credential:<id>` 不会继承 `created_by` 行范围 |
| 既有 list 分页约定 | `datapermission.go` L70–95、`account_self.go` L298–364：默认 pageSize 20、最大 100、`{items,total,page,pageSize}` | 通过。D-003 §7 复用该形状 |
| A-003 未冒充 independent closure | A-003 L23、L37–39：proposed fixed；权威开放投影仍为 3，待本条 | 通过 |

## 对照成功标准（本轮 scope = 设计可冻结性）

| 标准 / 关闭条件 | 本轮 | 证据 |
|-----------------|------|------|
| F-001：0044/0045 归属、三事件字面量、correlation-safe rebuild、测试 | **fixed** | D-003 §1 |
| F-002：`created_by` 与 `DeleteUser`/batch 交互及错误码 | **fixed** | D-003 §2（无 FK；删除行为不变） |
| F-003：NOCASE UNIQUE / trim / 稳定重名码 | **fixed** | D-003 §3、§8 |
| F-004：`sui_sc_` 先于 JWT/`devSession`，失败只 401 | **fixed** | D-003 §4 |
| F-005：create/revoke 审计与行变更同事务；use/`last_used_at` best-effort | **fixed** | D-003 §5 |
| F-006：完整 user-only 清单 + machine `self` 作用对象 | **fixed** | D-003 §6 |
| F-007：吊销行可见、分页、32-char hex id | **fixed** | D-003 §7、§8 |
| 成功标准 3（可审计且不泄密）在 schema 层可落地 | **设计可达** | 0045 扩 CHECK + 管理审计 fail-closed；secret/hash/header 禁入仍继承 D-002 |
| I-002～I-004 最晚阶段 = S0 结束前 | 设计已闭合；登记仍 `collecting` | 见下表 |

## A-002 disposition

### F-001 · 0044 未冻结 `operation_log.event` CHECK 扩展 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | high / required |
| 状态 | **fixed** |
| closure | D-003 §1；A-003 L29 |
| evidence | D-003：0044 仅 `core.auth-session`/`service_credentials`；0045 归 `core.operationlog`，在 0043 CHECK 全集上追加 `service-credentials.create/use/revoke`；rebuild 步骤与 `migrateOperationLogWalletJobs` 同构；tests 固定顺序、历史事件/correlation 保留、新事件可写、未知事件仍拒 |

A-002 要求的三件事——归属、编号、correlation 重建——均已成文。选择 0045 而不是把 CHECK 塞进 0044，符合「事件扩展历来由 `core.operationlog` rebuild」的现网事实。本条确认的是设计可落地，不是 0045 已实现。

### F-002 · `created_by ON DELETE RESTRICT` 与用户删除未对齐 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / required |
| 状态 | **fixed** |
| closure | D-003 §2；A-003 L30 |
| evidence | D-003 删除 users FK；`created_by` 与 `operation_log.actor_id` 同为持久历史标识；`DeleteUser`/batch 不级联、不阻断、不吊销；metadata 只回 `createdBy` id、读取不强制 join |

A-002 列出的三种写法（拒删 / 先吊销 / 仅未吊销 RESTRICT）是示例，不是穷尽合法集。无 FK + 独立生命周期与 D-001/D-002「凭据不随创建者删除或禁用而失效」一致，且直接消除 `DeleteUser` 被 SQLite RESTRICT 打成未冻存储错误的路径。本条接受该第四种明确语义为 `fixed`。

### F-003 · 名称大小写不敏感唯一缺少 DDL 与冲突码 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / required |
| 状态 | **fixed** |
| closure | D-003 §3、§8；A-003 L31 |
| evidence | `name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100)`；handler 存 trim；`ErrCredentialNameTaken` → 既有 `400 INVALID_CREATE_FIELD` + `fieldErrors=[{field:"name", reason:"name already exists"}]`；并发 `Build`/`build` 失败方不得为 `INTERNAL` |

A-002 点名的 DDL、规范化与「不扩 catalog 的既有码」均已冻结。列级 NOCASE UNIQUE 可区分 `service_credentials.name` 与 `token_hash` 的 UNIQUE 失败信息，映射可实现。

### F-004 · `sui_sc_` 分流必须先于 JWT / `devSession` — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / recommended |
| 状态 | **fixed** |
| closure | D-003 §4；A-003 L32 |
| evidence | 前缀判断在 `ParseAccessToken` 与任何 `devSession` fallback **之前**；匹配后 unknown/malformed/expired/revoked 一律 `401 UNAUTHENTICATED`；非前缀保持现有 JWT + 显式 dev-session |

这正中 `auth.go` L459–463 的现网回落。S2 必须按此改 Middleware；本条关闭的是设计缺口，不是实现。

### F-005 · create/revoke 审计不应与 use/`last_used_at` 同为 best-effort — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / recommended |
| 状态 | **fixed** |
| closure | D-003 §5；A-003 L33 |
| evidence | create/revoke 的行变更与 operation log（含既有 recorder 会写的 correlation sidecar）必须在**同一个** `Store` transaction；失败回滚并 `500 INTERNAL`；重复 revoke 已有 `revoked_at` 则 204 且不再追加审计；`use`/`last_used_at` 仍 best-effort |

现网 `Recorder.RecordOperation` 自开事务，而 `Store.WithTx` 不可重入；D-003 要求 authsession 接收不依赖 operationlog 包的 tx callback、operationlog 提供「使用既有 tx 的 recorder」，包边界与原子性可同时成立。不得在外层 tx 内再调现有 `RecordOperation`。

### F-006 · 自助隔离清单与 data-permission `self` — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / recommended |
| 状态 | **fixed** |
| closure | D-003 §6；A-003 L34 |
| evidence | 六类现网 `IdentityFrom` 自助消费者改 `UserIdentityFrom` 并逐项测机器凭据 401；permission-gated 路径保留通用 identity；`self` 作用于 `service-credential:<id>`，不模拟 `created_by` |

对照全仓 `auth.IdentityFrom` 调用点，D-003 清单完整：`account.go`、`account_self.go`（profile/password/sessions）、`account_avatar.go` 的 authenticated mutation、`mfa.go` self-service（`requireIdentity`，不含 admin-reset）、`notifications.go`、`wallet_self.go`。`resources`/`export`/`import`/`upload` 正确排除。

### F-007 · list/detail 已吊销行、分页与主键形状 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | low / recommended |
| 状态 | **fixed** |
| closure | D-003 §7、§8；A-003 L35 |
| evidence | id = 16 random bytes / 32-char lowercase hex，`CHECK (length(id)=32)`，非法或不存在均 `404`；list 默认含 active/expired/revoked 并回 `revokedAt`/`lastUsedAt`；`page`/`pageSize` 默认 1/20、最大 100，`createdAt DESC, id DESC`，`{items,total,page,pageSize}` |

与既有管理 list 形状一致。DDL 以长度约束主键、charset 由生成与路径校验保证，对 S0 足够。

## Findings

本轮 **无新 required / recommended finding**。A-001 F-001～F-003 recommended 仍是 S1/S2/S3 实施门，不升格，也不被本条重开。

## 信息项核对（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持 |
| I-002 | required | S0 结束前 | collecting | **设计可转 `verified`**：principal 形状、prefix-before-devSession、完整 user-only 清单、machine `self` 已写入 D-002§3 + D-003§4/§6 |
| I-003 | required | S0 结束前 | collecting | **设计可转 `verified`**：0044 DDL、无 FK `created_by`、NOCASE 唯一、hash-only/过期/吊销、id/分页已写入 D-003§1–3/§7–8 |
| I-004 | required | S0 结束前 | collecting | **设计可转 `verified`**：0045 三事件 + correlation rebuild、管理审计同事务、稳定重名码、R5 覆盖仍继承 D-002 |
| I-005 | required | S3 关门 | verified | 维持（bytes/pin 仍须 S3 复证） |
| I-006 | required | S0 结束前 | verified | 维持；本条即 A-002 的 independent finding-closure |

无 `deferred`。无用户书面 `accepted-residual`。`00-meta` 证据列仍写「待 D-002/A-*」，属编排器更新债务，不是开放 required finding。

## 必改项汇总

**无。** required = 0。

A-002 F-001～F-003 required 与 F-004～F-007 recommended 均可按 P-003 `fixed` 合法闭合。

## 与既有意见的异同

| 点 | A-002 | A-003 | 本意见 |
|----|-------|-------|--------|
| F-001 | required / open | proposed fixed | **fixed**（0045 + 0043 同构 rebuild 可重复核对） |
| F-002 | required / open | proposed fixed | **fixed**（无 FK 为合法明确语义，不是漏写） |
| F-003 | required / open | proposed fixed | **fixed** |
| F-004～F-007 | recommended / open | proposed fixed | **fixed** |
| I-002～I-004 | 维持 collecting | 仍 collecting，待本条 | 设计充分；状态改写归 `/govern` |
| D-003 状态 | 要求修订后复审 | `proposed`，不自闭 | 同意保持至编排器响应；本条不改 decision |
| A-001 pass | 不同意「设计已完整」 | 用 D-003 补洞 | **同意 A-002 原判**；补洞后可放行 S0 |
| verdict | conditional | pass（响应，不自闭） | **pass**（closure） |

不是 P-004.2 冲突：A-003 未宣称已闭合，本条确认其候选成立。

## 结论 + 建议给编排器/用户的下一步

D-003 作为 A-002 的修订契约**足以**关闭 F-001～F-007：0045 按 0043 先例做 correlation-safe CHECK 扩展；`created_by` 无 FK 使 `DeleteUser` 行为保持不变；NOCASE UNIQUE 映射到既有 `INVALID_CREATE_FIELD`；`sui_sc_` 先于 JWT/`devSession`；管理审计与行变更同事务；user-only 清单与现网 `IdentityFrom` 自助入口对齐，且 `self` 不继承创建者。S1 代码尚未开始，本条不把设计闭合读成已交付。

建议 `/govern` 下一句：响应 A-004，将 A-002 F-001～F-007 标为 `fixed`，把 I-002～I-004 改为 `verified`（证据改为 D-003 + A-004），将 D-003 标为 `accepted`，关闭 S0 并放行 S1。不要改写 D-003 正文。S1 须按 D-003 实现 0044/0045（若届时 44/45 已被其他 migration 占用则另开决策重编号，不得 silently 复用），并证明管理审计失败会回滚凭据行。

## 声明

本意见不修改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 业务代码。响应与信息项闭合由 `/govern` 处理。
