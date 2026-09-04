---
doc_type: goal-audit
id: A-006-r2-c2-implementation-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: Grok
provider: Grok
audit_type: execution-facts
scope: R2 C2 实现（telegram v67 additive migration、runtime 回读权威性、settings PATCH 部分更新与 persist-then-memory、mode/URL 校验、secret 加密与 export/status 暴露、identity/catalog/fresh/reopen 与 C2 composition/config/runtime 测试；对照 D-001、A-003、A-004、A-005）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-006 · R2 C2 实现独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：Grok
- **provider**：Grok
- **类型** / **scope**：execution-facts · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 C2 实现。只核 `apps/api/internal/config/config.go`、`apps/api/internal/channel/telegram/runtime.go`、`settings_handler.go`、`apps/api/modules/channel/telegram/migration/migration.go`、`apps/api/internal/composition/composition.go`、相关 migration identity/catalog 测试与 C2 composition/config/runtime 测试；对照 D-001、A-003、A-004、A-005
- **verdict**：pass
- **open required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001～A-005 原文均未改写。未复跑测试；证据来自当前文件直接核对。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：v67 是否 additive 且 fresh/upgrade/reopen/identity/catalog 一致；已有 DB row（含 mode/URL 为空）是否绝对优先于 seed；PATCH 部分更新是否仅在持久化成功后发布内存；mode/URL 校验、secret 加密、export/status 暴露；C2 兼容性与关门遗漏
- **excluded**：C3 Bot API / connection manager / Fx `OnStop` drain；C4 Admin UI / heartbeat；C5 Fake Bot API 矩阵；全量测试复跑；其他工作区；把 `progress: 1/5` 或 A-005 self `pass` 当作完成证据

## 信息项与门禁

| 项 | 状态 | 本条核对 |
|----|------|----------|
| I-033-014（required，最晚 C1） | verified | C1 已裁决；本条核的是 C2 代码是否落实 D-001「行存在后 DB authoritative、仅无行 seed、PATCH 部分更新」 |
| I-033-015～016 | verified | 最晚 C3；不在本 C2 scope |
| I-033-017～018 | non-blocking open | 最晚 C3；不构成本条必改 |
| 到期 required 信息项 | 无 | 无 residual / overruled |
| 共享资料 | none | 未当事实 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| v66 DDL 未改写；v67 只 `ALTER TABLE … ADD COLUMN` | `migration.go` L11–16（v66 仍无 mode/URL）、L28–31、L69–91 |
| v66 checksum 盐与冻结值仍为 `telegram_config`；v67 新 identity | `migration.go` L72–73、L76–86；`migrate_test.go` L739–741（v66 `e330d785…` 未改，v67 `50eb9a44…` / `telegram_config_connection`） |
| fresh / reopen / catalog head / identity 指向 v67 且无新表 | `migrate_test.go` L124–126、L192–193；`restart_test.go` L52–53；`identity.go` L93、L96–97、L113；`identity_test.go` L118 `67: {}`；`provider_test.go` L63 |
| 无行 → YAML/env seed；有行 → 读库覆盖内存 seed | `runtime.go` L59–60、L105–107、L116–128、L131–156；`composition.go` L853 |
| 空 mode/URL 的既有行不回落 seed | `runtime.go` L142–154；`composition_telegram_test.go` L378–427（UPDATE 成空后重启，status 为 polling + 空 URL，忽略 stale webhook seed） |
| PATCH 指针部分更新；未提供字段保留当前值 | `settings_handler.go` L22–28、L76–107 |
| 持久化成功后才写内存 | `runtime.go` L243–278（`Run` 失败在 L266–268 返回，L271–277 内存更新在其后） |
| mode 仅 polling/webhook；空 mode 规范为 polling | `config.go` L303–316、L656–659、L785–793；`runtime.go` L160–175 |
| 非空 URL 必须是绝对 http(s) origin（无 path/query/fragment/credentials/空白） | `config.go` L1336–1353、L794–797；`runtime.go` L177–188；`config_test.go` L204–214 |
| token/secret 仍 `EncryptSecret`；mode/URL 明文非敏感列 | `runtime.go` L118–125、L244–263 |
| Status / GET 不返回 token/secret 原文 | `runtime.go` L21–30、L282–301；`settings_handler.go` L55–61、L114–115 |
| 导出树只有 mode/URL，不含 bot_token / webhook_secret | `configpkg.go` L88–91、L237–238、L248–253 |
| composition 使用带 settings 的构造器 | `composition.go` L853 |
| C2 未声称 C3/C4 已完成 | A-005 L27–29；E-004 L18；本条同意 |

## 对照成功标准（C2 适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| D-001 / I-033-014：无行 YAML/env seed；行存在后 DB authoritative（空值也算）；PATCH 可改 mode/URL | **已达成（C2 代码）** | `runtime.go` L116–156；`settings_handler.go` L91–109；`composition_telegram_test.go` L326–376、L378–427 |
| v67 additive；v66 保留；fresh/reopen/identity/catalog 一致 | **已达成** | 见上表；`lockedHeadExtraTables[67]` 为空切片符合「无新对象的 additive ALTER」 |
| mode/URL 非法 fail closed；不复用 `runtime.mode` / `auth.public_base_url` | **已达成（装载与 PATCH 格式）** | `config.go` L246–252、L790–797；`config_test.go` L171–216；非法 mode `long-polling` 与带 path/query 的 URL 置 `LoadError` |
| token/secret 加密；export/status 不泄漏密钥 | **已达成** | `runtime.go` L118–125、L293–301；`configpkg.go` L88–91、L251–252 |
| A-003 F-001（已有行 + 新列为空不得当第二次 YAML seed） | **代码已落实**（finding 仍 open，待 `/govern` 标 closed） | `runtime.go` L116–129 仅 `count == 0` 才 INSERT seed；有行走 L131–156 |
| webhook 缺 URL/secret 不得 `setWebhook` / 进入 `running` | **不在 C2** | D-002 L51、D-003 L22；C2 允许把不完整 webhook 行写入 DB（见 F-001） |
| C3 Bot API / connection manager / 40s polling client | **不在本条 scope** | composition 仍把 `Manager` 接到密钥热切换 `RuntimeManager`（`composition.go` L832、L853） |

## 四点核对

### 1) v67 additive 且 fresh/upgrade/reopen/identity/catalog 一致

**结论：成立。**

- v66 `CREATE TABLE telegram_config` 仍只有 `id/bot_token_enc/webhook_secret_enc/updated_at`（`migration.go` L11–16、L19–26）。未把 mode/URL 塞进 v66，checksum 盐仍是 `"0066:telegram-config:v1"`（L72）。
- v67 是两条 `ADD COLUMN`：`mode TEXT NOT NULL DEFAULT 'polling'`、`webhook_public_base_url TEXT NOT NULL DEFAULT ''`（L28–31）。SQLite/PG 共用同一 ALTER 列表（L51–66）。这是 additive，不是 rebuild。
- 升级已有 v66 行时，SQL DEFAULT 会把 mode 写成 `polling`、URL 写成 `''`，然后 runtime 按「有行」读库、忽略 YAML。这与 A-003 F-001「空/polling、忽略 YAML URL」及 I-033-012 缺省 polling 相容；不是第二次 seed。
- catalog 连续编号 67、Key/Name `telegram_config_connection`、冻结 checksum `50eb9a447a0b2ebf1dd77b4f780dcbedc7cd38300430f88a4ddf217149797182`（`migrate_test.go` L741）。
- `completeFingerprintCatalogHead = 67`（`identity.go` L93）；`lockedHeadExtraTables[67] = {}`（`identity_test.go` L118）正确记录「无新表」。`completeLostLedgerTables` 仍只含 `telegram_config`（`identity.go` L113）。
- fresh 与 reopen 断言 `len==67` 且尾为 v67（`migrate_test.go` L124–126、L192–193；`restart_test.go` L52–53）。provider 贡献 `v66 + v67`（`provider_test.go` L63）。

未发现把 v66 改义、跳号、或把 v67 登记成新表的 identity 缺口。缺少「停在 v66 已有 token 行再 apply v67」的集成测试，见 F-002，不升级为 required。

### 2) 已有 DB row（含空 mode/URL）绝对优先于 seed；PATCH 部分更新仅在持久化成功后发布内存

**结论：值覆盖路径成立；persist-then-memory 在 `UpdateSettings` 成立。**

- `count == 0` 才用 seed INSERT（`runtime.go` L116–128）。`count > 0` 解密并采用 DB 的 token/secret/mode/URL，注释写明 including empty（L105–107、L131–156）。
- 空 mode 经 `normalizeTelegramMode` 成为内存中的 `polling`（L142–154、L160–165），**不** overlay YAML。`TestTelegramRuntime_EmptyConnectionSettingsRemainAuthoritative` 把已有行 UPDATE 成空后用 webhook/stale URL 重启，断言 `Mode==polling` 且 URL 为空（`composition_telegram_test.go` L397–426）。
- 非空权威：`TestTelegramRuntime_ConnectionSettingsPersistenceAndAuthority` 在 seed 后 `UpdateSettings` 再以 stale seed 重启，读到 live 值（L348–375）。
- composition 把 YAML/env 只当构造参数传入（`composition.go` L853），权威在 `initPersistence`。
- PATCH：`updateSettingsRequest` 四个指针字段（`settings_handler.go` L22–28）；`nil` 保留 `Get*` 当前值（L76–107）；然后只调用 `UpdateSettings`（L109）。
- `UpdateSettings`：先 `updateMu`，先校验，有 runner 时先加密再 `Run` UPSERT；`err != nil` 则 return，内存赋值在成功之后（`runtime.go` L228–278）。

限制（recommended，非本条必改）：

- 构造器在读库前校验 seed mode/URL（`runtime.go` L68–75）。无效 YAML 会在权威行被读取前让启动失败。生产 `config.Load` 同样 fail closed（`config.go` L790–797），所以这是装载纪律，不是「有效 seed 覆盖有效 DB」。
- 部分更新的读合并在 `updateMu` 之外（`settings_handler.go` L76–79 vs L109）。单请求 persist-then-memory 成立；并发 PATCH 可能丢失字段（F-005）。
- HTTP PATCH + DB 重启的证据在 composition 的 `UpdateSettings` 测试，不在 `runtime_test.go`（该测试 `NewRuntimeManager` 无 runner，L63）。persist 失败不改内存无直接测试（F-003）。代码顺序本身可核对。

### 3) mode/URL 校验、secret 加密、导出/status 暴露

**结论：C2 格式校验、加密与暴露边界成立。**

- 装载：空 mode → polling；非法 mode / 非法 origin → `LoadError`（`config.go` L656–659、L785–797、L1254–1258；`config_test.go` L192–214）。Telegram mode 独立于 `runtime.mode`（`config.go` L246–248）。
- PATCH：非法 mode → 400 `invalid telegram mode`；非法 URL → 400 `invalid webhook public base URL`（`settings_handler.go` L92–106）。空 URL 合法（`runtime.go` L177–181；`config.go` L1342–1344），与 D-001「polling 可不配置 URL」一致。
- webhook + 空 URL / 空 secret **可被持久化**。D-002 L51 / D-003 L22 的 fail-closed 在 `setWebhook`，属 C3（F-001）。
- token/secret：`mail.EncryptSecret`；空明文加密为空串（`secrets.go` L59–62；`runtime.go` L118–125、L244–251）。mode/URL 不加密。
- Status JSON 字段为 `configured/token_set/secret_set/mode/webhook_public_base_url/captured_*`，无 token/secret 键（`runtime.go` L22–30、L293–301）。GET/PATCH 都 `Encode(Status())`。
- 导出 `cfgTree.Telegram` 仅 mode 与 URL（`configpkg.go` L88–91）。`telegram.bot_token` / `telegram.webhook_secret` 登记在 sensitive 表（L251–252），默认 YAML 无这两键故 `TestExportDefaultShape` 仍期望 `secrets.exclude == 2`（`configpkg_test.go` L59–67）。值不会进入导出树。`TELEGRAM_MASTER_KEY` 不在 YAML/导出树。

### 4) 兼容性与 C2 关门遗漏

**结论：C2 合同范围内无未闭合 high/med required；C3 门禁不得被本条放行。**

已覆盖 C2 成功标准中的配置 schema、递增 migration、runtime 回读、settings API 与密钥边界。下列不是 C2 必改，但是关门前不得漏看：

- A-003 F-002/F-003、I-033-015/016、GOAL-002 D-003 的 30s/40s client 与独立 connection manager：C3。
- webhook 缺 URL/secret 不得进入 `running`：C3（本条 F-001）。
- I-033-017/018、Admin UI：C3/C4。
- A-003 F-001 代码已落实，ledger 仍为 recommended open，需 `/govern` 标 `closed/fixed`。
- 本条不把 A-005 的 git SHA 或 `progress: 1/5` 当作闭合证据。

## Findings

### F-001 · C2 允许持久化「mode=webhook 且 URL/secret 为空」；fail-closed 仍在 C3

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | D-001 L22；D-002 L51；D-003 L22 |
| evidence | `runtime.go` L177–181、L234–241；`settings_handler.go` L91–107；`config.go` L1342–1344 |

空 URL 在校验函数中直接 `return nil`。PATCH 可以把 mode 改成 `webhook` 而 URL/secret 仍空。这与「PATCH 部分更新」相容，也与 D-002「缺 URL 则不调用 `setWebhook`、不进入 `running`」相容——那是连接建立，不是 settings 行约束。C3 必须按 D-002/D-003 fail closed，不能因为 C2 能存下该行就标 `running`。不构成本条 required。

### F-002 · 缺少 v66 已有 token 行再 apply v67 的升级集成测试

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-003 F-001 |
| evidence | `migration.go` L28–31；`composition_telegram_test.go` L378–427 |

空列权威测试是在 **已经是 v67 的库** 上 `UPDATE mode/url=''`，不是「v66 行 + ALTER」。代码路径（`count > 0` 读库）与 ALTER DEFAULT（`polling` / `''`）足以支持升级语义，但没有把「旧行加列」钉成测试。建议 C5 或补测时加一条；不阻断 C2。

### F-003 · HTTP PATCH 落库+重启与 persist 失败不改内存缺少直接测试

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | `runtime_test.go` L63、L145–158（无 `TxRunner`）；`composition_telegram_test.go` L348（直接 `UpdateSettings`）；`runtime.go` L266–277 |

`UpdateSettings` 的成功路径顺序可从源码重复核对。失败 runner、以及 `handleUpdate` → DB → 进程重启的端到端测试不在当前 C2 测试里。不因此否定 persist-then-memory 实现。

### F-004 · 导出/serve 装载不校验 telegram.mode/URL；默认 exclude 测试未锁定 Telegram 密钥键名

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | `server/config.go` L105–113（解析但不拷贝、不校验）；`configpkg.go` L237–238、L248–253；`configpkg_test.go` L66–74 |

生产 API 走 `internal/config.Load`，非法 mode/URL fail closed。`schema-ui` 导出走 `server.LoadConfig`，非法 telegram 字段仍可打进包。默认导出树没有 token/secret 字段，泄漏面不在值拷贝。建议 overlay 导出测试断言 `bot_token`/`webhook_secret` 键名不出现在 `config` 段。

### F-005 · 部分更新读合并在 `updateMu` 之外

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | `settings_handler.go` L76–109；`runtime.go` L226–230 |

`UpdateSettings` 注释称串行化可避免 mixed snapshot。它串行化的是「整行写入 + 内存发布」，不是 handler 里对当前值的读合并。并发 PATCH 一个改 mode、一个改 token 时，后写者可能用过期 mode/URL 整行覆盖。单请求语义仍符合 D-001。C3 热切换前可把读合并移进同一把锁。

## 必改项汇总

无。本条 **open required = 0**。

无到期且影响 C2 的 required 信息项。无 `accepted-residual` / `user-overruled`。未发现与 D-001「行存在后空值也权威、仅无行 seed、PATCH 部分更新、token/secret 加密」相反的实现。

## 与既有意见的异同

| 项 | A-003 independent | A-004 self response | A-005 self | 本条 independent |
|----|-------------------|---------------------|------------|------------------|
| 原文是否保留 | C1 pass / open required 0 | 响应 A-003，放行 C2/C3 入口 | C2 self pass / open required 0 | 未改写 A-001～A-005 |
| A-003 F-001 已有行空列不得 YAML overlay | recommended，待 C2 测试 | 转入 C2 计划 | 称 composition 回归已固定 | **代码+测试已落实**；ledger 仍 open，待 `/govern` 闭合 |
| A-003 F-002/F-003 | C3/C4 recommended | 同左 | 明确未做 | 同意：不在 C2 完成面 |
| C2 生产实现 | 当时未开始 | 允许开工 | 称已完成 | **同意 C2 合同已落地**；不放行 C3 |
| verdict | pass（C1） | pass（响应） | pass（C2 self） | **pass（C2 independent）** |
| open required | 0 | 0 | 0 | **0** |

本条与 A-005 在 v67 additive、DB 空值权威、persist-then-memory、密钥不进 status/export 上一致。差异：本条补了 webhook 不完整行、升级测试形态、导出校验、并发读合并等 recommended，且不以 A-005 的 checkpoint SHA 或未复跑的测试日志为证据。A-005 的 `pass` 与本条 `pass` 不是 P-004 冲突。

## 结论 + 建议给编排器/用户的下一步

R2 C2 实现可重复核对：**verdict = pass**，**open required = 0**。v67 为 additive；已有行（含空 mode/URL）不回落 seed；`UpdateSettings` 在持久化成功后才发布内存；mode/URL 格式校验与密钥暴露边界成立。

建议 `/govern`：

1. 响应本条：I-033-014 保持 `verified`；将 A-003 F-001 标为 C2 `closed/fixed`（代码+`EmptyConnectionSettingsRemainAuthoritative`）。
2. 不要把本条 recommended（F-001～F-005）当成 C2 必改阻断；F-001 必须在 C3 `setWebhook` 路径 fail closed。
3. 仅在用户确认后把 GOAL-003 C2 检查点标完成并重算 `progress`。本独立审不改检查点。
4. 不要把本条当作 C3 放行。

## 声明

本意见 `source: independent`，`auditor`/`provider` 为 Grok。不修改 status/progress/检查点/goal-tree/decision 正文或生产代码。响应与是否将 C2 标完成由 `/govern` 处理。
