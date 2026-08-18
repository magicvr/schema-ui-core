---
id: A-007-r6-s3-closeout-independent
goal: GOAL-007-r6-api-token-service-credential
doc: audit-entry
record_id: A-007
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R6 S3 close-out；D-003、提交 aa00f33/ce8d952/1864f49/2a2d0dd、E-005、A-006；secret/hash 隔离、0044/0045 correlation-safe migration、事务 create/revoke audit、service prefix 在 JWT/dev fallback 前、scope ceiling 与 human-only 管理、全部 user-only consumer 隔离、R5 operational gate、Profile/Manifest/protocol 不变式及全量测试证据
audit_type: close-out
verdict: conditional
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews: A-006
---

# A-007 · R6 S3 independent close-out（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：close-out / S3
- **scope**：GOAL-007 最终 S3 cross independent 关门审计。核对 D-003（及 D-002 未被覆盖条款）、四条成功标准、I-001～I-006、A-001～A-006、提交 `aa00f33` / `ce8d952` / `1864f49` / `2a2d0dd`、E-005、A-006。重点：secret/hash 隔离、0044/0045 correlation-safe migration、事务 create/revoke audit、service prefix 在 JWT/dev fallback 前、scope ceiling 与 human-only 管理、全部 user-only consumer 隔离、R5 operational gate、Profile/Manifest/protocol 不变式及测试证据。
- **verdict**：**conditional**
- **required findings**：1（F-001 med）

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 覆盖本目标；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：GOAL-007 `00-meta` / D-001 / D-002（`superseded`，未被 D-003 覆盖的条款仍有效） / D-003（现行 `accepted`） / E-001～E-005 / A-001～A-006；Root R6 指针；VP-012；实现与测试：`service_credentials.go`（authsession + handler）、`auth.go` Middleware、`0044`/`0045` migration、`RecordOperationTx`、composition 权限注入与 use recorder、六处 `UserIdentityFrom`、R5 `WithOperationalGate`、相关 tests。
- **excluded**：不改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 业务代码 / 决策或执行正文；不读取或比较其他工作区上下文；不把 A-001 recommended 实施门升为新的 required（除本轮独立核出的契约偏差）。
- **共享资料**：无固定引用；不得当作事实或 finding 关闭依据。
- **本轮复验**：`apps/api` 定向 `go test -timeout 15m -count=1`：`authsession` 38.206s、`auth` 25.883s、`handler` 246.716s、`store` 43.710s、`operationlog` 7.605s、`composition` 25.590s、`docscheck` 0.503s，全部 ok。并发 `TestServiceCredentialConcurrent` `-count=5` ok。`kernel` / `manifest` ok。未重跑 `npm run build`：R6 四提交未改 protocol/kernel/manifest 源文件，工作树干净；I-005 以提交范围 + composition/kernel/manifest 测试复核。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-007 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R6 / VP-012 / Charter 的明显冲突 | VP-012 将 R6 限定为机器凭据管理面/作用域/吊销/审计且与用户会话分离；实施未新增 Profile/module/page/nav/fragment |
| Vision Review required | 本 scope 未见开放 required | 本意见不写 `docs/vision/reviews.md` |
| 既有 Goal 审计 | A-001～A-005 无开放 required；A-002 F-001～F-007 已由 A-004 `fixed`；A-006 self = pass，声明不代替 independent | `03-audit.md` |
| P-004 冲突 | 无互否必改项 | A-006 为 pass / required=0；本条新增 F-001 required，不是对同一 finding 一要一否 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 0044 仅建 `service_credentials`；DDL 与 D-003 §8 一致；`created_by` 无 users FK | `authsession/migration/migration.go` L135–150、L198–204、L273–279 | 通过。`name COLLATE NOCASE UNIQUE`、`id`/`token_prefix`/`token_hash` 长度 CHECK、两枚索引均在 |
| 0045 归 `core.operationlog`，在 0043 CHECK 上追加三事件；rebuild 与 0043 同构 | `operationlog/migration/migration.go` L366–372、L456–479、L538–550 | 通过。TEMP backup → DROP correlation → rebuild → 重建 correlation → restore → DROP backup |
| 0044/0045 顺序、历史 correlation 保留、新三事件可写、未知事件仍拒 | `store/migrate_test.go` L124–130、L619–621；`store/migrate_0045_test.go` L11–54；`operationlog/repository_test.go` L112–119 | 通过。fresh catalog 以 45/`operation_log_service_credentials` 结尾；0045 前写入的 `wallet.reconcile` correlation 保留；`records.purge` 仍 CHECK 失败 |
| hash-only：完整 raw（含前缀）入 SHA-256 hex；库中无 raw | `auth.go` L388–396、L402–405；`service_credentials_test.go` L17–54 | 通过。`tokenPrefix = raw[:15]`；`SELECT ... token_hash = raw OR token_prefix = raw` 计数为 0 |
| create/revoke 与 audit 同 Store transaction；失败回滚 | `service_credentials.go` L47–66、L140–169；`operationlog/repository.go` L149–181；`service_credentials_test.go` L81–96；`handler/service_credentials_test.go` L133–166 | 通过。`RecordOperationTx` 使用调用方 tx，不自开事务；forced failure 后行数为 0 / `revoked_at` 仍空；重复 revoke 不二次 audit |
| 并发 NOCASE 重名与单次 revoke 转换 | `service_credentials_test.go` L99–186；本轮 `-count=5` | 通过。`Build`/`build` 一方 `ErrCredentialNameTaken`；并发 revoke `changes=1` 且 `audits=1` |
| `sui_sc_` 在 `ParseAccessToken` 与 `devSession` fallback **之前**；失败只 401 | `auth.go` L493–516、L568–574；`auth_test.go` L63–130 | 通过。`devSession=true` 时 unknown/revoked 仍 401、`called=false`，不抬升 Dev Admin |
| principal 投影：synthetic id、空 roles、冻结 scopes | `auth.go` L575–581；`account/session.go` L24–35、L461–463 | 通过。`UserIdentityFrom` 拒绝 service principal |
| human-only 管理 + scope ceiling（catalog / 创建者子集 / 保留键） | `handler/service_credentials.go` L70–85、L152–169；`service_credentials_test.go` L58–60、L91–130 | 通过。机器凭据 GET 管理面 403；unknown → `INVALID_PERMISSION_REF`；reserved / excess → `FORBIDDEN` |
| 六处 user-only 已改 `UserIdentityFrom`；permission-gated 保留 `IdentityFrom` | `account.go` L22；`account_self.go` L83；`account_avatar.go` L51；`mfa.go` L231；`notifications.go` L65；`wallet_self.go` L107；`resources.go` L289；`export.go` / `import.go` / `upload.go` | 通过实现。HTTP 黑盒仅测了 `/api/accounts/me`，见 F-004 |
| data-permission `self` 作用于机器 id，不继承 `created_by` | `resources.go` L303–316；identity `ID=service-credential:<id>` | 通过代码。无专用资源行测试 |
| R5 罩住 credential POST；GET 仍为非 mutation | `handler/operational.go` L17、L50–57；`r5_operational_gate_test.go` L75–81 | 通过。三种受控态 `POST /api/service-credentials` → 503 + `SERVICE_*` + correlation |
| Profile / 模块矩阵 / Manifest / 协议 pin 无交付变化 | `git diff --name-only aa00f33^..2a2d0dd` 不含 `kernel/`、`manifest/`、`apps/web/public/protocol`、`docs/schemas`；`composition.go` L435–461 仅 system-data + 中央挂载；`composition_test.go` L456–482 权限 +2、navigation 不变 | 通过。未新增 module/page/nav/fragment |
| `DeleteUser` / batch 不级联、不吊销凭据 | `users_repository.go` L264–268 | 通过。仍只清 `refresh_tokens` 再删 `users` |
| A-006 对隔离/事务/R5/迁移的核心实施主张可重复核对 | 上表 + 本轮测试 | 通过。A-006 未核对 create 响应字段名，见 F-001 |

## 对照成功标准

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. 机器凭据与用户 JWT/refresh 分离；secret 只创建返回一次；持久层只存不可逆 hash | **部分** | 分离、hash-only、一次性返回均成立；冻结线名字段是 `secret`，实现与测试锁定为 `token`，见 F-001 |
| 2. 有效/过期/吊销可判定；scope 只从既有 permission key 解析；越权 fail closed | **达成（测试缺口见 F-004）** | 中间件检查 `RevokedAt` 与 `ExpiresAt`；ceiling 三重约束有测试；过期路径无专用用例 |
| 3. 管理 mutation 受权限与 R5 保护；创建/使用/吊销审计含 actor/credential/correlation 且不泄密 | **达成（字面量缺口见 F-002/F-005）** | 同事务 create/revoke；use best-effort；audit detail 无 raw/hash；use detail 缺 scope count |
| 4. 定向与全量相关测试通过；Profile 默认集、模块矩阵、Manifest 装配与协议 pin 不变 | **达成（本轮复验范围内）** | 定向 API 包 + 并发 + kernel/manifest PASS；R6 提交未改协议资产 |

## Findings

### F-001 · create 201 一次性字段名为 `token`，不是冻结的 `secret`

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| 状态 | open |
| 关联 | I-004；成功标准 1；D-002 §4 |
| evidence | D-002 L55：`201；metadata + 一次性 \`secret\``。D-002 L68 将 `secret` 列为不得进入 metadata 的字段名。`handler/service_credentials.go` L205–207：`response["token"] = raw`。`handler/service_credentials_test.go` L33：`created["token"]`。D-003 未改写该字段名。 |
| closure | — |

语义上 raw 只在 create 返回一次、list/detail/audit 不泄密，这一点成立。R6 是横切**契约**目标：D-002 管理 API 表把一次性字段冻成 `secret`。按该契约编写的调用方会读不到 secret。S3 不能在字段名未修订或未获用户书面 residual/overrule 的情况下无条件关门。合法路径：把 JSON 键改回 `secret` 并改测试，或由 `/govern` 落一条决策把线名字段改为 `token`。

### F-002 · 重名 `fieldErrors.reason` 不是 D-003 冻结字面量

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 关联 | I-004；A-002 F-003 / D-003 §3 |
| evidence | D-003 L41：`fieldErrors=[{field:"name", reason:"name already exists"}]`。`handler/service_credentials.go` L198：`Reason: "already in use"`。测试只断言 `INVALID_CREATE_FIELD`（`service_credentials_test.go` L48）。 |
| closure | — |

错误码与 `field` 正确，不扩 catalog。reason 文本与修订契约不一致。不阻断关门，但契约字面量应对齐。

### F-003 · `scopes` 未强制 1～64 上限

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 关联 | D-002 §4（D-003 未覆盖，仍有效） |
| evidence | D-002 L59：`scopes` 为 1～64 个非空唯一 permission key。`handler/service_credentials.go` L152–156 只拒绝空集，无 `> 64` 检查。 |
| closure | — |

catalog / 创建者子集 / 保留键 ceiling 仍 fail closed。缺的是数量上界，不是越权。

### F-004 · user-only「逐项」HTTP 测试与过期凭据用例不完整

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 关联 | D-003 §6；A-001 F-002；成功标准 2 |
| evidence | D-003 L59–65 要求六类自助面「逐项测试机器凭据返回 401」。`service_credentials_test.go` L55–57 只打 `/api/accounts/me`。`auth_test.go` L113–129 覆盖 unknown/revoked，无 expired。过期判定代码在 `auth.go` L571。六处调用点已改 `UserIdentityFrom`（见成果表）。 |
| closure | — |

实现隔离成立；A-006「未知/过期/吊销统一 401」对过期路径是代码审查结论，不是专用测试。补测可降低回归成本，不是未切换 helper。

### F-005 · `service-credentials.use` detail 缺少冻结的 scope count

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 关联 | D-002 §5 |
| evidence | D-002 L67：use 记录 actor、credential id、**scope count**、method/path、correlation。`composition.go` L242–245 与 `handler/testhelpers_test.go` L102–104 只有 `credentialId` / `method` / `path`。Actor 与 correlation 在 operation 行上，不泄密。 |
| closure | — |

## 必改项汇总

**1 条 required（med）**：F-001 — 将 create 201 一次性字段与 D-002 `secret` 对齐，或书面修订契约并留痕。

Recommended（不单独阻断）：F-002 reason 字面量；F-003 scopes 上限；F-004 补测；F-005 use detail `scopeCount`。

## 信息项核对（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持 |
| I-002 | required | S1/S2 实施 | verified | 维持；principal / prefix-before-dev / user-only helper 已实施 |
| I-003 | required | S1 实施 | verified | 维持；0044 DDL、hash-only、无 FK、NOCASE、吊销语义已实施 |
| I-004 | required | S1/S2 实施 | verified | **实施大体成立**；F-001 是线名字段偏差，不把本项打回 `collecting`，但阻断无条件关门 |
| I-005 | required | S3 关门 | verified | **本条复核**：R6 提交未改 Profile/module/Manifest/protocol 源文件；composition/kernel/manifest 测试通过 |
| I-006 | required | S0/S3 审计 | verified | 维持；本条即 S3 independent |

无 `deferred`。无用户书面 `accepted-residual`。无到期未关闭 required 信息项。

## 历史 finding 处置

| 原 finding | 原级别 | 本轮 | 说明 |
|------------|--------|------|------|
| A-002 F-001 | required / high | **保持 fixed** | 0045 correlation-safe CHECK 扩展已落地并有迁移测试 |
| A-002 F-002 | required / med | **保持 fixed** | `created_by` 无 FK；`DeleteUser` 行为未变 |
| A-002 F-003 | required / med | **保持 fixed**（字面量见本条 F-002） | NOCASE UNIQUE + `INVALID_CREATE_FIELD` 已落地 |
| A-002 F-004 | recommended | **实施闭合** | prefix 先于 JWT/`devSession`；有 devSession=true 用例 |
| A-002 F-005 | recommended | **实施闭合** | create/revoke 同事务；use/`last_used_at` best-effort |
| A-002 F-006 | recommended | **实现闭合；测试见 F-004** | 六处 helper 已切换 |
| A-002 F-007 | recommended | **实施闭合** | 32-char id、全状态 list、分页形状与既有 list 一致 |
| A-001 F-001 | recommended | **实施闭合** | 库查询证明 raw 未落盘；并发 revoke 单次转换 |
| A-001 F-002 | recommended | 部分：实现闭合，HTTP 逐项见 F-004 | |
| A-001 F-003 | recommended | **实施闭合** | R5 composition 黑盒 + Profile/Manifest 测试 |

## 与既有意见的异同

| 点 | A-006 self | 本意见 |
|----|------------|--------|
| hash-only / 0044/0045 / 同事务 audit / prefix / ceiling / R5 / 不变式 | pass | 同意；独立复跑定向包与并发测试 |
| create 一次性字段名 | 未点名（写「raw 只在 create 201 返回」） | **F-001 required**：线名是 `secret`，实现是 `token` |
| user-only 六面 | 「已切换 UserIdentityFrom」 | 同意实现；HTTP 只测了 `/api/accounts/me` |
| 过期 401 | 列为已覆盖 | 代码成立；无专用测试 |
| verdict | pass | **conditional** |

不是 P-004.2 冲突：A-006 未宣称线名字段已按 D-002 交付；本条补上契约字面量。

## 结论 + 建议给编排器/用户的下一步

S1/S2 的安全与隔离主张（hash-only、prefix 先分流、human-only、scope ceiling、同事务管理审计、R5、user-only helper、装配不变式）有可重复核对的实现与测试证据。A-006 作为 self 对实施主路径的判断成立，但不能代替本条。

**不能无条件关门**：F-001 required 仍开放。建议 `/govern` 下一句：响应 A-007，将 create 201 字段改为 D-002 的 `secret`（并改测试），或书面接受 `token` 为残余/驳回并修订 D-003；同时决定是否吸收 F-002～F-005。未闭合 F-001 前不得把 GOAL-007 标为 `done`。

## 声明

本意见不修改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 决策或执行正文 / 业务代码。响应由 `/govern` 处理。
