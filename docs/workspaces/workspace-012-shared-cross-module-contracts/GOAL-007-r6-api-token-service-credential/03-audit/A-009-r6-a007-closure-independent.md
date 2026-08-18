---
id: A-009-r6-a007-closure-independent
goal: GOAL-007-r6-api-token-service-credential
doc: audit-entry
record_id: A-009
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-007 F-001～F-005 finding-closure；提交 b6ebfec、E-006、A-008；现行代码/测试与 D-002/D-003 不变式
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
responds_to: A-007
reviews: A-008
---

# A-009 · A-007 F-001～F-005 关闭复核（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：ad-hoc / finding-closure
- **scope**：GOAL-007 R6；仅复核 A-007 F-001～F-005 是否已由提交 `b6ebfec`、E-006、A-008 与现行代码/测试按 P-003 `fixed` 可靠闭合。核对 D-002 未被覆盖条款（create 201 `secret`、scopes 1～64、use `scopeCount`）与 D-003 §3/§6（重名字面量、六类 user-only 逐项 401）。
- **verdict**：**pass**
- **required findings**：0

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 覆盖本目标；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：A-007 原文 F-001～F-005；A-008 候选响应；E-006；提交 `b6ebfec`（`b6ebfec824b959d0c330f6f3f3673d379d3d57b5`）及 HEAD 上现行实现/测试；D-002 §4/§5（`superseded`，未被 D-003 覆盖的条款仍有效）；D-003 §3/§6。
- **excluded**：不改 `status` / `progress` / `00-meta` / D-003 / E-* / goal-tree / Root / 业务代码；不重做 A-007 全量 S3 close-out；不读取或比较其他工作区上下文；不把 A-008 self 当作独立闭合。
- **共享资料**：无固定引用；不得当作事实或 finding 关闭依据。
- **本轮复验**：独立复跑 `apps/api` `go test -timeout 15m -count=1`：
  - `./internal/handler` `-run TestServiceCredentialManagementAndAuthentication|TestServiceCredentialCreateRejectsUnknownReservedAndExcessScopes|TestServiceCredentialRequiredAuditFailureRollsBack` — ok 5.013s
  - `./internal/auth` `-run TestServiceCredentialMiddlewarePrecedesDevFallback|TestNewServiceCredentialTokenContract` — ok 2.228s
  - `./internal/auth` 全包 — ok 15.582s
  - `./internal/composition` 全包 — ok 16.331s
  未重跑 `go test ./...` 全量；finding-closure 以锁定 F-001～F-005 的定向用例与相关包为准。`b6ebfec..HEAD` 未再改这 6 个实现/测试文件；其后 `dd39c94` 仅为文档响应。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-007 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R6 / VP-012 / Charter 的明显冲突 | 整改只对齐已冻结线名与测试，未新增 Profile/module/page/nav/fragment |
| Vision Review required | 本 scope 未见开放 required | 本意见不写 `docs/vision/reviews.md` |
| 既有 Goal 审计 | A-001～A-006 无开放 required；A-007 independent = conditional（F-001 required/med；F-002～F-005 recommended）；A-008 self = pass（实现侧 proposed fixed，不自闭） | `03-audit.md` |
| P-004 冲突 | 无待裁冲突 | A-008 对五条全部走 `fixed`，无 residual/overrule，也无一要一否 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 提交 `b6ebfec` 只改 6 个实现/测试文件，信息与 A-008/E-006 一致 | `git show --stat b6ebfec`：`auth.go`、`auth_test.go`、`composition.go`、`service_credentials.go`、`service_credentials_test.go`、`testhelpers_test.go`（+62/−12） | 通过。信息 `fix(workspace-012): close R6 audit findings` |
| HEAD 现行代码仍是该整改 | `git diff --name-only b6ebfec HEAD` 对上述 6 文件为空；后续 `dd39c94` 为 docs | 通过 |
| create 201 一次性字段为冻结的 `secret`；`token`/`tokenHash` 不出现 | `handler/service_credentials.go` L205–207；`service_credentials_test.go` L35–36 | 通过。本轮定向 handler 测试 ok |
| list/detail 元数据不含 raw/`secret`/`tokenHash` | 同测试 L40–46；`serviceCredentialRow`/`serviceCredentialAuditRow` 无这些键 | 通过 |
| 重名 `fieldErrors.reason` 为 D-003 字面量 | `service_credentials.go` L198；`errorcatalog.FieldError` `json:"reason"` | 通过实现。测试仍只锁 `INVALID_CREATE_FIELD`，见 F-002 处置 |
| scopes 在 normalize 后拒绝空集与 `> 64` | `service_credentials.go` L152–155；`service_credentials_test.go` L127–137（65 个唯一 `permission.%02d` → `INVALID_CREATE_FIELD`） | 通过。计数检查在 catalog 校验之前，故 65 未知 key 不会先变成 `INVALID_PERMISSION_REF` |
| 六类 user-only HTTP 表面机器凭据均为 401 | 测试 L57–71；夹具 `newWalletSelfEnv` + `mountMFASurface` 实际挂上六条路由；六处 helper 仍为 `UserIdentityFrom` | 通过。本轮定向测试 ok |
| expired credential 专用 401 | `auth.go` L572；`auth_test.go` L131–149（`ExpiresAt: now.Add(-time.Minute)`，`called=false`） | 通过。本轮 auth 定向与全包 ok |
| use audit detail 含 `scopeCount`；生产 composition 与测试 recorder 同形 | `auth.go` L77–80、L586–589；`composition.go` L241–247；`testhelpers_test.go` L101–104；`service_credentials_test.go` L99–100（`"scopeCount":1`） | 通过。`scopes:["users.read"]` 与 count=1 一致；composition 全包 ok |
| 整改未改 hash-only / prefix-before-dev / 0044/0045 / 同事务 audit / R5 / 装配源 | `b6ebfec` 不含 migration、`operational.go`、kernel/manifest/protocol | 通过。本条不重审 A-007 已通过的非 finding 主张 |

## 对照成功标准（本轮 scope = A-007 finding 闭合）

| 标准 / 关闭条件 | 本轮 | 证据 |
|-----------------|------|------|
| F-001：create 201 一次性字段与 D-002 `secret` 对齐 | **fixed** | 实现 + 测试锁定 `secret`，并否定 `token`/`tokenHash` |
| F-002：重名 `fieldErrors.reason` = `name already exists` | **fixed** | handler 字面量与 D-003 §3 一致；`FieldError.reason` JSON 标签正确 |
| F-003：`scopes` 强制 1～64 个唯一 permission key | **fixed** | normalize 后空集与 65 上界均 `INVALID_CREATE_FIELD` |
| F-004：六类 user-only 逐项 HTTP 401 + expired 401 | **fixed** | 六路由黑盒 + 过期中间件用例 |
| F-005：`service-credentials.use` detail 含 `scopeCount` | **fixed** | auth 投影、composition/test recorder、handler 审计断言 |
| 成功标准 1 的线名缺口 | **已补** | A-007「部分」仅因 `token`/`secret`；本轮字段名已对齐。本条不替代完整 S3 close-out |
| I-004 线名字段偏差 | 实施已对齐；登记维持 `verified` | 不把本项打回 `collecting` |

## A-007 disposition

### F-001 · create 201 一次性字段名为 `token`，不是冻结的 `secret` — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | med / **required** |
| 状态 | **fixed** |
| closure | 提交 `b6ebfec`；E-006 §1；A-008 表 F-001；未走 residual/overrule |
| evidence | D-002 L55 冻结 `secret`。`handler/service_credentials.go` L206：`response["secret"] = raw`。`service_credentials_test.go` L35–36：读 `created["secret"]`，并要求 `created["token"]` / `created["tokenHash"]` 为 nil。L43：list/detail 不得含 raw、`tokenHash` 或 `"secret"` 键。本轮定向 handler 测试通过。 |

语义上 raw 仍只在 create 返回一次。契约线名现与 D-002 一致；按该契约读 `secret` 的调用方可拿到一次性值。

### F-002 · 重名 `fieldErrors.reason` 不是 D-003 冻结字面量 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | low / recommended |
| 状态 | **fixed** |
| closure | `b6ebfec`；E-006 §2；A-008 表 F-002 |
| evidence | D-003 L41：`fieldErrors=[{field:"name", reason:"name already exists"}]`。`service_credentials.go` L198：`Reason: "name already exists"`，错误码仍为既有 `INVALID_CREATE_FIELD`。`errorcatalog.FieldError` 的 `Reason` 序列化为 `reason`。 |

观察（不新开 finding）：`service_credentials_test.go` L48–50 仍只 `expectError(..., INVALID_CREATE_FIELD)`，未锁 `fieldErrors.reason`。原 finding 是实现字面量偏差，现已与 D-003 对齐，可由源码独立复核。

### F-003 · `scopes` 未强制 1～64 上限 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | low / recommended |
| 状态 | **fixed** |
| closure | `b6ebfec`；E-006 §2；A-008 表 F-003 |
| evidence | D-002 L59（D-003 未覆盖，仍有效）。`service_credentials.go` L152–155：`normalizedCredentialScopes` 之后 `len == 0 \|\| len > 64` → `INVALID_CREATE_FIELD`。`service_credentials_test.go` L127–137：65 个互异 key，期望 `400 INVALID_CREATE_FIELD`。计数检查先于 `ValidatePermissionKeys`，因此该用例锁的是数量上界而不是 catalog。 |

catalog / 创建者子集 / 保留键 ceiling 未被本整改改写。

### F-004 · user-only「逐项」HTTP 测试与过期凭据用例不完整 — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | low / recommended |
| 状态 | **fixed** |
| closure | `b6ebfec`；E-006 §4；A-008 表 F-004 |
| evidence | D-003 §6 六类清单。`service_credentials_test.go` L57–71 逐项：`GET /api/accounts/me`、`GET /api/account/profile`、`POST /api/account/avatar`、`GET /api/mfa/status`、`GET /api/notifications`、`GET /api/wallet/me`，均为 `401 UNAUTHENTICATED`。夹具改为 `newWalletSelfEnv` + `mountMFASurface`，六条路由均已挂载。实现仍走 `UserIdentityFrom`：`account.go` L22、`account_self.go` L83、`account_avatar.go` L51、`mfa.go` L231、`notifications.go` L65、`wallet_self.go` L107。`auth_test.go` L131–149 补 expired：过期 secret 401 且 `called=false`。本轮 handler 定向与 auth 包通过。 |

「六类」按 D-003 文件/表面各取一条代表路由；同一 helper 下的 password/sessions 等兄弟路径不另开 finding。

### F-005 · `service-credentials.use` detail 缺少冻结的 scope count — **fixed**

| 字段 | 值 |
|------|-----|
| 原级别 | low / recommended |
| 状态 | **fixed** |
| closure | `b6ebfec`；E-006 §3；A-008 表 F-005 |
| evidence | D-002 L67：use 记录 actor、credential id、**scope count**、method/path、correlation。`auth.go` `ServiceCredentialUse.ScopeCount`；`authenticateServiceCredential` L589：`len(credential.Scopes)`。生产 `composition.go` L244 与测试 `testhelpers_test.go` L103 均写入 `scopeCount`。`service_credentials_test.go` L99–100 断言 use detail 含 `"scopeCount":1`。审计 detail 仍禁止 raw secret。本轮 handler 定向与 composition 包通过。 |

## Findings

本轮 **无新 required / recommended finding**。A-007 F-001～F-005 全部按 P-003 `fixed` 合法闭合。

## 必改项汇总

**无。** required = 0。

## 信息项核对（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持 |
| I-002 | required | S1/S2 实施 | verified | 维持；F-004 补测不改变已 verified 的 principal/隔离结论 |
| I-003 | required | S1 实施 | verified | 维持 |
| I-004 | required | S1/S2 实施 | verified | **维持**；F-001 线名字段现与 D-002 对齐，不打回 `collecting` |
| I-005 | required | S3 关门 | verified | 本条不重做 S3 装配/协议关门复证；整改提交未改 kernel/manifest/protocol |
| I-006 | required | S0/S3 审计 | verified | 本条为 A-007 的 independent finding-closure，不是第二次完整 S3 close-out |

无 `deferred`。无用户书面 `accepted-residual`。无到期未关闭、影响本 scope 的 required 信息项。

## 与既有意见的异同

| 点 | A-007 | A-008 | 本意见 |
|----|-------|-------|--------|
| F-001 | required / open | proposed fixed（`secret`） | **fixed**；独立核对源码、测试与 `b6ebfec` |
| F-002 | recommended / open | proposed fixed | **fixed**（实现字面量）；测试仍只锁错误码，不新开 finding |
| F-003 | recommended / open | proposed fixed（65-scope） | **fixed** |
| F-004 | recommended / open | proposed fixed（六面 + expired） | **fixed**；夹具确挂六路由 |
| F-005 | recommended / open | proposed fixed（`scopeCount`） | **fixed**；生产 composition 与测试 recorder 同形 |
| verdict | conditional | pass（self 响应，不自闭） | **pass**（closure） |
| S3 关门 | 因 F-001 不能无条件关门 | 明确仍待本条 | 本条只闭合 finding；是否关 S3 / `done` 由 `/govern` 处理 |

不是 P-004.2 冲突：A-008 未宣称已独立闭合；本条确认其候选成立。

## 结论 + 建议给编排器/用户的下一步

A-007 唯一 required finding F-001，以及四条 recommended findings，均有可重复核对的实现、提交与本轮复跑测试。关闭路径均为 `fixed`，无 residual/overrule。

本条 **不是** 第二次完整 S3 close-out，也不改目标状态。建议 `/govern` 下一句：响应 A-009，将 A-007 F-001～F-005 标为 `fixed`，再决定是否关闭 S3 / 将 GOAL-007 标为 `done`（须另核对 I-005 关门证据与 A-007 已通过的非 finding 主张，不得只用本条 progress 或本条 verdict 替代）。

## 声明

本意见不修改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 决策或执行正文 / 业务代码。响应由 `/govern` 处理。
