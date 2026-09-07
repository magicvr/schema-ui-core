---
doc_type: goal-audit
id: A-014-independent-closure-review-f007-f008
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: open
source: independent
auditor: deepseek-v4-flash (dsh local independent audit)
audit_type: finding-closure
scope: focused finding-closure 复审 · A-012 两项 required findings（F-007 迁移 Apply 中途失败零残留 + corrected reopen；F-008 Telegram-enabled 真实组合根 + 结构化 Manifest）的关闭证据，核对 A-013（self）声称 closed ×2 是否可独立成立
verdict: pass
open_required: 0
version: 1.0.0
---

# A-014 · A-012 F-007/F-008 关闭证据独立复审

## 结论摘要

- **source**：independent（focused finding-closure）
- **verdict**：**pass**
- **open required**：**0**
- A-013（self）声称「两项 required 全部 closed（fixed ×2）」**可独立核对成立**：F-007/F-008 的必改要求逐条与现行代码证据对证通过，且本轮对 3 项定向测试、2 个全包套件与 `go build` 的实际执行全部通过；真实 PostgreSQL 测试**实际执行、未 skip**。
- 依据 A-012 的 AI 推荐门禁三核对项（①SQLite/PG Apply 中途失败零残留 + corrected reopen；②Telegram-enabled 真实组合根与同一 runtime 命令注册；③结构化 Manifest page/route/schemaUrl/navigation 断言），三项全部成立。
- 本意见只追加审计台账；不修改目标 `status` / `progress`、goal-tree、决策/执行台账、合同正文或业务代码。

## 范围与区间

- 被审目标：`docs/workspaces/workspace-031-digital-offer-entitlement/GOAL-005-r4-evidence-closeout/`（workspace-031 内）。
- 复审分母：A-012（`03-audit/A-012-independent-runtime-closure-reaudit.md`）的 **F-007**（required · medium）与 **F-008**（required · medium）及其必改要求。
- 响应记录：A-013（`03-audit/A-013-self-response-a012.md`，self · closed ×2）、决策 D-004、执行 E-008。
- 本复审**不**复核 A-012 已确认关闭的 A-006 范围（F-003～F-006），也不评审整个目标的关门结论——关门与否由 `/govern` 汇总本意见后按 P-003 处理。
- 只读本工作区上下文；未读取/比较其他工作区。

## 关闭证据核对表

### F-007 · required · medium · 0070 多语句 Apply 中途失败后「无 schema/ledger 残留 + corrected catalog reopen 成功」

| A-012 必改要求 | 现行证据路径 | 结论 |
|---|---|---|
| ① 在 `internal/store` 增加 runner 级失败测试：多语句 migration 首条创建对象后返回注入错误，断言该 migration 的 schema 对象与 ledger 行均不存在 | `apps/api/internal/store/migrate_failure_reopen_test.go`：`f007Catalog(true)` 的 v2 Apply 先执行 `CREATE TABLE f007_fail_items` + `CREATE INDEX idx_f007_fail_items_offer`（两条 DDL）后返回注入错误；`TestMigrateMidApplyFailureNoResidueThenReopen`（95-169 行）断言表不存在（109-112）、索引不存在（113-116）、ledger 恰为 `{1}`（117-134）——v1 已提交、v2 全部回滚 | **成立** |
| ② 修正后同版本/同目录重开：断言迁移成功且仅一条 ledger 记录 | 同一文件 `f007Catalog(false)`：与失败版本**同 version 2、同 name、同 checksum**（checksum 由 DDL 串计算，与注入错误无关），仅移除注入错误；重开断言 applied=`{1,2}`（149-155）、v2 ledger 恰 1 行（156-162）、`verifyIntegrity` 通过（146-148）、表与索引重建（163-168） | **成立** |
| ③ SQLite 必测；真实 PG 同语义进集成测试；无需 Down、无需 profile-gated | SQLite：`TestMigrateMidApplyFailureNoResidueThenReopen` PASS。PG：`TestPostgresMigrateMidApplyFailureNoResidueThenReopen`（174-249 行）走 `pgtest.DSN()`（凭据来自 `apps/api/configs/.env` 的 `PG_TEST_*`）+ `scratchDSN` 独立 scratch 库，断言 `to_regclass` 为空、ledger=1 → 重开后 v2=1 行、`string_agg` 版本=`1,2`；本轮实际执行 PASS 未 skip。未新增 Down、未改 profile-gated（D-004 §F-007.4） | **成立** |

事务边界复核（runner 实现与测试声明一致）：`apps/api/internal/store/migrate.go` `applyMigration`（108-132 行）把模块 Apply 与 ledger insert 放在同一 `*sql.Tx`，Apply 返回错误即 `Rollback`；PG 侧 `apps/api/internal/store/postgres.go` `applyMigrationPG`（154-173 行）同语义（`p.Run` 单事务）。测试断言与 A-013 描述逐字吻合（表+索引+ledger 行零残留、ledger=`{1}`、重开=v2 恰 1 行、`verifyIntegrity`）。

### F-008 · required · medium · Telegram-enabled 真实组合根 + 结构化 Manifest

| A-012 必改要求 | 现行证据路径 | 结论 |
|---|---|---|
| ① `biz.digital-offer + channel.telegram` 真实 `newAppWithOptions`/Fx 组合测试，用现有 fake/runtime seam，禁真实外网 Bot API | `apps/api/internal/composition/composition_digitaloffer_telegram_test.go`：`TestDigitalOfferTelegramCompositionRoot`（86-278 行）经真实 `newAppWithOptions`/Fx 装配，`fx.Supply(&telegramRuntimeOptions{APIBaseURL:"https://telegram.test", HTTPClient:&http.Client{Transport: client}})`；fake client（42-84 行）只应答 `getMe`/`deleteWebhook`/`sendMessage`，未知路径 fail closed——零外网依赖。seam 生产不变性：`composition.go:963` `NewHTTPSender(rt, options.HTTPClient, options.APIBaseURL)`；非测试路径 options 恒为零值（`composition.go:923`、`932`，全仓 `telegramRuntimeOptions{非零}` 仅出现在 `_test.go`），`NewHTTPSender` 对 nil client / 空 URL 默认真实 Bot API（`http_sender.go:39-51`） | **成立** |
| ② price/buy/entitlements handlers 注册在 webhook/Telegram surface 使用的同一 dispatcher；至少一个命令经该组合根可达 | 代码链：`buildTelegramRuntime` 构造**唯一** `disp`（`composition.go:951`）→ webhook `HandlerConfig.Dispatcher=disp`（977 行）→ `TelegramRuntime.DispatcherState=disp`（986 行）→ digital-offer 分支 `tgDispatcher = tr.DispatcherState`（631-633 行）→ `digitaloffer.Provider.Register` → `Service.RegisterTelegram`（`apps/api/modules/digitaloffer/service/service.go:748-764`）注册 price/buy/entitlements 三命令；webhook 路由 `POST /api/channel/telegram/webhook` 经 `telegrammodule.New(tr.Webhook,...)`（`composition.go:693` → `apps/api/modules/channel/telegram/provider.go:182-195`）挂到同一 mux，分发路径 `h.dispatcher.Dispatch`（`apps/api/internal/channel/telegram/webhook.go:262-267`）。测试证据：`tr.DispatcherState.HasBusinessHandlers()` 为真（129-134）；`/price` 经**真实 webhook HTTP 面**（secret 头校验 → bot identity → subject 映射 → inbound 持久化 → dispatcher → 出站回包）返回 200 且 fake client 捕获含「在售数字服务」的回包（140-150）；`/entitlements`、`/buy` 经同一 dispatcher 直驱，回包「您当前没有有效权益」「用法：/buy」被捕获（154-170） | **成立**（三命令同一 dispatcher + `/price` 经真实 webhook 面可达，超出「至少一个」最低要求） |
| ③ Manifest 响应解析为结构，断言两个 page 的 route、schemaUrl、navigation refs；保留 schema/HTTP fail-closed 断言 | 测试 179-217 行：`/.well-known/schema-ui/app-manifest.json` 响应 `json.Unmarshal` 为结构（`pages[].pageId/route/schemaUrl` + `navigation.sidebar[].pageRef`），断言 `digitaloffer-offers` / `digitaloffer-entitlements` 两 page 的 route 与 schemaUrl **精确匹配**、sidebar 两 pageRef 均存在——非子串搜索；`apps/api/modules/digitaloffer/manifest/fragment.json` 与断言值一致，`apps/api/internal/manifest/manifest.go` `Aggregate`（92-186 行）为真结构化聚合。fail-closed 保留：schema 匿名 401 / 认证 200 / 未知页匿名 401 / 认证 404（221-247），admin 401/200、public catalog 200（249-266），`DELETE /api/digitaloffer/offers/{id}` → 405/404 负向（268-277，含 A-012 非 required 建议 1） | **成立** |

## Findings

### F-001 · recommended · low · 两条组合根测试的 Manifest 断言强度不对齐；同一 dispatcher 可加显式指针断言

- **证据**：disabled 分支 `TestDigitalOfferCompositionRoot`（`composition_digitaloffer_test.go:116-118`）对 Manifest 仍是 `strings.Contains` 子串断言，而新的 Telegram-enabled 测试已升级为结构化断言（`composition_digitaloffer_telegram_test.go:179-217`）。另外，「同一 dispatcher」目前由单一构造保证（`buildTelegramRuntime` 只构造一次 `disp`）加功能可达性（webhook `/price` 回包被捕获）共同证明，测试未做 `tr.DispatcherState` 与 webhook 挂载 dispatcher 的指针同一显式断言。
- **影响**：低。两条测试的断言强度差异可能让未来对 disabled 分支的 Manifest 改动（如 route/schemaUrl 漂移）逃过结构化检查；指针同一断言属锦上添花。
- **建议**：可选——把 disabled 分支测试的 Manifest 断言升级为与 Telegram-enabled 测试相同的结构解析；为 Telegram-enabled 测试增加 dispatcher 指针同一断言。二者均为 recommended，不阻塞关门。

## 必改项汇总

- **open required：0**。A-012 的 F-007 / F-008（均 required · medium）关闭证据均成立。
- 无新增 required finding。

## 与既有意见的异同

- **A-013（self · response）声称 closed ×2 是否成立**：**成立**。A-013 对两项 finding 的断言细节与实际测试代码逐条吻合（SQLite：`f007_fail_items` 表与索引、ledger=`{1}`、`verifyIntegrity`；PG：`to_regclass` 为空、ledger=1 → 重开 =2、`string_agg`=`1,2`）；证据路径与 A-012 必改要求一一对应；未加 Down、未改 profile-gated，与 A-012 约束一致。
- 与 A-012（independent · conditional）的关系：本复审把 A-012 指出的「把实现推断当成失败路径/启用分支事实」的缺口逐一补实为可执行切片并实际执行通过；A-012 的 AI 推荐门禁三核对项全部满足。
- A-012 非 required 建议 2/3（`digitaloffer.offer.manage` / `digitaloffer.entitlement.void` 权限分离测试；registry/Provider descriptor 等价测试）仍未落实，但 A-012 明确不阻塞关门；维持 recommended 状态，不纳入本 verdict。
- A-013 自身 frontmatter `status: closed` 为响应记录的状态；正式 finding 闭合（fixed）已由 D-004 决策 + E-008 事实 + A-013 登记留痕，本复审确认其证据可核对成立，最终关门放行由 `/govern` 按 P-003 执行。

## 本轮实际验证

工作目录 `C:\Users\magicvr\Documents\Code\schema-ui-core\apps\api`（Go 1.26.0）；今日日期 2026-09-05。以下输出为本会话**实际执行**记录，非引用既有台账。

1. **SQLite 失败/reopen（F-007 ①）**
   `go test -count=1 ./internal/store -run 'TestMigrateMidApplyFailureNoResidueThenReopen$' -v`
   → `--- PASS: TestMigrateMidApplyFailureNoResidueThenReopen (0.04s)`；`ok github.com/magicvr/schema-ui-core/apps/api/internal/store 0.184s`

2. **真实 PostgreSQL 失败/reopen（F-007 ③）**
   `go test -count=1 ./internal/store -run 'TestPostgresMigrateMidApplyFailureNoResidueThenReopen$' -v`
   → `--- PASS: TestPostgresMigrateMidApplyFailureNoResidueThenReopen (0.31s)`；`ok ... 0.460s`
   **PG 实际执行、未 skip**：`apps/api/configs/.env` 中 `PG_TEST_PASSWORD` 非空（已核实，未打印值），`pgtest.DSN()` 生效；`scratchDSN`（`postgres_test.go:470-493`）创建独立 scratch 库并在 cleanup 时 drop；0.31s 耗时为真实 PG 往返。若环境不可用该测试会走 `t.Skip`（175-178 行），本回合未触发。

3. **Telegram-enabled 组合根（F-008）**
   `go test -count=1 ./internal/composition -run '^TestDigitalOfferTelegramCompositionRoot$' -v`
   → `--- PASS: TestDigitalOfferTelegramCompositionRoot (0.26s)`；`ok ... 0.495s`

4. **加固（可选）**
   - `go test -count=1 ./internal/store` → `ok ... 54.268s`
   - `go test -count=1 ./internal/composition` → `ok ... 26.815s`
   - `go build ./...` → 退出码 0

## 结论 + 建议给编排器/用户的下一步

- **verdict：pass（open required：0）**——A-012 两项 required findings（F-007/F-008）的关闭证据充分、可重复核对，A-013 的 closed ×2 声明成立。
- 建议 `/govern`：本复审通过后，若 GOAL-005 / Root GOAL-001 / VP-031 无其他开放必改或信息门禁，可按 A-012 建议重新关门（GOAL-005 C2 2/2、Root 4/4、VP-031 closed），并将 workspace-031 作为依赖生产迁移、Telegram 命令与 Admin surface 的后继 VP 的已验证前置。F-001 为 recommended 可选加固，不阻塞。
- 若后继 VP 完全不依赖上述能力而希望以书面 `accepted-residual` 限定范围，仍由 `/govern` 提交用户裁决（A-012 同口径）；本独立意见不替用户接受残余。

## 声明

本意见只追加到 GOAL-005 审计台账（A-014 + `03-audit.md` 索引行），**不修改**目标 `status` / `progress` / 检查点、`goal-tree.md`、`00-meta.md`、决策/执行台账、合同正文或业务代码；响应与关门由 `/govern` 处理。
