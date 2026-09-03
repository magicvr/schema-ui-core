---
doc_type: goal-audit
id: A-003-r4-independent-audit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-03
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R4 证据矩阵、F-001 必改复核与工作区关门全量审查
audit_type: stage-closeout
verdict: pass
open_required: 0
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-003 · R4 F-001 闭合复核与工作区关门独立交叉审计（independent）

## 1. 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-005-r4-evidence-closeout](../00-meta.md)（R4 证据与关门）；覆盖 Root [GOAL-001-telegram-channel-runtime](../../GOAL-001-telegram-channel-runtime/00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal=GOAL-001-telegram-channel-runtime`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan=VP-030-telegram-channel-runtime` v0.2.1） |
| source | `independent` |
| auditor | grok-build（grok 4.6 · reasoning high） |
| 类型 | `stage-closeout`（finding-closure / R4 关门交叉复审） |
| scope | A-002 **F-001** 必改整改复核；VP-030 退出判据 1～8 与证据矩阵；全量架构红线；真实代码与测试 |
| 对照 | [VP-030](../../../../vision/plans/VP-030-telegram-channel-runtime.md)；GOAL-002 [D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)；GOAL-004 [D-001](../../GOAL-004-r3-outbound-settings-limiter/01-decision/D-001-r3-outbound-and-settings-architecture.md) §2.3；GOAL-005 [A-002](A-002-r4-independent-audit.md)、[A-001](A-001-r4-self-audit.md)、[r4-evidence-matrix.md](../attachments/r4-evidence-matrix.md) |
| 方法 | 只读：工作区绑定 + 五件套 + 指定源码/测试逐行对照；复跑 `apps/api` `go test ./... -count=1`（**PASS**，2026-09-03，本独立审会话）；`git log` 核 Charter 未在本波改动；`git diff` 核 F-001 四处整改。**不**改 status / progress / goal-tree |
| verdict | **pass** |
| 开放 required | **0**（A-002 F-001 闭合证据充分，可重复核对） |

### 范围与区间

- **在 scope**：A-002 F-001（settings 端点鉴权）关闭证据；VP-030 方向级退出判据 1～8；R1/R2/R3 交付与审计闭环；架构红线（Charter / 默认集 / SDK / Redis / Mini App·Stars·FSM·付费命令 / `admin.users` / 内核 import）。
- **不在 scope**：改 Charter / VP status；实施修复；把本意见写成 `done`；把 recommended 项升格为关门阻断。
- **P-005**：I-030-001/002/003/006 `required` 均已在 R1 `verified`。I-030-005/007 在工作区 GOAL-003/004 D-001 为 `verified`；VP-030 正文仍标 `open`——对齐滞后，不构成本 scope 到期 required 信息门禁。无 `deferred required`。
- **共享资料**：无。未读取或比较其他工作区上下文。

### 声明

本意见 `source: independent`，**不**修改目标 `status` / `progress` / 检查点 / 方案正文 / goal-tree。F-001 的正式台账闭合（`fixed`）与是否将 GOAL-005 / Root 标 `done` 由 **`/govern`** 处理。

---

## 2. F-001 闭合核验（代码 + 测试证据）

A-002（independent · `fail` · open required = 1）指出：`GET`/`PATCH /api/channel/telegram/settings` 被登记为 `Public: true`，handler 无 JWT / `settings.read` / `settings.write`，组合根按 contributed handler 原样挂载。模块启用后，未认证调用者可读取掩码密钥并 PATCH 热切换 Bot Token 与 Webhook Secret。

本审对照 A-002 建议修复四条，独立复验工作区未提交 diff 与当前源码。

### 2.1 整改对照表

| A-002 建议修复 | 本审核验 | 判定 |
|----------------|----------|------|
| 1. settings 路由不要 `Public: true`（仅 webhook 保持 Public） | `apps/api/modules/channel/telegram/provider.go`：GET/PATCH settings `Public: false`（L64、L76）；webhook `Public: true`（L90） | **已落地** |
| 2. handler 或装配处包裹与 mail 同构鉴权：读 `settings.read`，写 `settings.write` | `settings_handler.go`：`auth.IdentityFrom` → 无身份 401；GET 缺 `settings.read` → 403；PATCH/PUT 缺 `settings.write` → 403。`composition.go` L605：`tgSettings := a.Middleware(telegraminternal.NewSettingsHandler(tgRuntime))`，与 `mail_admin.go`「Middleware + 权限」同构（mail 在 handler 外包一层；telegram 在 handler 内检查 + 装配处包 Middleware） | **已落地** |
| 3. 未认证 GET/PATCH → 401/403 测试；composition/provider 测试不得再把 settings 当 Public | `runtime_test.go` `TestSettingsHandler_AuthenticationAndPermissions` 覆盖 401 GET、401 PATCH、403 无 `settings.read`、403 只读 PATCH、200 Admin GET（掩码）与 PATCH（热切换且不覆盖未提交 secret）。**缺口（不升 required）**：`provider_test.go` 未断言 `Public: false`；`composition_telegram_test.go` dummy 仍把全部路由标 `Public: true`（dummy，非生产 Provider） | **核心测试已落地**；dummy/`Public` 断言见 R-003 |
| 4. 修复后由 `/govern` 按 `fixed` 闭合 | 本意见只核验关闭证据；正式 `fixed` 留痕由编排器写入 | **证据就绪，待 /govern** |

### 2.2 生产路径防御层（可重复核对）

1. **JWT 门**（装配）：`Authenticator.Middleware` 无 Bearer → 401 `UNAUTHENTICATED`（`auth.go` L590–603，fail-closed）。settings handler 在 `newServer` 中先被 Middleware 包裹再交给 Provider，因此 contributed `mux.Handle` 即使仍不读取 `Public` 字段，settings 也不再以裸 handler 入网。webhook **未**包 Middleware，保持 Public 入站合同。
2. **身份门**（handler）：`IdentityFrom` 失败 → 401。即使有人绕过 Middleware 直调 handler，匿名请求仍被拒绝。
3. **权限门**（handler）：GET 要求 `settings.read`，PATCH/PUT 要求 `settings.write`，否则 403。对齐 mail 管理面权限名。
4. **声明门**（Provider）：settings `Public: false`，仅 webhook `Public: true`。`Public` 字段在 `composition.go` L648–653 装配环仍未被读取——与 A-002 观察相同——但第 1～3 层已关闭原漏洞。不单独重开 required。

### 2.3 测试证据（本审复跑）

`TestSettingsHandler_AuthenticationAndPermissions` 本审核对源码断言：

| 用例 | 期望 | 源码位置 |
|------|------|----------|
| 未认证 GET | 401 | `runtime_test.go` L68–74 |
| 未认证 PATCH（body 含 `bot_token`） | 401 | L76–82 |
| 已认证、无 `settings.read` 的 GET | 403 | L84–92 |
| 仅 `settings.read` 的 PATCH | 403 | L94–102 |
| Admin `settings.read`+`settings.write` GET | 200；掩码末 4 位 | L104–123 |
| 同上 PATCH 只改 token | 200；secret 保持不变 | L125–142 |

本审复跑（2026-09-03，本独立审会话）：

```text
apps/api$ go test ./internal/channel/telegram/ ./modules/channel/telegram/ ./internal/composition/ ./kernel/ -count=1
ok  .../internal/channel/telegram     2.045s
ok  .../modules/channel/telegram      0.707s
ok  .../internal/composition          22.700s
ok  .../kernel                        0.746s

apps/api$ go test ./... -count=1
# 全量 PASS（含 telegram / composition / kernel / handler）
```

工作区 diff（未提交，本审核验内容而非 commit 身份）：

- `settings_handler.go`：+IdentityFrom / 权限检查
- `provider.go`：settings `Public: true` → `false`
- `composition.go`：`NewSettingsHandler` → `a.Middleware(...)`
- `runtime_test.go`：原 `TestSettingsHandler` 扩展为鉴权矩阵

**F-001 关闭结论**：生产路径上「未认证可热切换密钥」已不可复现；401/403/200 有自动化回归。关闭证据充分、可重复核对。建议编排器按 **`fixed`** 合法闭合。无新的 high required。

---

## 3. VP-030 判据 1～8 与架构红线全量核验结论

对照权威：VP-030「方向级退出判据」；合同分母 GOAL-002 D-002；GOAL-004 D-001 §2.3。本审独立复验，不采信矩阵结论。

| 判据 | A-002 结论 | 本审结论 | 说明 |
|------|------------|----------|------|
| #1 Webhook 合同 | PASS | **PASS** | secret fail-closed 与测试仍在；进程内路由已装配 |
| #2 分发端口 | PASS | **PASS** | Register/Unregister + 未知命令回落有测试；R-001 仍为 recommended |
| #3 出站端口 | PASS | **PASS** | stdlib HTTP + mock 降级 + 公共面无 SDK 类型 |
| #4 身份映射 | PASS | **PASS** | `issuer=telegram` 幂等；不写 `admin.users`；不依赖 wallet HTTP |
| #5 设置与密钥 | FAIL（F-001） | **PASS** | 热切换 + 脱敏 + **Admin 鉴权已补**（见 §2）；密钥不进配置包明文 |
| #6 限流评估落盘 | PASS | **PASS** | VRev-070 评估 + 三桶请求计数落地并有测试 |
| #7 边界保持 | PASS | **PASS** | Charter / 默认集 / SDK / Redis / Mini App 红线保持 |
| #8 审计闭合 | FAIL（本审当时 open required=1） | **PASS** | 历史阶段 required=0；A-002 F-001 关闭证据充分；本意见 open required = 0 |

### 3.1 判据 #1 · Webhook 合同 — PASS

合同：secret 校验 fail-closed；无/错 secret 不可当合法 Update；有测试。本审未发现回归。

| 合同项 | 实现 / 测试 | 判定 |
|--------|-------------|------|
| `POST /api/channel/telegram/webhook` | `webhook.go`；Provider 同路径；`composition.go` `plan.HasModule("channel.telegram")` 装配 | PASS |
| 无 token → 503 | `TestWebhook_UnconfiguredToken_Returns503` | PASS |
| 缺/错/空 secret → 401 | `subtle.ConstantTimeCompare`；`TestWebhook_SecretValidation_FailClosed` | PASS |
| 畸形 JSON → 400 | `TestWebhook_MalformedJSON_Returns400` | PASS |
| 合法 Update → 200 | 命令/callback/未知命令用例 | PASS |

Webhook 仍为 `Public: true` 且**未**包 Admin Middleware，符合 D-002 §2「Public，无 Admin JWT」。

### 3.2 判据 #2 · 分发端口 — PASS

Register/Unregister、未知命令 `DefaultTelegramUnknownCommandText`、未知 callback 静默 200、kernel 无 SDK 类型：与 A-002 一致，测试仍绿。R-001（Dispatcher/Sender 未注入其他模块；未启用时无 `ErrTelegramDisabled` stub）保持 recommended，不阻断字面退出。

### 3.3 判据 #3 · 出站端口 — PASS

`HTTPSender` stdlib POST `sendMessage`、10s 超时、无 token 降级 `CaptureSender`、`Validate` fail-closed、按钮仅 `callback_data`：测试名与 A-002 一致，本审全量 `go test` 覆盖对应包。

### 3.4 判据 #4 · 身份映射 — PASS

`GetOrCreateSubject(ctx, "telegram", userID, now)` 幂等；telegram 路径无 `INSERT INTO users`；`composition.go` 独立 `subject.NewStore(st)`，在 `admin.wallet` HTTP 分支之外。`TestWebhook_SubjectMappingIdempotency` 仍在。

### 3.5 判据 #5 · 设置与密钥 — PASS（F-001 闭合后）

相对 A-002 的变化是管理面鉴权：

| 合同项 | 本审 | 判定 |
|--------|------|------|
| Admin 可配置 token/secret | GET `settings.read` / PATCH `settings.write`；未认证 401；未授权 403；Admin 200 热切换 | PASS |
| 密钥 fail-closed（入站） | 空 secret 仍 401 | PASS |
| 只读脱敏 | GET 返回 `token_masked` / `secret_masked`（末 4 位）；PATCH 不覆盖未提交字段 | PASS（R-002 仍建议对齐 mail 只回 `*Set`） |
| 不进配置包明文 | `cmd/schema-ui` 无 telegram 导出字段；`.env.example` 仅注释占位 | PASS（R-005：YAML 仍可持明文，不构成配置包泄漏） |

GOAL-004 D-001 §2.3「需要 Admin 权限或在 Public 外部网络下默认不暴露」现已满足。

### 3.6 判据 #6 · 限流评估落盘 — PASS

VRev-070 §6（进程内够用、不需要 Redis）+ IP 60 / Chat 30 / User 20 + 每次入站 `Record`、永不 `Clear`。`go.mod` 无 Redis 客户端。四条限流测试仍在。

### 3.7 判据 #7 · 边界保持 — PASS

见 §3.9。未发现红线突破。

### 3.8 判据 #8 · 审计闭合 — PASS

| 目标 | 审计台账 | 本审复核 |
|------|----------|----------|
| GOAL-002 | A-001 self `pass`，0 required | 合同 + kernel 端口仍在；阶段关门可维持 |
| GOAL-003 | A-002 independent fail（装配 F-001）→ A-003 `fixed` | BuiltinModules + composition 装配仍在 |
| GOAL-004 | A-001 self `pass`；设置鉴权由 R4 A-002 发现并在本波修复 | 出站/热切换仍在；鉴权缺口已补，不追溯改 R3 status |
| GOAL-005 | A-001 self pass；A-002 independent **fail**（F-001）；本 A-003 复核闭合 | **open required = 0** |

证据矩阵 #8 曾预支「关门双审 0 required」（A-002 R-004）。**当前**独立复审后 required 已归零，判据 #8 字面条件满足。矩阵/执行台账仍建议 `/govern` 回写 A-003 与测试命令（R-004，不阻断）。

### 3.9 架构红线全量核验

| 红线 | 核验 | 判定 |
|------|------|------|
| 未改 Charter | `docs/vision/charter.md` 仍 `schema-ui-core-admin-foundation@0.4.0`，`primary_workspace=workspace-001-mvp-admin-foundation`；最近提交 `1694dea7`（2026-08-31 / 显示 2026-09-01 +0800），工作区文件无 Charter diff | PASS |
| 未进 `mvp` / `admin` / `demo` 默认集 | `profileDefaults` 三档均无 `channel.telegram`；`TestTelegramModule_RegisterContributionsIntegration` 显式断言三档 | PASS |
| 候选集可启用、默认不启用 | `BuiltinModules()` 含 `channel.telegram`；仅 `plan.HasModule` 装配 | PASS |
| 无第三方 Telegram SDK | `apps/api/go.mod` 无 `go-telegram` / `tgbotapi` / `telebot`；入出站均为 stdlib | PASS |
| 无 Redis 依赖 | `go.mod` 无 redis；限流为 VP-027 进程内 | PASS |
| 无 Mini App / Stars / FSM / 付费命令 | 按钮类型禁止 WebApp；无 Stars/Payments/对话引擎/`/buy` | PASS |
| 无独立 Bot 进程 / 长轮询生产 / 多 bot | 同进程 webhook；未见 long-poll 生产路径 | PASS |
| 不污染 `admin.users` | 只写 `subjects (issuer, external_id)` | PASS |
| 内核未直接 import 实现细节 | `kernel/telegram.go` 仅 stdlib；实现 import 仅 composition（装配根，合法） | PASS |
| 不重开 VP-017/026/027/028/029 | 本波 F-001 四处 diff 仅 telegram settings 鉴权，未见历史 VP 回潮 | PASS |
| 密钥不进配置包明文 | 配置包无 telegram 字段；`.env.example` 注释占位 | PASS |

### 3.10 工作区绑定与信息门禁

- `workspace.md`：`id` / `root_goal` / `canonical_scope` / `primary_plan` 一致；共享资料 `none`。
- 工作区文档指针仍写「R4 待立项 GOAL-005」（A-002 I-002），属文档漂移，不阻断关门。
- Root `03-audit.md` 仍为空索引（A-002 I-004）。本意见按用户指定写入 GOAL-005；Root 关门响应应由 `/govern` 在 Root 索引登记本 A-003 的 Q2 链接。

---

## 4. Findings

### Required

无。A-002 **F-001** 关闭证据充分，建议编排器按 `fixed` 合法闭合。

### Recommended（继承 A-002，均不阻断关门）

| ID | 严重度 | 状态 | 本审 |
|----|--------|------|------|
| R-001 | med | open | Dispatcher/Sender 仍未作为可注入端口交给其他模块；未启用时无 `ErrTelegramDisabled` stub。判据 #2 字面仍 PASS。 |
| R-002 | low | open | GET 仍回显密钥末 4 位，弱于 mail `*Set` 布尔。鉴权修复后风险下降，建议后续对齐。 |
| R-003 | low | open | `TestTelegramChannelComposition` 仍用 dummy handler，且 dummy 把 settings 也标 `Public: true`；`provider_test.go` 未锁 `Public: false`。生产路径已由真实 Provider + Middleware + handler 关闭。 |
| R-004 | low | open | GOAL-005 `02-execution` 仍仅 E-001；建议补 E-00N 记录 A-002 fail、F-001 四处整改、本审 `go test ./...` PASS 与本 A-003。 |
| R-005 | low | open | YAML `telegram.bot_token` / `webhook_secret` 仍可持明文；配置包未导出。 |

### Informational

| ID | 说明 |
|----|------|
| I-001 | VP-030 仍将 I-030-005/007 标 `open`；工作区已 `verified`。建议 `/govern` 或 `/vision` 回写。 |
| I-002 | `workspace.md` / Root `00-meta` R4 行仍写「待立项 GOAL-005」。 |
| I-003 | 无 Admin Schema「Bot 渠道 tab」；R3 D-001 已收窄为 HTTP API。 |
| I-004 | Root GOAL-001 `03-audit` 尚无条目。关门时 `/govern` 应在 Root 索引登记本 A-003（Q2）。 |
| I-005 | 本审复跑 `apps/api` `go test ./... -count=1` **PASS**（2026-09-03）。F-001 四处改动目前在工作区未提交 diff 中，内容已核验。 |
| I-006 | `composition.go` 装配环仍不读取 `RouteContribution.Public`。本波靠 Middleware 预包装 + handler 自检闭合 settings；后续若有新的 contributed 管理面，仍可能重复踩坑。不升 required。 |

---

## 5. 必改项汇总

| ID | 严重度 | 建议 | 本审状态 | 是否阻断 GOAL-005 / Root 关门 |
|----|--------|------|----------|------------------------------|
| A-002 F-001 | high | required | **关闭证据充分**（建议 `/govern` 记 `fixed`） | **否**（闭合后） |

开放 required = **0**。无到期未关闭的 required 信息项。recommended 不构成 P-003 关门门禁。

---

## 6. 与既有意见的异同

| 来源 | 异同 |
|------|------|
| GOAL-005 A-002 independent `fail` | **同意**当时 F-001 成立。本审复核整改后**撤销阻断**：判据 #5/#8 现为 PASS。不重开装配类 R2 F-001。 |
| GOAL-005 A-001 self `pass` | 当时未打开 Provider `Public` 与 mail 鉴权先例，过宽。经 A-002 指出并修复后，自审对 #1–#4/#6/#7 的判断可维持。 |
| 证据矩阵 / D-001 | #1–#4、#6、#7 仍可复核。#5 现因鉴权补齐而与矩阵 PASS 对齐；#8 在本独立复审后才真正满足「双审 required=0」。 |
| GOAL-004 A-001 self | 出站/热切换/限流核账仍同意 PASS；设置鉴权已由 R4 补上。 |

无 P-004 意见冲突：仅一条 required，且关闭证据单向成立。

---

## 7. 综合结论与放行 GOAL-005 及 Root GOAL-001 关门建议

**verdict: pass。** A-002 唯一 required F-001（settings 端点无 Admin 鉴权）已在生产路径上以三层防护闭合（Middleware JWT + handler `IdentityFrom` + `settings.read`/`settings.write`），401/403/200 有自动化回归，全量 `apps/api` 测试绿。VP-030 判据 1～8 现均可核对为达成；架构红线保持。开放 required = 0。

**放行建议（给编排器 / 用户）**：

1. 用 **`/govern`** 合并响应本 A-003：将 A-002 F-001 记为 **`fixed`**（证据指向 settings handler / Provider `Public: false` / `a.Middleware` 挂载 / `TestSettingsHandler_AuthenticationAndPermissions` + 本审 `go test ./...` PASS）。
2. R-001～R-005 可修可记接受，**不阻断**本波关门。
3. **可以**将 GOAL-005 与 Root GOAL-001 标为 `done`（检查点收口、goal-tree 同步、Root `03-audit` 登记本 A-003 的 Q2 链接）。独立审本身不改 status。
4. 建议顺手：补 GOAL-005 执行事实 E-00N；回写 VP-030 I-030-005/007 与关门记录；修正 workspace/Root「待立项 GOAL-005」指针。
5. F-001 四处代码目前为未提交工作区改动——关门响应时应纳入显式 owned paths 的 Git checkpoint（若任务启用），避免只关门文档、丢失鉴权修复。

建议编排器下一句：

```text
/govern 响应 GOAL-005 A-003（independent pass，F-001 关闭证据充分）：将 A-002 F-001 记 fixed；开放 required=0；放行 GOAL-005 与 Root GOAL-001 关门，并在 Root 03-audit 登记本意见 Q2 链接
```
